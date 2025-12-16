package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/storage"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.todoList.Todos)-1 {
				m.cursor++
			}
		case " ": // Space to toggle completion
			todos := m.todoList.List()
			if len(todos) > 0 && m.cursor < len(todos) {
				selectedTodo := todos[m.cursor]
				selectedID := selectedTodo.ID

				err := m.todoList.Toggle(selectedID)
				if err != nil {
					m.err = err
					return m, nil
				}

				// Save to file
				if err := storage.SaveTodos(m.todoList, m.filename); err != nil {
					m.err = err
				}

				// Re-sync cursor after reorder
				todos = m.todoList.List()
				if idx := findTodoIndexByID(todos, selectedID); idx >= 0 {
					m.cursor = idx
				}
			}

		case "d", "x": // Delete todo
			todos := m.todoList.List()
			if len(todos) > 0 && m.cursor < len(todos) {
				selectedID := todos[m.cursor].ID

				err := m.todoList.Delete(selectedID)
				if err != nil {
					m.err = err
					return m, nil
				}

				// Save to file
				if err := storage.SaveTodos(m.todoList, m.filename); err != nil {
					m.err = err
				}

				// Clamp cursor safely
				todos = m.todoList.List()
				if m.cursor >= len(todos) && m.cursor > 0 {
					m.cursor--
				}
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}
