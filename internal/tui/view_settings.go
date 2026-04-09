package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) viewSettings() string {
	navbar := m.renderSimpleNavbar("Instellingen")

	currentPath := m.InstallBase
	if m.settingsInput != "" {
		currentPath = m.settingsInput
	}

	var pathField string
	if m.editingPath {
		pathField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1).
			Render(currentPath + "█")
	} else {
		pathField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(0, 1).
			Render(currentPath)
	}

	var editHint string
	if m.editingPath {
		editHint = StyleSubtle.Render("Typ het pad en druk [Enter] om op te slaan")
	} else {
		editHint = StyleSubtle.Render("[Enter] om te bewerken")
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		StyleBold.Render("Standaard installatiepad:"),
		"", pathField, editHint,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(body),
		"",
		StyleHelp.Render("  [B/Esc] terug"),
	)
}
