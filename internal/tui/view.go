package tui

import "fmt"

// View is het centrale render-punt — delegeert naar het juiste scherm.
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
	case ScreenDeleteConfirm:
		return m.viewDeleteConfirm()
	}
	return ""
}

// formatTitle is een gedeelde helper.
func formatTitle(name, badge string) string {
	return fmt.Sprintf("%-30s %s", name, badge)
}
