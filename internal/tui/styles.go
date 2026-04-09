package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary    = lipgloss.Color("#5C7AEA")
	ColorPrimaryDim = lipgloss.Color("#3D5BD9")
	ColorGreen      = lipgloss.Color("#4CAF50")
	ColorGray       = lipgloss.Color("#6B7280")
	ColorDim        = lipgloss.Color("#374151")
	ColorSurface    = lipgloss.Color("#1E2433")
	ColorWhite      = lipgloss.Color("#F9FAFB")
	ColorRed        = lipgloss.Color("#EF4444")
)

var (
	// Navbar
	StyleNavbar = lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorWhite)

	StyleNavBrand = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPrimary).
			Padding(0, 3)

	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPrimary).
			Padding(0, 3)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorGray).
				Background(ColorSurface).
				Padding(0, 3)

	StyleNavDivider = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(ColorSurface)

	// Content
	StyleTitle  = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(ColorPrimary).Padding(0, 2)
	StyleSubtle = lipgloss.NewStyle().Foreground(ColorGray)

	StyleBadgeInstalled = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)
	StyleBadgeAvailable = lipgloss.NewStyle().Foreground(ColorGray)

	StyleButton = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 2)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(1)

	StyleBold   = lipgloss.NewStyle().Bold(true)
	StyleDanger = lipgloss.NewStyle().Bold(true).Foreground(ColorRed)
)
