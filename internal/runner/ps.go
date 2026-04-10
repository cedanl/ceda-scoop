package runner

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

// runPSCommand voert een inline PowerShell commando uit.
// stdout en stderr worden concurrent gelezen om deadlocks te voorkomen.
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

	var mu sync.Mutex
	var lines []string
	var errBuf bytes.Buffer
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if t := scanner.Text(); t != "" {
				mu.Lock()
				lines = append(lines, t)
				mu.Unlock()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(&errBuf, stderr)
	}()

	wg.Wait()
	cmdErr := cmd.Wait()

	if errBuf.Len() > 0 {
		scanner := bufio.NewScanner(&errBuf)
		for scanner.Scan() {
			if t := scanner.Text(); t != "" {
				mu.Lock()
				lines = append(lines, t)
				mu.Unlock()
			}
		}
	}

	return lines, cmdErr
}
