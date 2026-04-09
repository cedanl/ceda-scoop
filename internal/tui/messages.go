package tui

// Screen representeert het actieve TUI-scherm.
type Screen int

const (
	ScreenSplash Screen = iota
	ScreenStore
	ScreenLibrary
	ScreenDetail
	ScreenInstall
	ScreenSettings
	ScreenDeleteConfirm
)

// Bubble Tea messages.
type SplashDoneMsg struct{}
type InstallStepMsg struct{ Line string }
type InstallDoneMsg struct{ Err string }
type InstallProgressMsg struct{ Pct float64 }
type installTickMsg struct{}
type DeleteDoneMsg struct{ Err string }
