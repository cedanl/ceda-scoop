package tui

import "github.com/charmbracelet/lipgloss"

// ASCII art gegenereerd met figlet (font: ANSI Shadow)
const cedaASCII = `
 ██████╗███████╗██████╗  █████╗
██╔════╝██╔════╝██╔══██╗██╔══██╗
██║     █████╗  ██║  ██║███████║
██║     ██╔══╝  ██║  ██║██╔══██║
╚██████╗███████╗██████╔╝██║  ██║
 ╚═════╝╚══════╝╚═════╝ ╚═╝  ╚═╝`

func (m Model) viewSplash() string {
	logo := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render(cedaASCII)

	store := lipgloss.NewStyle().
		Foreground(ColorWhite).
		Bold(true).
		Render("  Store")

	tagline := StyleSubtle.Render("Tooling voor Nederlands hoger onderwijs")

	// Geen spinner meer — wacht gewoon op Enter
	hint := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("[ Druk op Enter om te beginnen ]")

	content := lipgloss.JoinVertical(lipgloss.Center,
		logo+store,
		"",
		tagline,
		"",
		"",
		hint,
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
