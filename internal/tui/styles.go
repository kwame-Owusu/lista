package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kwame-Owusu/lista/internal/config"
)

// Semantic colors derived from theme
var (
	bgMain  lipgloss.Color
	bgAlt   lipgloss.Color
	fgMain  lipgloss.Color
	fgMuted lipgloss.Color
	border  lipgloss.Color

	errorCol   lipgloss.Color
	warningCol lipgloss.Color
	successCol lipgloss.Color
	infoCol    lipgloss.Color

	priorityHigh   lipgloss.Color
	priorityMedium lipgloss.Color
	priorityLow    lipgloss.Color

	accentPrimary   lipgloss.Color
	accentSecondary lipgloss.Color
)

// Styles
var (
	titleStyle             lipgloss.Style
	modalStyle             lipgloss.Style
	selectedStyle          lipgloss.Style
	itemStyle              lipgloss.Style
	completedStyle         lipgloss.Style
	completedSelectedStyle lipgloss.Style
	highPriorityStyle      lipgloss.Style
	mediumPriorityStyle    lipgloss.Style
	lowPriorityStyle       lipgloss.Style
	helpStyle              lipgloss.Style
	cursorStyle            lipgloss.Style
	errorStyle             lipgloss.Style
	timeAgoStyle           lipgloss.Style

	// TUI layout styles
	contentStyle         lipgloss.Style
	wordmarkStyle        lipgloss.Style
	taglineStyle         lipgloss.Style
	summaryStyle         lipgloss.Style
	emptyStateStyle      lipgloss.Style
	progressStyle        lipgloss.Style
	progressEmptyStyle   lipgloss.Style
	modalHeaderStyle     lipgloss.Style
	fieldFocusStyle      lipgloss.Style
	fieldLabelStyle      lipgloss.Style
	checkboxDoneStyle    lipgloss.Style
	checkboxPendingStyle lipgloss.Style
	noteMarkerStyle      lipgloss.Style
	selectedRowStyle     lipgloss.Style
)

// contentMaxWidth is the widest the TUI content column may grow to so the app
// reads as a deliberate frame even on very wide terminals.
const contentMaxWidth = 80

func InitStyles(theme config.Theme) {
	//  Color mapping
	bgMain = lipgloss.Color(theme.Background)
	bgAlt = lipgloss.Color(theme.BackgroundAlt)

	fgMain = lipgloss.Color(theme.TextPrimary)
	fgMuted = lipgloss.Color(theme.TextMuted)
	border = lipgloss.Color(theme.Border)

	errorCol = lipgloss.Color(theme.Error)
	warningCol = lipgloss.Color(theme.Warning)
	successCol = lipgloss.Color(theme.Success)
	infoCol = lipgloss.Color(theme.Info)

	priorityHigh = lipgloss.Color(theme.PriorityHigh)
	priorityMedium = lipgloss.Color(theme.PriorityMedium)
	priorityLow = lipgloss.Color(theme.PriorityLow)

	accentPrimary = lipgloss.Color(theme.Accent)
	accentSecondary = lipgloss.Color(theme.AccentSecondary)

	//  Styles

	titleStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true).
		MarginBottom(1)

	modalStyle = lipgloss.NewStyle().
		Foreground(fgMain).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(1, 3)

	itemStyle = lipgloss.NewStyle().
		Foreground(fgMain).
		Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(bgMain).
		Background(accentSecondary).
		Bold(true).
		Padding(0, 1)

	completedStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Strikethrough(true).
		Padding(0, 1)

	completedSelectedStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Background(bgMain).
		Strikethrough(true).
		Padding(0, 1)

	highPriorityStyle = lipgloss.NewStyle().
		Foreground(priorityHigh).
		Bold(true)

	mediumPriorityStyle = lipgloss.NewStyle().
		Foreground(priorityMedium).
		Bold(true)

	lowPriorityStyle = lipgloss.NewStyle().
		Foreground(priorityLow)

	helpStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		MarginTop(1)

	cursorStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(errorCol).
		Bold(true)

	timeAgoStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Italic(true)

	// TUI layout styles
	contentStyle = lipgloss.NewStyle().
		Padding(1, 2)

	wordmarkStyle = lipgloss.NewStyle().
		Foreground(accentPrimary).
		Bold(true)

	taglineStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Italic(true)

	summaryStyle = lipgloss.NewStyle().
		Foreground(fgMuted)

	emptyStateStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Bold(true)

	progressStyle = lipgloss.NewStyle().
		Foreground(successCol)

	progressEmptyStyle = lipgloss.NewStyle().
		Foreground(bgAlt)

	modalHeaderStyle = lipgloss.NewStyle().
		Foreground(bgMain).
		Background(accentPrimary).
		Bold(true)

	fieldFocusStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentPrimary)

	fieldLabelStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Bold(true)

	checkboxDoneStyle = lipgloss.NewStyle().
		Foreground(successCol).
		Bold(true)

	checkboxPendingStyle = lipgloss.NewStyle().
		Foreground(fgMuted)

	noteMarkerStyle = lipgloss.NewStyle().
		Foreground(fgMuted).
		Bold(true)

	selectedRowStyle = lipgloss.NewStyle().
		Foreground(fgMain).
		Background(bgAlt)

	//  CLI styles
	headerStyle = lipgloss.NewStyle().
		Foreground(fgMain).
		Bold(true)

	completedText = lipgloss.NewStyle().
		Foreground(fgMuted).
		Strikethrough(true)

	pendingText = lipgloss.NewStyle().
		Foreground(infoCol).
		Bold(true)

	mutedText = lipgloss.NewStyle().
		Foreground(fgMuted)

	normalText = lipgloss.NewStyle().
		Foreground(fgMain)
}

func GetPriorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "High":
		return highPriorityStyle
	case "Medium":
		return mediumPriorityStyle
	default:
		return lowPriorityStyle
	}
}
