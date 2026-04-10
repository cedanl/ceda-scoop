package tui

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cedanl/ceda-scoop/internal/runner"
)

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

	case runTickMsg:
		if !m.runDone {
			m.runElapsed = time.Since(m.runStart)
			if len(m.runSteps) > 0 {
				stepPct := float64(m.runCurrentStep) / float64(len(m.runSteps))
				elapsed := m.runElapsed.Seconds()
				withinStep := (1 - (1 / (1 + elapsed/10))) * (1.0 / float64(len(m.runSteps))) * 0.8
				target := stepPct + withinStep
				cap := float64(m.runCurrentStep+1) / float64(len(m.runSteps))
				if target > cap {
					target = cap
				}
				if target > m.runPercent {
					m.runPercent = target
					cmds = append(cmds, m.progress.SetPercent(target))
				}
			}
			cmds = append(cmds, runTickCmd())
		}

	case RunStepDoneMsg:
		installPath := ""
		if m.selectedRepo != nil {
			installPath = m.selectedRepo.installPath
		}

		if msg.Err != "" {
			// Stap mislukt
			m.runDone = true
			m.runErr = msg.Err
			m.runFailedStep = msg.StepIdx
			m.runElapsed = time.Since(m.runStart)
			cmds = append(cmds, m.progress.SetPercent(float64(msg.StepIdx)/float64(len(m.runSteps))))

		} else {
			// Stap geslaagd

			// Na detect-stap: breidt stappenlijst uit met type-specifieke stappen
			if m.runSteps[msg.StepIdx].StepName == "detect" {
				pt := msg.DetectedType
				m.runProjectType = pt

				var typeSteps []runner.RunStep
				switch pt {
				case runner.ProjectTypeR:
					typeSteps = runner.RSteps
				case runner.ProjectTypeUV:
					typeSteps = runner.UVSteps
				default:
					// Onbekend type — stop
					m.runDone = true
					m.runErr = "Projecttype niet herkend (geen renv.lock of uv.lock gevonden)"
					m.runFailedStep = msg.StepIdx
					m.runElapsed = time.Since(m.runStart)
					cmds = append(cmds, m.progress.SetPercent(float64(msg.StepIdx)/float64(len(m.runSteps))))
					break
				}

				// Voeg type-stappen toe aan de lijst
				m.runSteps = append(m.runSteps, typeSteps...)
			}

			m.runCurrentStep = msg.StepIdx + 1
			pct := float64(m.runCurrentStep) / float64(len(m.runSteps))
			m.runPercent = pct
			cmds = append(cmds, m.progress.SetPercent(pct))

			if m.runCurrentStep >= len(m.runSteps) {
				m.runDone = true
				m.runElapsed = time.Since(m.runStart)
				cmds = append(cmds, m.progress.SetPercent(1.0))
			} else {
				cmds = append(cmds, doRunStepCmd(m.runCurrentStep, m.runSteps[m.runCurrentStep], installPath))
			}
		}

	case DeleteDoneMsg:
		if msg.Err == "" {
			m.refreshItems()
			m.selectedRepo = nil
			m.ActiveTab = 0
			m.CurrentScreen = ScreenStore
		} else {
			m.CurrentScreen = ScreenDetail
		}

	case tea.KeyMsg:
		newM, cmd := m.handleKey(msg)
		return newM, tea.Batch(append(cmds, cmd)...)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	blocking := m.CurrentScreen == ScreenInstall || m.CurrentScreen == ScreenRun
	if (key == "ctrl+c" || key == "q") && !blocking {
		return m, tea.Quit
	}

	switch m.CurrentScreen {

	case ScreenSplash:
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
				return m.beginRun()
			}
		case "f":
			if m.selectedRepo != nil && m.selectedRepo.installed {
				runner.OpenInExplorer(m.selectedRepo.installPath)
			}
		case "d":
			if m.selectedRepo != nil && m.selectedRepo.installed {
				m.CurrentScreen = ScreenDeleteConfirm
			}
		}

	case ScreenRun:
		if key == "enter" && m.runDone {
			m.CurrentScreen = ScreenDetail
			m.runDone = false
			m.runErr = ""
			m.runFailedStep = -1
			m.runCurrentStep = 0
			m.runSteps = nil
			m.runPercent = 0
			m.runElapsed = 0
		}

	case ScreenInstall:
		if key == "enter" && m.installDone {
			m.ActiveTab = 1
			m.CurrentScreen = ScreenLibrary
			m.installDone = false
			m.installErr = ""
			m.installPercent = 0
			m.installElapsed = 0
		}

	case ScreenDeleteConfirm:
		switch key {
		case "y", "enter":
			if m.selectedRepo != nil {
				return m, doDeleteCmd(m.selectedRepo.installPath)
			}
		case "n", "esc", "b":
			m.CurrentScreen = ScreenDetail
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
	m.installPercent = 0
	m.installStart = time.Now()
	m.installElapsed = 0

	return m, tea.Batch(
		m.spinner.Tick,
		installTickCmd(),
		doInstallCmd(repo.RepoURL, destPath, repo.Script),
	)
}

func (m Model) beginRun() (Model, tea.Cmd) {
	if m.selectedRepo == nil {
		return m, nil
	}

	// Start altijd met de common stappen — detect breidt later uit
	steps := append([]runner.RunStep{}, runner.CommonSteps...)

	m.CurrentScreen = ScreenRun
	m.runSteps = steps
	m.runCurrentStep = 0
	m.runDone = false
	m.runErr = ""
	m.runFailedStep = -1
	m.runPercent = 0
	m.runStart = time.Now()
	m.runElapsed = 0
	m.runProjectType = runner.ProjectTypeUnknown

	installPath := m.selectedRepo.installPath

	return m, tea.Batch(
		m.spinner.Tick,
		runTickCmd(),
		doRunStepCmd(0, steps[0], installPath),
	)
}
