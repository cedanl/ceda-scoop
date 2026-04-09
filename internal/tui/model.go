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

type repoItem struct {
	repo        Repo
	installed   bool
	installPath string
}

type Model struct {
	CurrentScreen Screen
	ActiveTab     int

	repoItems    []repoItem
	selectedCard int
	selectedRepo *repoItem
	InstallBase  string

	spinner  spinner.Model
	progress progress.Model
	logView  viewport.Model

	// Install
	installDone    bool
	installErr     string
	installPercent float64
	installStart   time.Time
	installElapsed time.Duration

	// Run — stap-gebaseerd
	runSteps       []runner.RunStep // alle stappen voor dit project
	runCurrentStep int              // welke stap loopt nu (0-indexed)
	runDone        bool
	runErr         string          // foutmelding van mislukte stap
	runFailedStep  int             // index van mislukte stap (-1 = geen)
	runStart       time.Time
	runElapsed     time.Duration
	runPercent     float64
	runProjectType runner.ProjectType

	// Project type picker (bij ambiguïteit)
	pickerSelected int // 0 = R, 1 = uv

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
		runFailedStep: -1,
	}
}

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
	items := m.filteredItems()
	if m.selectedCard >= len(items) {
		m.selectedCard = len(items) - 1
	}
	if m.selectedCard < 0 {
		m.selectedCard = 0
	}
}

// ── Tick cmds ─────────────────────────────────────────────────────────────────

func installTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return installTickMsg{}
	})
}

func runTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return runTickMsg{}
	})
}

// ── Install cmd ───────────────────────────────────────────────────────────────

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

// ── Delete cmd ────────────────────────────────────────────────────────────────

func doDeleteCmd(installPath string) tea.Cmd {
	return func() tea.Msg {
		if err := os.RemoveAll(installPath); err != nil {
			return DeleteDoneMsg{Err: err.Error()}
		}
		return DeleteDoneMsg{Err: ""}
	}
}

// ── Run cmds — één stap per Cmd ───────────────────────────────────────────────

// doRunStepCmd voert één stap uit en stuurt RunStepDoneMsg terug.
func doRunStepCmd(stepIdx int, step runner.RunStep) tea.Cmd {
	return func() tea.Msg {
		err := runner.ExecuteStep(step)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		return RunStepDoneMsg{StepIdx: stepIdx, Err: errStr}
	}
}
