package tui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cedanl/ceda-scoop/internal/runner"
)

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.logView.Width = m.width - 4
		m.logView.Height = m.height - 10

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		cmds = append(cmds, cmd)

	case installTickMsg:
		if !m.installDone {
			m.installElapsed = time.Since(m.installStart)
			elapsed := m.installElapsed.Seconds()
			target := 1 - (1 / (1 + elapsed/15))
			if target > 0.90 {
				target = 0.90
			}
			if target > m.installPercent {
				m.installPercent = target
				cmds = append(cmds, m.progress.SetPercent(target))
			}
			cmds = append(cmds, installTickCmd())
		}

	case InstallStepMsg:
		m.installLog = append(m.installLog, msg.Line)
		m.logView.SetContent(strings.Join(m.installLog, "\n"))
		m.logView.GotoBottom()

	case InstallProgressMsg:
		m.installPercent = msg.Pct
		cmds = append(cmds, m.progress.SetPercent(msg.Pct))

	case InstallDoneMsg:
		m.installDone = true
		m.installErr = msg.Err
		m.installElapsed = time.Since(m.installStart)
		if msg.Err == "" && m.selectedRepo != nil {
			m.selectedRepo.installed = true
			m.selectedRepo.installPath = filepath.Join(m.InstallBase, m.selectedRepo.repo.ID)
			m.refreshItems()
		}
		cmds = append(cmds, m.progress.SetPercent(1.0))

	case DeleteDoneMsg:
		if msg.Err == "" {
			m.refreshItems()
			m.selectedRepo = nil
			m.ActiveTab = 0
			m.CurrentScreen = ScreenStore
		} else {
			// Fout tonen — ga terug naar detail
			m.CurrentScreen = ScreenDetail
		}

	case tea.KeyMsg:
		newM, cmd := m.handleKey(msg)
		return newM, tea.Batch(append(cmds, cmd)...)
	}

	return m, tea.Batch(cmds...)
}

// ── handleKey ─────────────────────────────────────────────────────────────────

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	if (key == "ctrl+c" || key == "q") && m.CurrentScreen != ScreenInstall {
		return m, tea.Quit
	}

	switch m.CurrentScreen {

	case ScreenSplash:
		// Wacht op Enter — geen auto-timer
		if key == "enter" || key == " " {
			m.CurrentScreen = ScreenStore
		}

	case ScreenStore, ScreenLibrary:
		library := m.CurrentScreen == ScreenLibrary
		items := m.filteredItems()
		n := len(items)

		switch key {
		case "tab":
			if library {
				m.ActiveTab = 0
				m.CurrentScreen = ScreenStore
			} else {
				m.ActiveTab = 1
				m.CurrentScreen = ScreenLibrary
			}
			m.selectedCard = 0
		case "s":
			m.CurrentScreen = ScreenSettings
		case "up", "k":
			if m.selectedCard > 0 {
				m.selectedCard--
			}
		case "down", "j":
			if n > 0 && m.selectedCard < n-1 {
				m.selectedCard++
			}
		case "enter":
			if n > 0 && m.selectedCard < n {
				it := items[m.selectedCard]
				m.selectedRepo = &it
				m.CurrentScreen = ScreenDetail
			}
		}

	case ScreenDetail:
		switch key {
		case "esc", "b":
			m.CurrentScreen = m.tabScreen()
			m.selectedRepo = nil
		case "i":
			if m.selectedRepo != nil && !m.selectedRepo.installed {
				return m.beginInstall()
			}
		case "o":
			if m.selectedRepo != nil && m.selectedRepo.installed {
				runner.OpenInExplorer(m.selectedRepo.installPath)
			}
		case "d":
			if m.selectedRepo != nil && m.selectedRepo.installed {
				m.CurrentScreen = ScreenDeleteConfirm
			}
		}

	case ScreenDeleteConfirm:
		switch key {
		case "y", "enter":
			if m.selectedRepo != nil {
				path := m.selectedRepo.installPath
				return m, doDeleteCmd(path)
			}
		case "n", "esc", "b":
			m.CurrentScreen = ScreenDetail
		}

	case ScreenInstall:
		if key == "enter" && m.installDone {
			m.ActiveTab = 1
			m.CurrentScreen = ScreenLibrary
			m.installLog = nil
			m.installDone = false
			m.installPercent = 0
			m.installElapsed = 0
		}

	case ScreenSettings:
		switch key {
		case "esc", "b":
			m.CurrentScreen = m.tabScreen()
			m.editingPath = false
		case "enter":
			if m.editingPath && m.settingsInput != "" {
				m.InstallBase = m.settingsInput
				m.editingPath = false
				m.refreshItems()
			} else {
				m.editingPath = true
			}
		case "backspace":
			if m.editingPath && len(m.settingsInput) > 0 {
				m.settingsInput = m.settingsInput[:len(m.settingsInput)-1]
			}
		default:
			if m.editingPath && len(key) == 1 {
				m.settingsInput += key
			}
		}
	}

	return m, nil
}

// beginInstall initialiseert de installatiestatus en geeft de juiste Cmds terug.
func (m Model) beginInstall() (Model, tea.Cmd) {
	if m.selectedRepo == nil {
		return m, nil
	}

	repo := m.selectedRepo.repo
	destPath := filepath.Join(m.InstallBase, repo.ID)
	_ = os.MkdirAll(m.InstallBase, 0755)

	m.CurrentScreen = ScreenInstall
	m.installDone = false
	m.installErr = ""
	m.installLog = []string{}
	m.installPercent = 0
	m.installStart = time.Now()
	m.installElapsed = 0

	return m, tea.Batch(
		m.spinner.Tick,
		installTickCmd(),
		doInstallCmd(repo.RepoURL, destPath, repo.Script),
	)
}
