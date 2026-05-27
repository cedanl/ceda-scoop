package main

import (
	"fmt"
	"os"

	"github.com/cedanl/ceda-scoop/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

const Version = "0.2.1"

func main() {
	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kon debug.log niet openen:", err)
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	p := tea.NewProgram(ui.New(Version), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
