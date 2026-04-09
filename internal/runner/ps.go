package runner

import (
	"bufio"
	"io"
	"os/exec"
	"runtime"
)

// runPSCommand voert een PowerShell commando uit en returnt output + error.
func runPSCommand(command string) ([]string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe",
			"-NoProfile", "-NonInteractive",
			"-ExecutionPolicy", "Bypass",
			"-Command", command,
		)
	} else {
		cmd = exec.Command("pwsh",
			"-NoProfile", "-NonInteractive",
			"-Command", command,
		)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var lines []string
	for _, r := range []io.ReadCloser{stdout, stderr} {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if t := scanner.Text(); t != "" {
				lines = append(lines, t)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return lines, err
	}
	return lines, nil
}
