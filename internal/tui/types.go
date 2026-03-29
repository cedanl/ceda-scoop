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
)

// Repo beschrijft een CEDA tool die geïnstalleerd kan worden.
type Repo struct {
	ID          string
	Name        string
	Description string
	RepoURL     string
	Script      string // naam van het bootstrap script
}

// Catalog is de hardcoded lijst van beschikbare CEDA tools.
// Voeg hier nieuwe repos toe als ze beschikbaar komen.
var Catalog = []Repo{
	{
		ID:          "1cijferho",
		Name:        "1CijferHO",
		Description: "Snel en zorgvuldig aan de slag met 1cijferHO-data voor onderwijsanalyses en onderzoek.",
		RepoURL:     "https://github.com/cedanl/1cijferho",
		Script:      "setup.ps1",
	},
	{
		ID:          "no-fairness-without-awareness",
		Name:        "No Fairness Without Awareness",
		Description: "Het NFWA (No Fairness Without Awareness) package is een R package ontwikkeld op basis van het onderzoek van het lectoraat Learning Technology & Analytics (LTA) van De Haagse Hogeschool.",
		RepoURL:     "https://github.com/cedanl/no-fairness-without-awareness",
		Script:      "setup.ps1",
	},
	// TODO: voeg hier meer repos toe
}

// Bubble Tea messages.
type SplashDoneMsg struct{}
type InstallStepMsg struct{ Line string }
type InstallDoneMsg struct{ Err string }
type InstallProgressMsg struct{ Pct float64 }
