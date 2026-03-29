package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// DefaultInstallBase geeft het standaard installatiepad terug (~/<username>/ceda).
func DefaultInstallBase() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "ceda")
}

// OpenInExplorer opent het gegeven pad in de bestandsverkenner van het OS.
func OpenInExplorer(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	_ = cmd.Start()
}

// Clone doet een git clone van repoURL naar destPath.
func Clone(repoURL, destPath string) error {
	cmd := exec.Command("git", "clone", repoURL, destPath)
	return cmd.Run()
}

// Bootstrap voert het platform-specifieke setup script uit in repoPath.
// Op Windows: PowerShell, op macOS/Linux: sh
// TODO: scripts inbakken via go:embed en tijdelijk uitpakken voor uitvoering
func Bootstrap(repoPath, scriptName string) error {
	scriptPath := filepath.Join(repoPath, scriptName)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	default:
		cmd = exec.Command("sh", scriptPath)
	}

	cmd.Dir = repoPath
	return cmd.Run()
}
