package tui

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	var b strings.Builder

	// Title with style
	title := titleStyle.Render("✨ LISTA - Your CLI task manager")
	b.WriteString(title + "\n\n")

	if m.err != nil {
		errorMsg := errorStyle.
			Render(fmt.Sprintf("⚠ Error: %v", m.err))
		b.WriteString(errorMsg + "\n\n")
	}

	todos := m.todoList.List()

	if len(todos) == 0 {
		empty := itemStyle.Foreground(muted).Render("No todos yet. Add one to get started!")
		b.WriteString(empty + "\n")
	}

	for i, todo := range todos {
		// Cursor indicator
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("▶ ")
		} else {
			cursor = "  "
		}

		// Checkbox
		checkbox := "○"
		if todo.Completed {
			checkbox = "✓"
		}

		content := fmt.Sprintf(
			"%s %s [%s]",
			checkbox,
			todo.Title,
			todo.Priority,
		)

		// Priority badge
		priorityBadge := getPriorityStyle(todo.Priority.String()).
			Render(fmt.Sprintf("[%s]", todo.Priority))

		// Todo title
		todoTitle := todo.Title

		// Build the line
		var line string
		if m.cursor == i {
			// Selected item - highlight entire line
			line = cursor + selectedStyle.Render(content)
			if todo.Completed {
				line = cursor + completedSelectedStyle.Render(content)
			}
		} else if todo.Completed {
			// Completed item
			line = cursor + completedStyle.Render(content)
		} else {
			// Normal item
			line = fmt.Sprintf("%s %s %s %s", cursor, checkbox,
				itemStyle.Render(todoTitle), priorityBadge)
		}

		b.WriteString(line + "\n")
	}

	// Help text
	help := helpStyle.Render("\n↑/↓: navigate • space: toggle • a: add • d: delete • e: edit • q: quit")
	b.WriteString(help)

	return b.String()
}
