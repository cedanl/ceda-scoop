package tui

import (
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewDetail() string {
	if m.selectedRepo == nil {
		return ""
	}
	r := m.selectedRepo

	navbar := m.renderSimpleNavbar(r.repo.Name)

	var badge string
	if r.installed {
		badge = lipgloss.NewStyle().Bold(true).Foreground(ColorGreen).Render("✓ Geïnstalleerd")
	} else {
		badge = lipgloss.NewStyle().Foreground(ColorGray).Render("○ Beschikbaar")
	}

	name := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render(r.repo.Name)

	descWidth := m.width - 12
	if descWidth < 30 {
		descWidth = 30
	}
	desc := lipgloss.NewStyle().Foreground(ColorWhite).Width(descWidth).
		Render(wrapText(r.repo.Description, descWidth))

	url := StyleSubtle.Render("⎘  " + r.repo.RepoURL)

	var actionsView string
	if r.installed {
		actionsView = lipgloss.JoinVertical(lipgloss.Left,
			"",
			StyleSubtle.Render("Pad:  "+r.installPath),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top,
				StyleButton.Render(" [O] Openen "),
				"  ",
				StyleButton.Render(" [F] Folder "),
				"  ",
				StyleButton.Render(" [U] Updaten "),
				"  ",
				StyleDanger.Render(" [D] Verwijderen "),
			),
		)
	} else {
		actionsView = lipgloss.JoinVertical(lipgloss.Left,
			"",
			StyleSubtle.Render("Installeert naar:  "+filepath.Join(m.InstallBase, r.repo.ID)),
			"",
			StyleButton.Render(" [I] Installeren "),
		)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(1, 3).
		Width(m.width - 6).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, name, "   ", badge),
			"",
			desc,
			"",
			url,
			actionsView,
		))

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(box),
		"",
		StyleHelp.Render("  [B/Esc] terug  [Q] afsluiten"),
	)
}
