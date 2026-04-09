package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// NavTab beschrijft een tab in de navbar.
type NavTab struct {
	Label    string
	Active   bool
	Shortcut string // optioneel, bv "[S]"
}

// renderNavbar bouwt een volledige-breedte navbar met brand + tabs + rechts shortcuts.
func (m Model) renderNavbar(tabs []NavTab, rightHints string) string {
	brand := StyleNavBrand.Render(" CEDA Store ")

	// Tabs als pills
	var tabParts []string
	for _, t := range tabs {
		if t.Active {
			tabParts = append(tabParts, StyleTabActive.Render(" "+t.Label+" "))
		} else {
			tabParts = append(tabParts, StyleTabInactive.Render(" "+t.Label+" "))
		}
	}
	tabRow := strings.Join(tabParts, StyleNavDivider.Render("│"))

	// Rechts-uitgelijnde hints
	right := lipgloss.NewStyle().
		Foreground(ColorDim).
		Background(ColorSurface).
		Padding(0, 2).
		Render(rightHints)

	// Breedte van de spacer = totale breedte - brand - tabs - rechts
	brandW := lipgloss.Width(brand)
	tabW := lipgloss.Width(tabRow)
	rightW := lipgloss.Width(right)
	spacerW := m.width - brandW - tabW - rightW
	if spacerW < 0 {
		spacerW = 0
	}
	spacer := lipgloss.NewStyle().Background(ColorSurface).Render(strings.Repeat(" ", spacerW))

	navbar := lipgloss.JoinHorizontal(lipgloss.Top,
		brand,
		StyleNavDivider.Render("  "),
		tabRow,
		spacer,
		right,
	)

	// Scheidingslijn onder de navbar
	divider := lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(strings.Repeat("─", m.width))

	return lipgloss.JoinVertical(lipgloss.Left, navbar, divider)
}

// renderSimpleNavbar voor schermen zonder tabs (detail, install, settings).
func (m Model) renderSimpleNavbar(subtitle string) string {
	brand := StyleNavBrand.Render(" CEDA Store ")

	sub := lipgloss.NewStyle().
		Foreground(ColorGray).
		Background(ColorSurface).
		Padding(0, 2).
		Bold(false).
		Render(subtitle)

	spacerW := m.width - lipgloss.Width(brand) - lipgloss.Width(sub)
	if spacerW < 0 {
		spacerW = 0
	}
	spacer := lipgloss.NewStyle().Background(ColorSurface).Render(strings.Repeat(" ", spacerW))

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, brand, sub, spacer)

	divider := lipgloss.NewStyle().
		Foreground(ColorDim).
		Render(strings.Repeat("─", m.width))

	return lipgloss.JoinVertical(lipgloss.Left, navbar, divider)
}
