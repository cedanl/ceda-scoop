package tui

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// View is het centrale render-punt — Bubble Tea roept dit elke frame aan.
func (m Model) View() string {
	switch m.CurrentScreen {
	case ScreenSplash:
		return m.viewSplash()
	case ScreenStore:
		return m.viewBrowser(false)
	case ScreenLibrary:
		return m.viewBrowser(true)
	case ScreenDetail:
		return m.viewDetail()
	case ScreenInstall:
		return m.viewInstall()
	case ScreenSettings:
		return m.viewSettings()
	}
	return ""
}

func (m Model) viewSplash() string {
	logo := StyleTitle.Render(" CEDA Store ")
	sub := StyleSubtle.Render("Tooling voor Nederlands hoger onderwijs")
	hint := StyleHelp.Render("Druk op Enter om te beginnen...")
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, logo, "", sub, "", hint),
	)
}

func (m Model) viewBrowser(library bool) string {
	storeTab := StyleTabInactive.Render("Store")
	libTab := StyleTabInactive.Render("Library")
	if !library {
		storeTab = StyleTabActive.Render("Store")
	} else {
		libTab = StyleTabActive.Render("Library")
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		StyleTitle.Render(" CEDA Store "),
		"  ",
		storeTab, libTab,
	)

	var listView string
	if !library {
		if len(m.storeList.Items()) == 0 {
			listView = StyleSubtle.Render("\n  Geen tools beschikbaar.")
		} else {
			listView = m.storeList.View()
		}
	} else {
		if len(m.libList.Items()) == 0 {
			listView = StyleSubtle.Render("\n  Nog niets geïnstalleerd. Ga naar Store om te beginnen.")
		} else {
			listView = m.libList.View()
		}
	}

	help := StyleHelp.Render("↑↓ navigeer  [Enter] details  [Tab] wissel tab  [S] instellingen  [Q] afsluiten")

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(listView),
		help,
	)
}

func (m Model) viewDetail() string {
	if m.selectedRepo == nil {
		return ""
	}
	r := m.selectedRepo

	header := StyleTitle.Render(" CEDA Store ")
	title := StyleBold.Copy().Foreground(ColorWhite).Render(r.repo.Name)

	var status, actions string
	if r.installed {
		status = StyleBadgeInstalled.Render("✓ Geïnstalleerd")
		actions = lipgloss.JoinVertical(lipgloss.Left,
			"",
			StyleSubtle.Render("Pad: "+r.installPath),
			"",
			StyleButton.Render(" [O] Openen in Verkenner "),
			"",
			StyleButton.Render(" [U] Updaten "),
			"",
			StyleDanger.Render(" [D] Verwijderen "),
		)
	} else {
		status = StyleBadgeAvailable.Render("─ Nog niet geïnstalleerd")
		actions = lipgloss.JoinVertical(lipgloss.Left,
			"",
			StyleSubtle.Render("Installeert naar: "+filepath.Join(m.InstallBase, r.repo.ID)),
			"",
			StyleButton.Render(" [I] Installeren "),
		)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		title, status, "",
		StyleSubtle.Render(r.repo.Description),
		StyleSubtle.Render("📦 "+r.repo.RepoURL),
		actions,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(body),
		"",
		StyleHelp.Render("[B/Esc] terug  [Q] afsluiten"),
	)
}

func (m Model) viewInstall() string {
	if m.selectedRepo == nil {
		return ""
	}

	header := StyleTitle.Render(" CEDA Store ")
	title := StyleBold.Render("Installeren: " + m.selectedRepo.repo.Name)

	var statusLine string
	if m.installDone {
		if m.installErr == "" {
			statusLine = StyleBadgeInstalled.Render("✓ Klaar! Je vindt het nu in je Library.")
		} else {
			statusLine = StyleDanger.Render("✗ Fout: " + m.installErr)
		}
	} else {
		statusLine = fmt.Sprintf("%s Bezig...", m.spinner.View())
	}

	logBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Render(m.logView.View())

	var hint string
	if m.installDone {
		hint = StyleHelp.Render("[Enter] terug naar Library")
	} else {
		hint = StyleHelp.Render("Installatie bezig, even geduld...")
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		title, "", statusLine, "",
		m.progress.ViewAs(m.installPercent),
		"", logBox,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(body),
		hint,
	)
}

func (m Model) viewSettings() string {
	header := StyleTitle.Render(" CEDA Store — Instellingen ")

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
		header, "",
		lipgloss.NewStyle().PaddingLeft(2).Render(body),
		"",
		StyleHelp.Render("[B/Esc] terug"),
	)
}
