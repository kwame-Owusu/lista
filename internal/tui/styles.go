package tui

import "github.com/charmbracelet/lipgloss"

var (
	primary   = lipgloss.Color("#A7AAE1") // purple
	secondary = lipgloss.Color("#FFA4A4") // pink
	success   = lipgloss.Color("#C1E59F") // green
	warning   = lipgloss.Color("#FEEAC9") // amber
	danger    = lipgloss.Color("#FD7980") // red
	muted     = lipgloss.Color("#6B7280") // gray
	text      = lipgloss.Color("#F9FAFB") // light text

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primary).
			Padding(0, 1).
			MarginBottom(1)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3)

	// Selected item (cursor on it)
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(text).
			Background(primary).
			Padding(0, 1)

	completedSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(text).
				Background(muted).
				Padding(0, 1)

	// Normal item
	itemStyle = lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 1)

	// Completed task
	completedStyle = lipgloss.NewStyle().
			Foreground(muted).
			Strikethrough(true).
			Padding(0, 1)

	// Priority styles
	highPriorityStyle = lipgloss.NewStyle().
				Foreground(danger).
				Bold(true)

	mediumPriorityStyle = lipgloss.NewStyle().
				Foreground(warning).
				Bold(true)

	lowPriorityStyle = lipgloss.NewStyle().
				Foreground(success)

	// Help text at bottom
	helpStyle = lipgloss.NewStyle().
			Foreground(muted).
			MarginTop(1)

	// Status indicators
	cursorStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)
)

// Helper function to get priority style
func getPriorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "High":
		return highPriorityStyle
	case "Medium":
		return mediumPriorityStyle
	default:
		return lowPriorityStyle
	}
}
