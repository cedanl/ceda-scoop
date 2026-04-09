package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) viewProjectTypePicker() string {
	navbar := m.renderSimpleNavbar("Projecttype kiezen")

	title := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).
		Render("Beide projecttypen gevonden")

	sub := StyleSubtle.Render("Zowel renv.lock als uv.lock aanwezig — welk type wil je starten?")

	optR := m.pickerOption("R", "renv.lock  →  Positron + renv::restore", m.pickerSelected == 0)
	optUV := m.pickerOption("uv", "uv.lock  →  uv sync + project starten", m.pickerSelected == 1)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(1, 3).
		Width(m.width - 8).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			sub, "",
			optR, "",
			optUV,
		))

	hint := StyleHelp.Render("  ↑↓ selecteren  [Enter] bevestigen  [1] R  [2] uv  [Esc] terug")

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(box),
		"",
		lipgloss.NewStyle().PaddingLeft(2).Render(hint),
	)
}

func (m Model) pickerOption(key, desc string, active bool) string {
	if active {
		dot := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("◉")
		k := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(key)
		d := lipgloss.NewStyle().Foreground(ColorWhite).Render(desc)
		return dot + "  " + k + "  " + StyleSubtle.Render("·") + "  " + d
	}
	dot := StyleSubtle.Render("○")
	k := StyleSubtle.Render(key)
	d := StyleSubtle.Render(desc)
	return dot + "  " + k + "  " + StyleSubtle.Render("·") + "  " + d
}
