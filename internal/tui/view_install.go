package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewInstall() string {
	if m.selectedRepo == nil {
		return ""
	}

	navbar := m.renderSimpleNavbar("Installeren")
	width := m.width - 8

	// ── Status ────────────────────────────────────────────────────────────────
	var statusLine string
	if m.installDone {
		if m.installErr == "" {
			statusLine = lipgloss.NewStyle().Bold(true).Foreground(ColorGreen).Render("✓  Installatie voltooid")
		} else {
			statusLine = lipgloss.NewStyle().Bold(true).Foreground(ColorRed).Render("✗  " + m.installErr)
		}
	} else {
		statusLine = fmt.Sprintf("%s  Bezig met installeren...", m.spinner.View())
	}

	// ── Progress ──────────────────────────────────────────────────────────────
	// Bij done altijd 1.0 tonen, anders de lopende waarde
	displayPct := m.installPercent
	if m.installDone {
		displayPct = 1.0
	}

	pct := int(displayPct * 100)
	pctLabel := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render(fmt.Sprintf("%3d%%", pct))

	progressBar := lipgloss.NewStyle().Width(width - 6).Render(
		m.progress.ViewAs(displayPct),
	)
	progressRow := lipgloss.JoinHorizontal(lipgloss.Center, progressBar, "  ", pctLabel)

	// ── Stappen ───────────────────────────────────────────────────────────────
	// Drempels afgestemd op de asymptotische progress curve:
	// clone start ~0%, loopt naar ~60-70% na 15-20s, bootstrap daarna
	steps := []struct {
		label string
		pct   float64
	}{
		{"Verbinding maken", 0.03},
		{"Repository clonen", 0.30},
		{"Omgeving instellen", 0.75},
		{"Afronden", 0.95},
	}

	var stepsView string
	for _, s := range steps {
		var dot, label string
		done := m.installDone && m.installErr == ""
		if done || displayPct >= s.pct {
			dot = lipgloss.NewStyle().Foreground(ColorGreen).Render("●")
			if done {
				label = lipgloss.NewStyle().Foreground(ColorGreen).Render(s.label)
			} else {
				label = lipgloss.NewStyle().Foreground(ColorWhite).Render(s.label)
			}
		} else if displayPct >= s.pct-0.10 {
			dot = lipgloss.NewStyle().Foreground(ColorPrimary).Render("◉")
			label = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(s.label)
		} else {
			dot = StyleSubtle.Render("○")
			label = StyleSubtle.Render(s.label)
		}
		stepsView += dot + "  " + label + "\n"
	}

	// ── Elapsed ───────────────────────────────────────────────────────────────
	elapsed := m.installElapsed
	elapsedStr := fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
	elapsedLabel := StyleSubtle.Render("⏱  " + elapsedStr)

	// ── Repo naam ─────────────────────────────────────────────────────────────
	repoName := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render(m.selectedRepo.repo.Name)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(1, 2).
		Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			repoName, "",
			statusLine, "",
			progressRow, "",
			stepsView,
			elapsedLabel,
		))

	var hint string
	if m.installDone {
		hint = StyleHelp.Render("[Enter] terug naar Library")
	} else {
		hint = StyleHelp.Render("Even geduld — installatie loopt...")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(box),
		"",
		lipgloss.NewStyle().PaddingLeft(2).Render(hint),
	)
}
