package theme

import "github.com/charmbracelet/lipgloss"

var (
	Color = struct {
		Primary    lipgloss.Color
		Secondary  lipgloss.Color
		Accent     lipgloss.Color
		Error      lipgloss.Color
		Success    lipgloss.Color
		Warning    lipgloss.Color
		Background lipgloss.Color
		Foreground lipgloss.Color
	}{
		Primary:    lipgloss.Color("212"),
		Secondary:  lipgloss.Color("243"),
		Accent:     lipgloss.Color("86"),
		Error:      lipgloss.Color("196"),
		Success:    lipgloss.Color("46"),
		Warning:    lipgloss.Color("226"),
		Background: lipgloss.Color("235"),
		Foreground: lipgloss.Color("255"),
	}

	TabActiveStyle = lipgloss.NewStyle().
		Foreground(Color.Primary).
		Bold(true).
		Border(lipgloss.RoundedBorder(), false, false, true, false).
		BorderForeground(Color.Primary).
		Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
		Foreground(Color.Secondary).
		Padding(0, 2)

	PaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Color.Secondary).
		Padding(1, 2).
		Background(Color.Background)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(Color.Primary).
		Bold(true).
		MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(Color.Secondary).
		Italic(true)
)
