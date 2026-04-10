package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectType geeft aan of een repo R of uv is.
type ProjectType string

const (
	ProjectTypeR       ProjectType = "r"
	ProjectTypeUV      ProjectType = "uv"
	ProjectTypeUnknown ProjectType = "unknown"
)

// RunStep beschrijft één stap in de run-flow.
type RunStep struct {
	Label    string
	StepName string
}

// CommonSteps zijn altijd de eerste stappen, ongeacht projecttype.
var CommonSteps = []RunStep{
	{Label: "Scoop controleren", StepName: "scoop-check"},
	{Label: "Core dependencies installeren", StepName: "core-deps"},
	{Label: "Scoop buckets toevoegen", StepName: "buckets"},
	{Label: "Projecttype detecteren", StepName: "detect"},
}

// RSteps zijn de stappen voor R-projecten.
var RSteps = []RunStep{
	{Label: "R, Rtools en Positron installeren", StepName: "r-install"},
	{Label: "Rtools paden instellen", StepName: "r-paths"},
	{Label: "Positron R interpreter instellen", StepName: "r-positron"},
	{Label: "R packages installeren via renv", StepName: "r-sync"},
	{Label: "Project openen in Positron", StepName: "r-run"},
}

// UVSteps zijn de stappen voor uv-projecten.
var UVSteps = []RunStep{
	{Label: "uv installeren", StepName: "uv-install"},
	{Label: "pyproject.toml configureren", StepName: "uv-config"},
	{Label: "Packages installeren via uv sync", StepName: "uv-sync"},
	{Label: "Project starten", StepName: "uv-run"},
}

// StepsForType geeft de volledige stap-lijst voor een bekend projecttype.
func StepsForType(pt ProjectType) []RunStep {
	switch pt {
	case ProjectTypeR:
		return append(CommonSteps, RSteps...)
	case ProjectTypeUV:
		return append(CommonSteps, UVSteps...)
	default:
		return CommonSteps
	}
}

// DetectFromOutput leest de TYPE: regel uit de detect-stap output.
// Returnt ProjectTypeUnknown als het niet gevonden wordt.
func DetectFromOutput(output string) ProjectType {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, "TYPE:") {
			t := strings.TrimPrefix(line, "TYPE:")
			switch strings.TrimSpace(t) {
			case "r":
				return ProjectTypeR
			case "uv":
				return ProjectTypeUV
			}
		}
	}
	return ProjectTypeUnknown
}

// ExecuteStep voert één stap uit via ceda-run.ps1.
// Returnt (output, error). Output is altijd gevuld voor context bij fouten.
func ExecuteStep(step RunStep, installPath string) (string, error) {
	scriptPath := filepath.Join(modulesDir(), "ceda-run.ps1")

	if _, err := os.Stat(scriptPath); err != nil {
		return "", fmt.Errorf("ceda-run.ps1 niet gevonden op: %s", scriptPath)
	}

	cmd := exec.Command("powershell.exe",
		"-NoProfile", "-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-File", scriptPath,
		"-Step", step.StepName,
		"-Root", installPath,
	)

	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		lines := parseOutput(output)
		if len(lines) > 0 {
			return output, fmt.Errorf("%s\n%s", err.Error(), strings.Join(lines, "\n"))
		}
		return output, err
	}
	return output, nil
}

// modulesDir geeft het pad naar scripts/windows/modules.
func modulesDir() string {
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "scripts", "windows", "modules")
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			return candidate
		}
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, "scripts", "windows", "modules")
}

func parseOutput(raw string) []string {
	var lines []string
	for _, l := range strings.Split(raw, "\n") {
		if t := strings.TrimRight(l, "\r"); strings.TrimSpace(t) != "" {
			lines = append(lines, t)
		}
	}
	return lines
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
