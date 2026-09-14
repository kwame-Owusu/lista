package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kwame-Owusu/lista/internal/models"
)

// contentPadding is the horizontal gutter lipgloss applies on each side of the
// content column (see contentStyle).
const contentPadding = 2

func (m model) View() string {
	if m.confirmDelete {
		return m.renderDeleteModal()
	}

	if m.confirmPurge {
		return m.renderPurgeModal()
	}

	if m.showHelp {
		return m.renderHelpOverlay()
	}

	if m.addingTodo {
		return m.renderForm("Add New Todo")
	}

	if m.editingTodo {
		return m.renderForm("Edit Todo")
	}

	var b strings.Builder
	todos := m.todoList.Todos

	b.WriteString(m.renderTitle())
	b.WriteString(m.renderError())
	b.WriteString(m.renderTodos(todos))
	b.WriteString(m.renderSummary(todos))
	b.WriteString(m.renderHelp())

	return contentStyle.Render(b.String())
}

// innerWidth returns the usable width of the content column, clamped between a
// narrow floor and the content max width.
func (m model) innerWidth() int {
	avail := contentMaxWidth - 2*contentPadding
	if m.width > 0 && m.width-2*contentPadding < avail {
		avail = m.width - 2*contentPadding
	}
	if avail < 20 {
		avail = 20
	}
	return avail
}

// truncate shortens a string to w runes, appending an ellipsis when cut.
func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w <= 1 {
		return "…"
	}
	return string(r[:w-1]) + "…"
}

func (m model) renderTitle() string {
	title := wordmarkStyle.Render("✦ LISTA")
	tagline := taglineStyle.Render("your CLI todo manager")

	b := strings.Builder{}
	b.WriteString(title + "  " + tagline + "\n")
	b.WriteString(cursorStyle.Render(strings.Repeat("─", m.innerWidth())) + "\n\n")
	return b.String()
}

func (m model) renderError() string {
	if m.err == nil {
		return ""
	}
	return errorStyle.Render(fmt.Sprintf("⚠ %v", m.err)) + "\n\n"
}

func (m model) renderTodos(todos []models.Todo) string {
	if len(todos) == 0 {
		msg := emptyStateStyle.Render("No todos yet. Add one to get started!")
		hint := helpStyle.Render("press  a  to add · ? for help")
		return lipgloss.NewStyle().Width(m.innerWidth()).Align(lipgloss.Center).
			Render(msg+"\n"+hint) + "\n"
	}

	var b strings.Builder
	for i, todo := range todos {
		b.WriteString(m.renderTodoLine(i, todo) + "\n")
	}
	return b.String()
}

// todoRowWidths returns the column widths for a todo row plus the total inner
// width the row should span.
func (m model) todoRowWidths() (titleCol, prioCol, timeCol, total int) {
	prioCol = 9  // "[Medium]" (8) + a gap
	timeCol = 17 // "added 999d ago" (14) + breathing room

	titleCol = m.innerWidth() - (3 + 2 + prioCol + 2 + timeCol)
	if titleCol > 40 {
		titleCol = 40
	}
	if titleCol < 8 {
		titleCol = 8
	}

	return titleCol, prioCol, timeCol, m.innerWidth()
}

func (m model) renderTodoLine(i int, todo models.Todo) string {
	selected := m.cursor == i
	titleCol, prioCol, timeCol, total := m.todoRowWidths()

	// Cursor marker (2 wide).
	cursor := "  "
	if selected {
		cursor = "\u203a "
	}

	// Checkbox with status color (2 wide).
	checkbox := "○"
	checkboxStyle := checkboxPendingStyle
	if todo.Completed {
		checkbox = "✓"
		checkboxStyle = checkboxDoneStyle
	}

	// Title (titleCol wide).
	title := truncate(todo.Title, titleCol)

	// Note marker (2 wide).
	note := " "
	if len(todo.Notes) > 0 {
		note = "•"
	}

	// Priority badge (prioCol wide) and timestamp (timeCol wide).
	badge := fmt.Sprintf("[%s]", todo.Priority)
	timeAgo := todo.TimeAgo()

	if selected {
		return m.renderSelectedRow(cursor, checkboxStyle, checkbox, title, note, badge, timeAgo, titleCol, prioCol, timeCol, total, todo)
	}

	var b strings.Builder
	b.WriteString(cursor)
	b.WriteString(lipgloss.NewStyle().Width(2).Render(checkboxStyle.Render(checkbox)))

	titleStyle := lipgloss.NewStyle().Foreground(fgMain)
	if todo.Completed {
		titleStyle = lipgloss.NewStyle().Foreground(fgMuted).Strikethrough(true)
	}
	b.WriteString(titleStyle.Render(lipgloss.NewStyle().Width(titleCol).Render(title)))
	b.WriteString(GetPriorityStyle(todo.Priority.String()).Render(lipgloss.NewStyle().Width(prioCol).Render(badge)))
	b.WriteString(noteMarkerStyle.Render(lipgloss.NewStyle().Width(2).Render(note)))
	if timeAgo != "" {
		b.WriteString(timeAgoStyle.Render(lipgloss.NewStyle().Width(timeCol).Render(timeAgo)))
	} else {
		b.WriteString(lipgloss.NewStyle().Width(timeCol).Render(""))
	}
	return b.String()
}

// renderSelectedRow draws a todo row with a full-width highlight spanning the
// content column, in the pending or completed variant.
func (m model) renderSelectedRow(cursor string, checkboxStyle lipgloss.Style, checkbox, title, note, badge, timeAgo string, titleCol, prioCol, timeCol, total int, todo models.Todo) string {
	var row lipgloss.Style
	var titleRow lipgloss.Style
	if todo.Completed {
		row = lipgloss.NewStyle().
			Foreground(fgMuted).
			Background(bgMain)
		titleRow = row.Strikethrough(true)
	} else {
		row = selectedRowStyle
		titleRow = row
	}

	// Dark text reads better on the bright highlight, so keep the status glyphs
	// in their semantic colors but switch neutral text to the dark ink.
	cbStyle := checkboxStyle
	if !todo.Completed {
		cbStyle = lipgloss.NewStyle().Foreground(bgMain).Bold(true)
	}
	badgeStyle := GetPriorityStyle(todo.Priority.String()).Bold(true)
	timeStyle := timeAgoStyle.Foreground(row.GetForeground())

	var b strings.Builder
	b.WriteString(row.Render(cursor))
	b.WriteString(row.Render(lipgloss.NewStyle().Width(2).Render(cbStyle.Render(checkbox))))
	b.WriteString(titleRow.Render(lipgloss.NewStyle().Width(titleCol).Render(title)))
	b.WriteString(row.Render(lipgloss.NewStyle().Width(prioCol).Render(badgeStyle.Render(badge))))
	b.WriteString(row.Render(lipgloss.NewStyle().Width(2).Render(note)))
	if timeAgo != "" {
		b.WriteString(row.Render(lipgloss.NewStyle().Width(timeCol).Render(timeStyle.Render(timeAgo))))
	} else {
		b.WriteString(row.Render(lipgloss.NewStyle().Width(timeCol).Render("")))
	}

	used := lipgloss.Width(b.String())
	if used < total {
		b.WriteString(row.Render(strings.Repeat(" ", total-used)))
	}
	return b.String()
}

func (m model) renderSummary(todos []models.Todo) string {
	if len(todos) == 0 {
		return "\n"
	}

	completed := 0
	for _, t := range todos {
		if t.Completed {
			completed++
		}
	}

	barCol := 28
	if m.innerWidth()-24 < barCol {
		barCol = m.innerWidth() - 24
	}
	if barCol < 3 {
		barCol = 3
	}

	fill := int(float64(barCol)*float64(completed)/float64(len(todos)) + 0.5)
	if fill > barCol {
		fill = barCol
	}

	done := progressStyle.Render(strings.Repeat("█", fill))
	todo := progressEmptyStyle.Render(strings.Repeat("░", barCol-fill))
	bar := lipgloss.NewStyle().Padding(0, 1).Render(done + todo)

	count := fmt.Sprintf("%d of %d complete", completed, len(todos))

	line := lipgloss.JoinHorizontal(lipgloss.Top, bar, count)
	return "\n" + summaryStyle.Render(line) + "\n\n"
}

func (m model) renderHelp() string {
	return helpStyle.Render(
		"↑/k ↓/j navigate · space toggle · a add · e edit · d/x delete · c purge · u undo · ? help · q quit",
	) + "\n"
}

func (m model) renderDeleteModal() string {
	todos := m.todoList.Todos
	var title string
	if idx := findTodoIndexByID(todos, m.deleteID); idx >= 0 {
		title = todos[idx].Title
	}

	body := fmt.Sprintf(
		"Delete %q?\n\n%s",
		title,
		cursorStyle.Render("y: confirm • n / esc: cancel"),
	)
	return m.renderModal("Delete todo", body)
}

func (m model) renderPurgeModal() string {
	body := fmt.Sprintf(
		"Purge all completed todos?\n\n%s",
		cursorStyle.Render("y: confirm • n / esc: cancel"),
	)
	return m.renderModal("Purge completed", body)
}

func (m model) renderModal(headerTitle, body string) string {
	width := lipgloss.Width(body)
	header := m.renderModalHeader(headerTitle, width)

	full := header + "\n" + body
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(full),
	)
}

func (m model) renderModalHeader(headerTitle string, width int) string {
	return modalHeaderStyle.Width(width).Render(headerTitle)
}

func (m model) helpRow(rawKeys, action string, keyWidth int) string {
	return cursorStyle.Render(fmt.Sprintf("%-*s", keyWidth, rawKeys)) + "  " + itemStyle.Render(action)
}

func maxWidth(keys []string) int {
	w := 0
	for _, k := range keys {
		if l := len([]rune(k)); l > w {
			w = l
		}
	}
	return w
}

func (m model) renderHelpOverlay() string {
	keys := []string{
		"↑ / k", "↓ / j", "space", "a", "e", "d / x", "c",
		"y / enter", "n / esc",
		"tab / shift+tab", "← / → / h / l", "enter / ctrl+s", "u", "?", "q / ctrl+c",
	}
	actions := []string{
		"move up", "move down", "toggle complete", "add todo", "edit todo", "delete todo", "purge completed",
		"confirm", "cancel",
		"next / previous field", "change priority", "save form", "undo last action", "show this help", "quit",
	}

	kw := maxWidth(keys)

	var lines []string
	for i, k := range keys {
		lines = append(lines, m.helpRow(k, actions[i], kw))
	}

	content := cursorStyle.Render("KEYBINDINGS") + "\n\n" +
		strings.Join(lines, "\n") + "\n\n" +
		helpStyle.Render("esc / ? to close")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)
}

func (m model) renderPrioritySelector() string {
	var b strings.Builder

	priorityLabel := fieldLabelStyle.Render("Priority:")
	if m.focusedField == fieldPriority {
		priorityLabel = cursorStyle.Render("▸ Priority:")
	}
	b.WriteString(priorityLabel + "\n")

	for _, p := range priorityOptions {
		var style lipgloss.Style
		if p == m.priority {
			if m.focusedField == fieldPriority {
				style = selectedStyle
			} else {
				style = itemStyle.Foreground(fgMain)
			}
		} else {
			style = itemStyle.Foreground(fgMuted)
		}
		b.WriteString("  " + style.Render(p.String()))
		b.WriteString("  ")
	}
	b.WriteString("\n\n")

	return b.String()
}

func (m model) renderForm(formTitle string) string {
	var b strings.Builder

	// Title field
	b.WriteString(m.renderFieldLabel("Title", m.focusedField == fieldTitle))
	b.WriteString(m.renderInput(m.titleInput.View(), m.focusedField == fieldTitle) + "\n\n")

	// Priority field
	b.WriteString(m.renderPrioritySelector())

	// Notes field
	b.WriteString(m.renderFieldLabel("Notes (optional)", m.focusedField == fieldNotes))
	b.WriteString(m.renderInput(m.notesInput.View(), m.focusedField == fieldNotes) + "\n\n")

	// Help text
	helpText := "tab: next field • ←/→/h/l: change priority • enter/ctrl+s: save • esc: cancel"
	b.WriteString(helpStyle.Render(helpText))

	return m.renderModal("✦ "+formTitle, b.String())
}

func (m model) renderFieldLabel(label string, focused bool) string {
	if focused {
		return cursorStyle.Render("▸ "+label) + "\n"
	}
	return fieldLabelStyle.Render("  "+label) + "\n"
}

func (m model) renderInput(view string, focused bool) string {
	if focused {
		return fieldFocusStyle.Render(view)
	}
	return view
}
