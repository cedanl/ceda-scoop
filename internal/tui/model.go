package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cedanl/ceda-scoop/internal/runner"
)

// repoItem implementeert list.Item voor de Bubble Tea lijst.
type repoItem struct {
	repo        Repo
	installed   bool
	installPath string
}

func (i repoItem) Title() string {
	badge := StyleBadgeAvailable.Render("─ Beschikbaar")
	if i.installed {
		badge = StyleBadgeInstalled.Render("✓ Geïnstalleerd")
	}
	return fmt.Sprintf("%-30s %s", i.repo.Name, badge)
}

func (i repoItem) Description() string { return i.repo.Description }
func (i repoItem) FilterValue() string  { return i.repo.Name }

// Model is de centrale Bubble Tea applicatiestatus.
type Model struct {
	CurrentScreen Screen
	ActiveTab     int // 0 = Store, 1 = Library

	repoItems    []repoItem
	selectedRepo *repoItem
	InstallBase  string

	storeList list.Model
	libList   list.Model
	spinner   spinner.Model
	progress  progress.Model
	logView   viewport.Model

	installLog     []string
	installDone    bool
	installErr     string
	installPercent float64

	settingsInput string
	editingPath   bool

	width  int
	height int
}

// ── Initialisatie ─────────────────────────────────────────────────────────────

func buildRepoItems(base string) []repoItem {
	items := make([]repoItem, len(Catalog))
	for idx, r := range Catalog {
		path := filepath.Join(base, r.ID)
		_, err := os.Stat(path)
		items[idx] = repoItem{repo: r, installed: err == nil, installPath: path}
	}
	return items
}

func toListItems(items []repoItem, onlyInstalled bool) []list.Item {
	var out []list.Item
	for _, it := range items {
		if onlyInstalled && !it.installed {
			continue
		}
		out = append(out, it)
	}
	return out
}

// InitialModel bouwt het startmodel. Geëxporteerd zodat main.go het kan aanroepen.
func InitialModel() Model {
	base := runner.DefaultInstallBase()
	items := buildRepoItems(base)

	storeList := list.New(toListItems(items, false), list.NewDefaultDelegate(), 0, 0)
	storeList.SetShowTitle(false)
	storeList.SetShowHelp(false)
	storeList.SetShowStatusBar(false)
	storeList.SetFilteringEnabled(false)

	libList := list.New(toListItems(items, true), list.NewDefaultDelegate(), 0, 0)
	libList.SetShowTitle(false)
	libList.SetShowHelp(false)
	libList.SetShowStatusBar(false)
	libList.SetFilteringEnabled(false)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(ColorPrimary)

	return Model{
		CurrentScreen: ScreenSplash,
		repoItems:     items,
		InstallBase:   base,
		storeList:     storeList,
		libList:       libList,
		spinner:       sp,
		progress:      progress.New(progress.WithDefaultGradient()),
		logView:       viewport.New(0, 0),
	}
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(1500*time.Millisecond, func(t time.Time) tea.Msg {
			return SplashDoneMsg{}
		}),
	)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentH := m.height - 6
		m.storeList.SetSize(m.width-4, contentH)
		m.libList.SetSize(m.width-4, contentH)
		m.logView.Width = m.width - 4
		m.logView.Height = contentH - 4

	case SplashDoneMsg:
		m.CurrentScreen = ScreenStore

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		cmds = append(cmds, cmd)

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
		if msg.Err == "" && m.selectedRepo != nil {
			m.selectedRepo.installed = true
			m.selectedRepo.installPath = filepath.Join(m.InstallBase, m.selectedRepo.repo.ID)
			m.refreshLists()
		}

	case tea.KeyMsg:
		newM, cmd := m.handleKey(msg)
		return newM, tea.Batch(append(cmds, cmd)...)
	}

	switch m.CurrentScreen {
	case ScreenStore:
		var cmd tea.Cmd
		m.storeList, cmd = m.storeList.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenLibrary:
		var cmd tea.Cmd
		m.libList, cmd = m.libList.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenInstall:
		var cmd tea.Cmd
		m.logView, cmd = m.logView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleKey verwerkt toetsaanslagen per scherm.
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	if (key == "ctrl+c" || key == "q") && m.CurrentScreen != ScreenInstall {
		return m, tea.Quit
	}

	switch m.CurrentScreen {

	case ScreenSplash:
		if key == "enter" || key == " " {
			m.CurrentScreen = ScreenStore
		}

	case ScreenStore:
		switch key {
		case "tab", "l":
			m.ActiveTab = 1
			m.CurrentScreen = ScreenLibrary
		case "s":
			m.CurrentScreen = ScreenSettings
		case "enter":
			if i, ok := m.storeList.SelectedItem().(repoItem); ok {
				it := i
				m.selectedRepo = &it
				m.CurrentScreen = ScreenDetail
			}
		default:
			var cmd tea.Cmd
			m.storeList, cmd = m.storeList.Update(msg)
			return m, cmd
		}

	case ScreenLibrary:
		switch key {
		case "tab", "h":
			m.ActiveTab = 0
			m.CurrentScreen = ScreenStore
		case "s":
			m.CurrentScreen = ScreenSettings
		case "enter":
			if i, ok := m.libList.SelectedItem().(repoItem); ok {
				it := i
				m.selectedRepo = &it
				m.CurrentScreen = ScreenDetail
			}
		default:
			var cmd tea.Cmd
			m.libList, cmd = m.libList.Update(msg)
			return m, cmd
		}

	case ScreenDetail:
		switch key {
		case "esc", "b":
			m.CurrentScreen = m.tabScreen()
			m.selectedRepo = nil
		case "i":
			if m.selectedRepo != nil && !m.selectedRepo.installed {
				m.startInstall()
				return m, m.spinner.Tick
			}
		case "o":
			if m.selectedRepo != nil && m.selectedRepo.installed {
				runner.OpenInExplorer(m.selectedRepo.installPath)
			}
		}

	case ScreenInstall:
		if key == "enter" && m.installDone {
			m.ActiveTab = 1
			m.CurrentScreen = ScreenLibrary
			m.installLog = nil
			m.installDone = false
			m.installPercent = 0
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
				m.refreshLists()
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

func (m *Model) tabScreen() Screen {
	if m.ActiveTab == 1 {
		return ScreenLibrary
	}
	return ScreenStore
}

func (m *Model) refreshLists() {
	m.repoItems = buildRepoItems(m.InstallBase)
	m.storeList.SetItems(toListItems(m.repoItems, false))
	m.libList.SetItems(toListItems(m.repoItems, true))
}

func (m *Model) startInstall() {
	m.CurrentScreen = ScreenInstall
	m.installDone = false
	m.installErr = ""
	m.installLog = []string{}
	m.installPercent = 0
	// TODO: vervang met runner.Clone + runner.Bootstrap
	// Stuur voortgang via InstallStepMsg / InstallProgressMsg / InstallDoneMsg
	go func() { _ = m.selectedRepo }()
}
