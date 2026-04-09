package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) viewDeleteConfirm() string {
	if m.selectedRepo == nil {
		return ""
	}

	navbar := m.renderSimpleNavbar("Verwijderen")

	name := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render(m.selectedRepo.repo.Name)
	warning := lipgloss.NewStyle().Bold(true).Foreground(ColorRed).Render("⚠  Weet je zeker dat je dit wilt verwijderen?")
	path := StyleSubtle.Render("Verwijdert: " + m.selectedRepo.installPath)

	yesBtn := lipgloss.NewStyle().
		Bold(true).Foreground(ColorWhite).Background(ColorRed).
		Padding(0, 2).MarginRight(2).
		Render(" [Y] Ja, verwijderen ")

	noBtn := lipgloss.NewStyle().
		Bold(true).Foreground(ColorWhite).
		Border(lipgloss.RoundedBorder()).BorderForeground(ColorGray).
		Padding(0, 2).
		Render(" [N] Annuleren ")

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yesBtn, noBtn)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorRed).
		Padding(1, 3).
		Width(m.width - 8).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			name, "",
			warning, "",
			path, "",
			buttons,
		))

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(box),
		"",
		StyleHelp.Render("  [Y/Enter] bevestigen  [N/Esc] annuleren"),
	)
}
