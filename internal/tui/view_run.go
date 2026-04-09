package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewRun() string {
	if m.selectedRepo == nil {
		return ""
	}

	navbar := m.renderSimpleNavbar("Openen")
	width := m.width - 8

	// ── Repo naam + type ──────────────────────────────────────────────────────
	repoName := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).
		Render(m.selectedRepo.repo.Name)

	typeBadge := lipgloss.NewStyle().Foreground(ColorGray).
		Render(string(m.runProjectType))

	nameRow := lipgloss.JoinHorizontal(lipgloss.Top, repoName, "   ", typeBadge)

	// ── Status ────────────────────────────────────────────────────────────────
	var statusLine string
	switch {
	case m.runDone && m.runErr == "":
		statusLine = lipgloss.NewStyle().Bold(true).Foreground(ColorGreen).
			Render("✓  Klaar — project is gestart")
	case m.runDone && m.runErr != "":
		statusLine = lipgloss.NewStyle().Bold(true).Foreground(ColorRed).
			Render("✗  Er ging iets mis")
	default:
		label := "Bezig..."
		if len(m.runSteps) > 0 && m.runCurrentStep < len(m.runSteps) {
			label = m.runSteps[m.runCurrentStep].Label + "..."
		}
		statusLine = fmt.Sprintf("%s  %s", m.spinner.View(), label)
	}

	// ── Progress ──────────────────────────────────────────────────────────────
	displayPct := m.runPercent
	if m.runDone && m.runErr == "" {
		displayPct = 1.0
	}
	pct := int(displayPct * 100)
	pctLabel := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).
		Render(fmt.Sprintf("%3d%%", pct))
	progressBar := lipgloss.NewStyle().Width(width - 6).
		Render(m.progress.ViewAs(displayPct))
	progressRow := lipgloss.JoinHorizontal(lipgloss.Center, progressBar, "  ", pctLabel)

	// ── Stappen ───────────────────────────────────────────────────────────────
	var stepsView strings.Builder
	for i, step := range m.runSteps {
		var dot, label string

		switch {
		case m.runDone && m.runFailedStep == i:
			// Deze stap is mislukt
			dot = lipgloss.NewStyle().Foreground(ColorRed).Render("✗")
			label = lipgloss.NewStyle().Foreground(ColorRed).Bold(true).Render(step.Label)

		case m.runDone && m.runErr != "" && i > m.runFailedStep:
			// Stappen na de fout — niet uitgevoerd
			dot = StyleSubtle.Render("○")
			label = StyleSubtle.Render(step.Label)

		case i < m.runCurrentStep || (m.runDone && m.runErr == ""):
			// Voltooid
			dot = lipgloss.NewStyle().Foreground(ColorGreen).Render("●")
			label = lipgloss.NewStyle().Foreground(ColorGreen).Render(step.Label)

		case i == m.runCurrentStep && !m.runDone:
			// Actief
			dot = lipgloss.NewStyle().Foreground(ColorPrimary).Render("◉")
			label = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(step.Label)

		default:
			// Wachtend
			dot = StyleSubtle.Render("○")
			label = StyleSubtle.Render(step.Label)
		}

		stepsView.WriteString(dot + "  " + label + "\n")
	}

	// ── Foutmelding box (alleen bij fout) ─────────────────────────────────────
	var errorBox string
	if m.runDone && m.runErr != "" {
		wrapped := wrapText(m.runErr, width-8)
		errorBox = "\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorRed).
			Foreground(ColorRed).
			Padding(0, 2).
			Width(width - 4).
			Render(wrapped)
	}

	// ── Elapsed ───────────────────────────────────────────────────────────────
	elapsed := m.runElapsed
	elapsedStr := fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
	elapsedLabel := StyleSubtle.Render("⏱  " + elapsedStr)

	// ── Box ───────────────────────────────────────────────────────────────────
	innerContent := lipgloss.JoinVertical(lipgloss.Left,
		nameRow, "",
		statusLine, "",
		progressRow, "",
		stepsView.String(),
		elapsedLabel,
	)
	if errorBox != "" {
		innerContent += errorBox
	}

	borderColor := ColorDim
	if m.runDone && m.runErr != "" {
		borderColor = ColorRed
	} else if m.runDone {
		borderColor = ColorGreen
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(width).
		Render(innerContent)

	// ── Hint ──────────────────────────────────────────────────────────────────
	var hint string
	if m.runDone {
		hint = StyleHelp.Render("[Enter] terug naar detail")
	} else {
		hint = StyleHelp.Render("Even geduld...")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(box),
		"",
		lipgloss.NewStyle().PaddingLeft(2).Render(hint),
	)
}
