package tui

import "fmt"

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
	case ScreenRun:
		return m.viewRun()
	case ScreenSettings:
		return m.viewSettings()
	case ScreenDeleteConfirm:
		return m.viewDeleteConfirm()
	case ScreenProjectTypePicker:
		return m.viewProjectTypePicker()
	}
	return ""
}

func formatTitle(name, badge string) string {
	return fmt.Sprintf("%-30s %s", name, badge)
}
