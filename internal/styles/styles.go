package styles

import "github.com/charmbracelet/lipgloss"

// Npuls design tokens — see docs/roadmap/001-v0.2.0.md
const (
	ColorOranje  = "#DD784B"
	ColorBlauw   = "#3D68EC"
	ColorGroen   = "#00AF81"
	ColorGeel    = "#F4D74B"
	ColorDark    = "#111827"
	ColorWit     = "#FFFFFF"
	ColorGray500 = "#6B7280"
	ColorGray700 = "#374151"
)

var (
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorOranje)).
		Background(lipgloss.Color(ColorDark)).
		Padding(0, 1)

	Body = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWit)).
		Background(lipgloss.Color(ColorDark))

	BorderBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOranje))

	Accent = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBlauw))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGroen))

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGeel))

	Subtle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGray500))

	Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDark)).
		Background(lipgloss.Color(ColorBlauw))
)
