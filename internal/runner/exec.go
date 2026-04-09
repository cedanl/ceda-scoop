package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// DefaultInstallBase geeft het standaard installatiepad terug (~/ceda).
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
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// Bootstrap voert het platform-specifieke setup script uit in repoPath.
func Bootstrap(repoPath, scriptName string, logCh chan<- string) error {
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
