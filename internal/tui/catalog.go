package tui

// Repo beschrijft een CEDA tool die geïnstalleerd kan worden.
type Repo struct {
	ID          string
	Name        string
	Description string
	RepoURL     string
	Script      string
}

// Catalog is de hardcoded lijst van beschikbare CEDA tools.
var Catalog = []Repo{
	{
		ID:          "1cijferho",
		Name:        "1CijferHO",
		Description: "Snel en zorgvuldig aan de slag met 1cijferHO-data voor onderwijsanalyses en onderzoek.",
		RepoURL:     "https://github.com/cedanl/1cijferho.git",
		Script:      "setup.ps1",
	},
	{
		ID:          "no-fairness-without-awareness",
		Name:        "No Fairness Without Awareness",
		Description: "Het NFWA (No Fairness Without Awareness) package is een R package ontwikkeld op basis van het onderzoek van het lectoraat Learning Technology & Analytics (LTA) van De Haagse Hogeschool.",
		RepoURL:     "https://github.com/cedanl/no-fairness-without-awareness.git",
		Script:      "setup.ps1",
	},
}
