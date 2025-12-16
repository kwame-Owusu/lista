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
				err := m.todoList.Toggle(selectedTodo.ID)
				if err != nil {
					m.err = err
					return m, nil
				}
				// Save to file
				err = storage.SaveTodos(m.todoList, m.filename)
				if err != nil {
					m.err = err
				}
			}

		case "d", "x": // Delete todo
			todos := m.todoList.List()
			if len(todos) > 0 && m.cursor < len(todos) {
				selectedTodo := todos[m.cursor]
				err := m.todoList.Delete(selectedTodo.ID)
				if err != nil {
					m.err = err
					return m, nil
				}
				// Adjust cursor if we deleted the last item
				if m.cursor >= len(m.todoList.List()) && m.cursor > 0 {
					m.cursor--
				}
				// Save to file
				err = storage.SaveTodos(m.todoList, m.filename)
				if err != nil {
					m.err = err
				}
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}
