package tui

import "github.com/cedanl/ceda-scoop/internal/runner"

// Screen representeert het actieve TUI-scherm.
type Screen int

const (
	ScreenSplash Screen = iota
	ScreenStore
	ScreenLibrary
	ScreenDetail
	ScreenInstall
	ScreenRun
	ScreenSettings
	ScreenDeleteConfirm
	ScreenProjectTypePicker
)

// Install messages
type SplashDoneMsg struct{}
type InstallDoneMsg struct{ Err string }
type installTickMsg struct{}
type DeleteDoneMsg struct{ Err string }

// Run messages
type RunStepDoneMsg struct {
	StepIdx      int
	Err          string
	DetectedType runner.ProjectType // alleen ingevuld na detect-stap
}
type runTickMsg struct{}
