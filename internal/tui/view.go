package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kwame-Owusu/lista/internal/models"
)

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
		return m.renderAddForm()
	}

	if m.editingTodo {
		return m.renderEditForm()
	}

	var b strings.Builder
	todos := m.todoList.Todos

	b.WriteString(m.renderTitle())
	b.WriteString(m.renderError())
	b.WriteString(m.renderTodos(todos))
	b.WriteString(m.renderCompletedCount(todos))
	b.WriteString(m.renderHelp())

	return b.String()
}

func (m model) renderTitle() string {
	title := titleStyle.Render(`
	╔════════════════╗
	║     LISTA      ║
	╚════════════════╝
	Your CLI todo manager
`)
	return title + "\n"
}

func (m model) renderError() string {
	if m.err == nil {
		return ""
	}
	return errorStyle.Render(fmt.Sprintf("⚠ Error: %v", m.err)) + "\n\n"
}

func (m model) renderTodos(todos []models.Todo) string {
	if len(todos) == 0 {
		return itemStyle.Render("No todos yet. Add one to get started!") + "\n"
	}

	var b strings.Builder
	for i, todo := range todos {
		b.WriteString(m.renderTodoLine(i, todo) + "\n")
	}
	return b.String()
}

func (m model) renderTodoLine(i int, todo models.Todo) string {
	// Cursor
	cursor := "  "
	if m.cursor == i {
		cursor = cursorStyle.Render("▶ ")
	}

	// Checkbox
	checkbox := "○"
	if todo.Completed {
		checkbox = "✓"
	}

	// Note indicator
	noteIndicator := ""
	if len(todo.Notes) > 0 {
		noteIndicator = "*"
	}

	content := fmt.Sprintf("%s %s [%s] %s", checkbox, todo.Title, todo.Priority, noteIndicator)
	priorityBadge := GetPriorityStyle(todo.Priority.String()).Render(fmt.Sprintf("[%s]", todo.Priority))
	todoTitle := todo.Title

	// Determine style
	if m.cursor == i {
		if todo.Completed {
			return cursor + completedSelectedStyle.Render(content)
		}
		timeAgo := todo.TimeAgo()
		if timeAgo != "" {
			timeAgo = " " + timeAgoStyle.Render(timeAgo)
		}
		return cursor + selectedStyle.Render(content) + timeAgo
	} else if todo.Completed {
		return cursor + completedStyle.Render(content)
	}
	return fmt.Sprintf("%s %s %s %s", cursor, checkbox, itemStyle.Render(todoTitle), priorityBadge)
}

func (m model) renderCompletedCount(todos []models.Todo) string {
	if len(todos) == 0 {
		return ""
	}

	completedCount := 0
	for _, todo := range todos {
		if todo.Completed {
			completedCount++
		}
	}
	countString := fmt.Sprintf("%v of %v complete", completedCount, len(todos))
	return helpStyle.Render(countString)
}

func (m model) renderHelp() string {
	return helpStyle.Render("\n↑/↓: navigate • space: toggle • a: add • u: undo • ?: help • q: quit")
}

func (m model) renderDeleteModal() string {
	todos := m.todoList.Todos
	var title string
	if idx := findTodoIndexByID(todos, m.deleteID); idx >= 0 {
		title = todos[idx].Title
	}

	modal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(
			fmt.Sprintf(
				"Delete \"%s\"?\n\n%s",
				title,
				cursorStyle.Render("y: confirm • n / esc: cancel"),
			),
		),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func (m model) renderPurgeModal() string {
	modal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(
			fmt.Sprintf(
				"Purge all completed todos?\n\n%s",
				cursorStyle.Render("y: confirm • n / esc: cancel"),
			),
		),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
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

	content := titleStyle.Render("KEYBINDINGS") + "\n\n" +
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

	priorityLabel := "Priority:"
	if m.focusedField == fieldPriority {
		priorityLabel = cursorStyle.Render("→ Priority:")
	} else {
		priorityLabel = itemStyle.Render("  Priority:")
	}
	b.WriteString(priorityLabel + "\n")

	for i, p := range priorityOptions {
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
		if i < len(priorityOptions)-1 {
			b.WriteString("  ")
		}
	}
	b.WriteString("\n\n")

	return b.String()
}

func (m model) renderAddForm() string {
	var b strings.Builder

	// Title
	formTitle := titleStyle.Render("✨ Add New Todo") + "\n\n"
	b.WriteString(formTitle)

	// Title field
	titleLabel := "Title:"
	if m.focusedField == fieldTitle {
		titleLabel = cursorStyle.Render("→ Title:")
	} else {
		titleLabel = itemStyle.Render("  Title:")
	}
	b.WriteString(titleLabel + "\n")
	b.WriteString(m.titleInput.View() + "\n\n")

	// Priority field
	b.WriteString(m.renderPrioritySelector())

	// Notes field
	notesLabel := "Notes (optional):"
	if m.focusedField == fieldNotes {
		notesLabel = cursorStyle.Render("→ Notes (optional):")
	} else {
		notesLabel = itemStyle.Render("  Notes (optional):")
	}
	b.WriteString(notesLabel + "\n")
	b.WriteString(m.notesInput.View() + "\n\n")

	// Help text
	helpText := "tab: next field • ←/→/h/l: change priority • enter/ctrl+s: save • esc: cancel"
	b.WriteString(helpStyle.Render(helpText))

	// Center the form
	content := b.String()
	formBox := modalStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formBox,
	)
}

func (m model) renderEditForm() string {
	var b strings.Builder

	// Title
	formTitle := titleStyle.Render("✨ Edit Todo") + "\n\n"
	b.WriteString(formTitle)

	// Title field
	titleLabel := "Title:"
	if m.focusedField == fieldTitle {
		titleLabel = cursorStyle.Render("→ Title:")
	} else {
		titleLabel = itemStyle.Render("  Title:")
	}
	b.WriteString(titleLabel + "\n")
	b.WriteString(m.titleInput.View() + "\n\n")

	// Priority field
	b.WriteString(m.renderPrioritySelector())

	// Notes field
	notesLabel := "Notes (optional):"
	if m.focusedField == fieldNotes {
		notesLabel = cursorStyle.Render("→ Notes (optional):")
	} else {
		notesLabel = itemStyle.Render("  Notes (optional):")
	}
	b.WriteString(notesLabel + "\n")
	b.WriteString(m.notesInput.View() + "\n\n")

	// Help text
	helpText := "tab: next field • ←/→/h/l: change priority • enter/ctrl+s: save • esc: cancel"
	b.WriteString(helpStyle.Render(helpText))

	// Center the form
	content := b.String()
	formBox := modalStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formBox,
	)
}
