package tui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6C5CE7")).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3A3A3")).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0, 0, 0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2).
			Margin(1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C5CE7"))

	PageIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FDCB6E")).
			PaddingRight(1)

	DatabaseIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00CEC9")).
				PaddingRight(1)
)
