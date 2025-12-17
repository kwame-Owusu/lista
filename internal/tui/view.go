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

	if m.addingTodo {
		return m.renderAddForm()
	}

	var b strings.Builder

	b.WriteString(m.renderTitle())
	b.WriteString(m.renderError())
	b.WriteString(m.renderTodos())
	b.WriteString(m.renderCompletedCount())
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

func (m model) renderAddTodo() string {
	return m.textarea.View()
}

func (m model) renderTodos() string {
	todos := m.todoList.List()
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
		noteIndicator = "󰈙"
	}

	content := fmt.Sprintf("%s %s [%s] %s", checkbox, todo.Title, todo.Priority, noteIndicator)
	priorityBadge := getPriorityStyle(todo.Priority.String()).Render(fmt.Sprintf("[%s]", todo.Priority))
	todoTitle := todo.Title

	// Determine style
	if m.cursor == i {
		if todo.Completed {
			return cursor + completedSelectedStyle.Render(content)
		}
		return cursor + selectedStyle.Render(content)
	} else if todo.Completed {
		return cursor + completedStyle.Render(content)
	}
	return fmt.Sprintf("%s %s %s %s", cursor, checkbox, itemStyle.Render(todoTitle), priorityBadge)
}

func (m model) renderCompletedCount() string {
	todos := m.todoList.List()
	if len(todos) == 0 {
		return ""
	}

	completedCount := len(m.todoList.GetCompleted())
	countString := fmt.Sprintf("%v of %v complete", completedCount, len(todos))
	return helpStyle.Render(countString)
}

func (m model) renderHelp() string {
	return helpStyle.Render("\n↑/↓: navigate • space: toggle • a: add • d: delete • e: edit • q: quit")
}

func (m model) renderDeleteModal() string {
	todos := m.todoList.List()
	var title string
	for _, t := range todos {
		if t.ID == m.deleteID {
			title = t.Title
			break
		}
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
	priorities := []string{"Low", "Medium", "High"}
	priorityLabel := "Priority:"
	if m.focusedField == fieldPriority {
		priorityLabel = cursorStyle.Render("→ Priority:")
	} else {
		priorityLabel = itemStyle.Render("  Priority:")
	}
	b.WriteString(priorityLabel + "\n")

	// Render priority options
	for i, p := range priorities {
		var style lipgloss.Style
		if i == m.priorityIndex {
			if m.focusedField == fieldPriority {
				style = selectedStyle
			} else {
				style = itemStyle.Foreground(brightYellow)
			}
		} else {
			style = itemStyle.Foreground(fgMuted)
		}
		b.WriteString("  " + style.Render(p))
		if i < len(priorities)-1 {
			b.WriteString("  ")
		}
	}
	b.WriteString("\n\n")

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
	helpText := "tab: next field • ↑/↓: change priority • enter/ctrl+s: save • esc: cancel"
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
