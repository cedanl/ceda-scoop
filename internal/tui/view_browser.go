package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewBrowser(library bool) string {
	navbar := m.renderNavbar([]NavTab{
		{Label: "Store", Active: !library},
		{Label: "Library", Active: library},
	}, "[S] instellingen  [Q] afsluiten")

	// ── Items voor deze tab ───────────────────────────────────────────────────
	var items []repoItem
	for _, it := range m.repoItems {
		if library && !it.installed {
			continue
		}
		items = append(items, it)
	}

	// ── Lege staat ────────────────────────────────────────────────────────────
	if len(items) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(ColorGray).
			PaddingLeft(2).
			MarginTop(2).
			Render("Nog niets geïnstalleerd. Ga naar Store om te beginnen.")
		help := StyleHelp.Render("  [Tab] wissel  [Q] afsluiten")
		return lipgloss.JoinVertical(lipgloss.Left, navbar, "", empty, help)
	}

	// ── Kaartjes ──────────────────────────────────────────────────────────────
	cardWidth := m.width - 6
	if cardWidth < 40 {
		cardWidth = 40
	}

	var cards []string
	for idx, it := range items {
		selected := idx == m.selectedCard

		var badge string
		if it.installed {
			badge = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true).Render("✓ Geïnstalleerd")
		} else {
			badge = lipgloss.NewStyle().Foreground(ColorGray).Render("○ Beschikbaar")
		}

		nameStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorWhite)
		if selected {
			nameStyle = nameStyle.Foreground(ColorPrimary)
		}
		name := nameStyle.Render(it.repo.Name)

		desc := lipgloss.NewStyle().
			Foreground(ColorGray).
			Width(cardWidth - 6).
			Render(wrapText(it.repo.Description, cardWidth-6))

		url := lipgloss.NewStyle().Foreground(ColorDim).Render("⎘  " + it.repo.RepoURL)

		inner := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, name, "   ", badge),
			"",
			desc,
			"",
			url,
		)

		borderColor := ColorDim
		if selected {
			borderColor = ColorPrimary
		}

		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(cardWidth).
			Render(inner)

		cards = append(cards, card)
	}

	cardsView := lipgloss.NewStyle().
		PaddingLeft(2).
		Render(strings.Join(cards, "\n"))

	help := StyleHelp.Render("  ↑↓ navigeer  [Enter] details  [Tab] wissel tab")

	return lipgloss.JoinVertical(lipgloss.Left,
		navbar, "",
		cardsView,
		help,
	)
}

// wrapText breekt tekst af op maxWidth tekens, op woordgrenzen.
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, w := range words {
		if current == "" {
			current = w
		} else if len(current)+1+len(w) <= maxWidth {
			current += " " + w
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}
