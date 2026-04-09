package tui

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cedanl/ceda-scoop/internal/runner"
)

// repoItem beschrijft een catalog-item met installatiestatus.
type repoItem struct {
	repo        Repo
	installed   bool
	installPath string
}

// Model is de centrale Bubble Tea applicatiestatus.
type Model struct {
	CurrentScreen Screen
	ActiveTab     int // 0 = Store, 1 = Library

	repoItems    []repoItem
	selectedCard int // index in de gefilterde lijst van de huidige tab
	selectedRepo *repoItem
	InstallBase  string

	spinner  spinner.Model
	progress progress.Model
	logView  viewport.Model

	installLog     []string
	installDone    bool
	installErr     string
	installPercent float64
	installStart   time.Time
	installElapsed time.Duration

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

// filteredItems geeft de items terug die zichtbaar zijn in de huidige tab.
func (m Model) filteredItems() []repoItem {
	if m.ActiveTab == 0 {
		return m.repoItems
	}
	var out []repoItem
	for _, it := range m.repoItems {
		if it.installed {
			out = append(out, it)
		}
	}
	return out
}

// InitialModel bouwt het startmodel.
func InitialModel() Model {
	base := runner.DefaultInstallBase()
	items := buildRepoItems(base)

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(ColorPrimary)

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithoutPercentage(),
	)

	return Model{
		CurrentScreen: ScreenSplash,
		repoItems:     items,
		InstallBase:   base,
		spinner:       sp,
		progress:      prog,
		logView:       viewport.New(0, 0),
	}
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (m *Model) tabScreen() Screen {
	if m.ActiveTab == 1 {
		return ScreenLibrary
	}
	return ScreenStore
}

func (m *Model) refreshItems() {
	m.repoItems = buildRepoItems(m.InstallBase)
	// Clamp selectedCard
	items := m.filteredItems()
	if m.selectedCard >= len(items) {
		m.selectedCard = len(items) - 1
	}
	if m.selectedCard < 0 {
		m.selectedCard = 0
	}
}

// installTickCmd stuurt elke 100ms een tick voor de elapsed timer + progress.
func installTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return installTickMsg{}
	})
}

// doInstallCmd voert clone + bootstrap uit als een echte tea.Cmd.
func doInstallCmd(repoURL, destPath, script string) tea.Cmd {
	return func() tea.Msg {
		if err := runner.Clone(repoURL, destPath); err != nil {
			return InstallDoneMsg{Err: "Clone mislukt: " + err.Error()}
		}
		scriptPath := filepath.Join(destPath, script)
		if _, err := os.Stat(scriptPath); err == nil {
			if err := runner.Bootstrap(destPath, script, make(chan<- string)); err != nil {
				return InstallDoneMsg{Err: "Bootstrap mislukt: " + err.Error()}
			}
		}
		return InstallDoneMsg{Err: ""}
	}
}

// doDeleteCmd verwijdert de installmap van een repo.
func doDeleteCmd(installPath string) tea.Cmd {
	return func() tea.Msg {
		if err := os.RemoveAll(installPath); err != nil {
			return DeleteDoneMsg{Err: err.Error()}
		}
		return DeleteDoneMsg{Err: ""}
	}
}
