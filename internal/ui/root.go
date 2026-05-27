package ui

import (
	"fmt"

	"github.com/cedanl/ceda-scoop/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	appName   = "CEDA Store"
	minWidth  = 80
	minHeight = 24
	maxWidth  = 120
)

type Model struct {
	width  int
	height int
	ready  bool
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
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
		}
	}
	return m, nil
}

// contentWidth returns the capped content width and left margin for centering.
func (m Model) contentWidth() (width, margin int) {
	w := m.width
	if w > maxWidth {
		w = maxWidth
	}
	margin = (m.width - w) / 2
	return w, margin
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
	body := pad.Render(styles.Body.
		Width(cw).
		Height(m.height - lipgloss.Height(header)).
		Render(""))

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}
