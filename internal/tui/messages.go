package tui

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

// Bubble Tea messages — install
type SplashDoneMsg struct{}
type InstallStepMsg struct{ Line string }
type InstallDoneMsg struct{ Err string }
type InstallProgressMsg struct{ Pct float64 }
type installTickMsg struct{}
type DeleteDoneMsg struct{ Err string }

// Bubble Tea messages — run
// RunStepDoneMsg wordt gestuurd als één PS-stap klaar is.
type RunStepDoneMsg struct {
	StepIdx int
	Err     string // leeg = succes
}
type runTickMsg struct{}
