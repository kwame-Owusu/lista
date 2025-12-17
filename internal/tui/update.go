package tui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/storage"
)

type msgTodoSaved struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle form input when adding todo
	if m.addingTodo {
		return m.updateAddForm(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			if m.confirmDelete {
				cmd = m.deleteTodo()
			}
		case "n", "esc":
			m.cancelDelete()
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.moveCursorUp()
		case "down", "j":
			m.moveCursorDown()
		case " ":
			cmd = m.toggleTodo()
		case "d", "x":
			m.confirmDeleteAtCursor()
		case "a":
			m.startAddTodo()
		}

	case msgTodoSaved:
		if msg.err != nil {
			m.err = msg.err
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

// updateAddForm handles input when the add form is open
func (m model) updateAddForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.cancelAddTodo()
			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab":
			m.cycleFocusedField(msg.String() == "shift+tab")
			return m, nil

		case "left":
			if m.focusedField == fieldPriority {
				m.cyclePriority(true) // cycle up
				return m, nil
			}

		case "right":
			if m.focusedField == fieldPriority {
				m.cyclePriority(false) // cycle down
				return m, nil
			}

		case "enter":
			if m.focusedField == fieldNotes {
				// Allow enter in notes field for newlines
				break
			}
			// Submit form if not in notes or use ctrl+s
			return m, m.submitNewTodo()

		case "ctrl+s":
			return m, m.submitNewTodo()
		}
	}

	// Update the focused input field
	switch m.focusedField {
	case fieldTitle:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case fieldNotes:
		m.notesInput, cmd = m.notesInput.Update(msg)
	}

	return m, cmd
}

// Helper methods
func (m *model) startAddTodo() {
	m.addingTodo = true
	m.focusedField = fieldTitle
	m.titleInput.SetValue("")
	m.notesInput.SetValue("")
	m.priorityIndex = 0
	m.titleInput.Focus()
	m.notesInput.Blur()
}

func (m *model) cancelAddTodo() {
	m.addingTodo = false
	m.titleInput.Blur()
	m.notesInput.Blur()
}

func (m *model) cycleFocusedField(reverse bool) {
	m.titleInput.Blur()
	m.notesInput.Blur()

	if reverse {
		m.focusedField = (m.focusedField - 1 + 3) % 3
	} else {
		m.focusedField = (m.focusedField + 1) % 3
	}

	switch m.focusedField {
	case fieldTitle:
		m.titleInput.Focus()
	case fieldNotes:
		m.notesInput.Focus()
	}
}

func (m *model) cyclePriority(up bool) {
	if up {
		m.priorityIndex = (m.priorityIndex - 1 + 3) % 3
	} else {
		m.priorityIndex = (m.priorityIndex + 1) % 3
	}
}

func (m *model) submitNewTodo() tea.Cmd {
	title := strings.TrimSpace(m.titleInput.Value())
	if title == "" {
		m.err = errors.New("Empty Title, cannot have empty title.")
		return nil
	}

	notes := strings.TrimSpace(m.notesInput.Value())

	// Convert priorityIndex to Priority
	var priority models.Priority
	switch m.priorityIndex {
	case 0:
		priority = models.Low
	case 1:
		priority = models.Medium
	case 2:
		priority = models.High
	}

	err := m.todoList.Add(title, priority, notes)
	if err != nil {
		m.err = err
		return nil
	}

	m.cancelAddTodo()
	return m.saveTodosCmd()
}

func (m *model) deleteTodo() tea.Cmd {
	err := m.todoList.Delete(m.deleteID)
	if err != nil {
		m.err = err
	}
	m.confirmDelete = false
	m.deleteID = 0

	if m.cursor >= len(m.todoList.List()) && m.cursor > 0 {
		m.cursor--
	}
	return m.saveTodosCmd()
}

func (m *model) toggleTodo() tea.Cmd {
	todos := m.todoList.List()
	if len(todos) == 0 || m.cursor >= len(todos) {
		return nil
	}
	selectedID := todos[m.cursor].ID
	if err := m.todoList.Toggle(selectedID); err != nil {
		m.err = err
		return nil
	}

	todos = m.todoList.List()
	if idx := findTodoIndexByID(todos, selectedID); idx >= 0 {
		m.cursor = idx
	}
	return m.saveTodosCmd()
}

func (m *model) saveTodosCmd() tea.Cmd {
	todos := m.todoList.Todos
	filename := m.filename
	return func() tea.Msg {
		err := storage.SaveTodos(todos, filename)
		return msgTodoSaved{err: err}
	}
}

func (m *model) cancelDelete() {
	m.confirmDelete = false
	m.deleteID = 0
}

func (m *model) confirmDeleteAtCursor() {
	todos := m.todoList.List()
	if len(todos) == 0 || m.cursor >= len(todos) {
		return
	}
	m.confirmDelete = true
	m.deleteID = todos[m.cursor].ID
}

func (m *model) moveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) moveCursorDown() {
	if m.cursor < len(m.todoList.Todos)-1 {
		m.cursor++
	}
}
