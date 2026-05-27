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
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return ""
	}

	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal te klein — minimaal %dx%d (nu %dx%d). Vergroot het venster.", minWidth, minHeight, m.width, m.height)
		return styles.Warning.Render(msg)
	}

	header := styles.Header.Width(m.width).Render(appName)
	body := styles.Body.
		Width(m.width).
		Height(m.height - lipgloss.Height(header)).
		Render("")

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}
