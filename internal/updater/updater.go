package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	apiURL  = "https://api.github.com/repos/cedanl/ceda-scoop/releases/latest"
	timeout = 5 * time.Second
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// CheckLatest fetches the latest release tag from GitHub. Returns "" on error.
func CheckLatest() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	return strings.TrimPrefix(r.TagName, "v"), nil
}

// Download fetches the release asset for the current OS/arch and writes it next to the current executable.
// On Unix: replaces the binary and returns the path to restart from.
// On Windows: saves as <binary>-update.exe and returns that path (caller shows restart message).
func Download(version string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	assetName := assetFilename(version)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return "", fmt.Errorf("asset %q niet gevonden in release %s", assetName, version)
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, _ = filepath.EvalSymlinks(exe)

	return downloadBinary(ctx, downloadURL, exe)
}

func assetFilename(version string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	osName := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}[goos]

	archName := map[string]string{
		"amd64": "x86_64",
		"arm64": "arm64",
	}[goarch]

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}

	return fmt.Sprintf("ceda-scoop_%s_%s_%s%s", version, osName, archName, ext)
}

func downloadBinary(ctx context.Context, url, currentExe string) (string, error) {
	resp, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	httpResp, err := http.DefaultClient.Do(resp)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()

	dir := filepath.Dir(currentExe)
	tmp, err := os.CreateTemp(dir, "ceda-scoop-update-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, httpResp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	if runtime.GOOS == "windows" {
		dest := strings.TrimSuffix(currentExe, ".exe") + "-update.exe"
		if err := os.Rename(tmp.Name(), dest); err != nil {
			os.Remove(tmp.Name())
			return "", err
		}
		return dest, nil
	}

	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	if err := os.Rename(tmp.Name(), currentExe); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return currentExe, nil
}
