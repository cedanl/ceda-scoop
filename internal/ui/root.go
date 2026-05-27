package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/cedanl/ceda-scoop/internal"
	"github.com/cedanl/ceda-scoop/internal/styles"
	"github.com/cedanl/ceda-scoop/internal/updater"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	appName   = "CEDA Store"
	minWidth  = 80
	minHeight = 24
	maxWidth  = 120
)

type updateState int

const (
	updateChecking updateState = iota
	updateCurrent
	updateAvailable
	updateDownloading
	updateReady    // Unix: binary replaced, restart to apply
	updateWindows  // Windows: *-update.exe ready, manual restart needed
	updateError
)

type Model struct {
	width       int
	height      int
	ready       bool
	showHelp    bool
	version     string
	latestVer   string
	updateState updateState
	updatePath  string
}

func New(version string) Model {
	return Model{
		version:     version,
		updateState: updateChecking,
	}
}

func (m Model) Init() tea.Cmd {
	return checkUpdateCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "u":
			if m.updateState == updateAvailable {
				m.updateState = updateDownloading
				return m, downloadCmd(m.latestVer)
			}
			if m.updateState == updateReady {
				return m, restartCmd(m.updatePath)
			}
		}

	case internal.UpdateCheckDoneMsg:
		if msg.Err != nil || msg.Latest == "" {
			m.updateState = updateError
			return m, nil
		}
		m.latestVer = msg.Latest
		if msg.Latest == m.version {
			m.updateState = updateCurrent
		} else {
			m.updateState = updateAvailable
		}

	case internal.UpdateDownloadDoneMsg:
		if msg.Err != nil {
			m.updateState = updateError
			return m, nil
		}
		m.updatePath = msg.Path
		if runtime.GOOS == "windows" {
			m.updateState = updateWindows
		} else {
			m.updateState = updateReady
		}
	}

	return m, nil
}

func (m Model) contentWidth() (width, margin int) {
	w := m.width
	if w > maxWidth {
		w = maxWidth
	}
	return w, (m.width - w) / 2
}

func (m Model) View() string {
	if !m.ready {
		return ""
	}

	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal te klein — minimaal %dx%d (nu %dx%d). Vergroot het venster.", minWidth, minHeight, m.width, m.height)
		return styles.Warning.Render(msg)
	}

	cw, margin := m.contentWidth()
	pad := lipgloss.NewStyle().PaddingLeft(margin)

	header := pad.Render(styles.Header.Width(cw).Render(appName))
	statusBar := pad.Render(m.renderStatusBar(cw))
	bodyHeight := m.height - lipgloss.Height(header) - lipgloss.Height(statusBar)
	var bodyContent string
	if m.showHelp {
		bodyContent = m.renderHelp()
	}
	body := pad.Render(styles.Body.Width(cw).Height(bodyHeight).Render(bodyContent))

	return lipgloss.JoinVertical(lipgloss.Left, header, body, statusBar)
}

func (m Model) renderStatusBar(width int) string {
	left := styles.Subtle.Render(fmt.Sprintf(" v%s", m.version))

	var right string
	switch m.updateState {
	case updateChecking:
		right = styles.Subtle.Render("update controleren…  ")
	case updateCurrent:
		right = styles.Success.Render("up to date  ")
	case updateAvailable:
		right = styles.Warning.Render(fmt.Sprintf("v%s beschikbaar — druk u om te downloaden  ", m.latestVer))
	case updateDownloading:
		right = styles.Accent.Render("downloaden…  ")
	case updateReady:
		right = styles.Success.Render("gedownload — druk u om te herstarten  ")
	case updateWindows:
		right = styles.Warning.Render("update klaar — herstart om toe te passen  ")
	case updateError:
		right = styles.Subtle.Render("update onbekend  ")
	}

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	bar := left + strings.Repeat(" ", gap) + right
	return styles.Body.Width(width).Render(bar)
}

func (m Model) renderHelp() string {
	bindings := []struct{ key, desc string }{
		{"q / ctrl+c", "afsluiten"},
		{"?", "help tonen/verbergen"},
		{"u", "update downloaden / herstarten"},
	}
	var lines []string
	for _, b := range bindings {
		key := styles.Accent.Render(fmt.Sprintf("  %-14s", b.key))
		lines = append(lines, key+styles.Subtle.Render(b.desc))
	}
	box := styles.BorderBox.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return "\n" + box
}

// checkUpdateCmd runs the GitHub API check off the main goroutine.
func checkUpdateCmd() tea.Msg {
	latest, err := updater.CheckLatest()
	return internal.UpdateCheckDoneMsg{Latest: latest, Err: err}
}

func downloadCmd(version string) tea.Cmd {
	return func() tea.Msg {
		path, err := updater.Download(version)
		return internal.UpdateDownloadDoneMsg{Path: path, Err: err}
	}
}

// restartCmd replaces the current process with a fresh copy (Unix only).
func restartCmd(exe string) tea.Cmd {
	return func() tea.Msg {
		_ = syscall.Exec(exe, os.Args, os.Environ())
		return nil
	}
}
