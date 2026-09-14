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
	// Refresh relative timestamps on the heartbeat regardless of mode.
	if _, ok := msg.(tickMsg); ok {
		return m, tick()
	}

	var cmd tea.Cmd

	// Handle form input when adding todo
	if m.addingTodo {
		return m.updateAddForm(msg)
	}

	if m.editingTodo {
		return m.updateEditForm(msg) // new update function to perform editing
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showHelp {
			switch msg.String() {
			case "?", "esc":
				m.showHelp = false
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, cmd
		}

		switch msg.String() {
		case "y", "enter":
			if m.confirmDelete {
				cmd = m.deleteTodo()
			} else if m.confirmPurge {
				cmd = m.purgeCompleted()
			}
		case "n", "esc":
			if m.confirmDelete {
				m.cancelDelete()
			} else if m.confirmPurge {
				m.cancelPurge()
			}
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
		case "c":
			m.startPurge()
		case "?":
			if !m.confirmDelete && !m.confirmPurge {
				m.showHelp = true
			}
		case "u":
			if !m.confirmDelete && !m.confirmPurge {
				cmd = m.performUndo()
			}
		case "a":
			m.startAddTodo()
		case "e":
			m.startEditTodo()
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

		case "left", "h":
			if m.focusedField == fieldPriority {
				m.cyclePriority(true) // cycle up
				return m, nil
			}

		case "right", "l":
			if m.focusedField == fieldPriority {
				m.cyclePriority(false) // cycle down
				return m, nil
			}

		case "enter":
			if m.focusedField == fieldNotes {
				// Allow enter in notes field for newlines
				break
			}
			// update todo if not in notes or use ctrl+s
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

func (m model) updateEditForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.cancelEditTodo()
			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab":
			m.cycleFocusedField(msg.String() == "shift+tab")
			return m, nil

		case "left", "h":
			if m.focusedField == fieldPriority {
				m.cyclePriority(true) // cycle up
				return m, nil
			}

		case "right", "l":
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
			return m, m.updateTodo()

		case "ctrl+s":
			return m, m.updateTodo()
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

func (m *model) startAddTodo() {
	m.addingTodo = true
	m.focusedField = fieldTitle
	m.titleInput.SetValue("")
	m.notesInput.SetValue("")
	m.priority = models.Low
	m.titleInput.Focus()
	m.notesInput.Blur()
}

func (m *model) startEditTodo() {
	todos := m.todoList.Todos
	if len(todos) == 0 || m.cursor >= len(todos) {
		return
	}

	todo := todos[m.cursor]

	m.editingTodo = true
	m.editingID = todo.ID
	m.focusedField = fieldTitle

	m.titleInput.SetValue(todo.Title)
	m.notesInput.SetValue(todo.Notes)
	m.priority = todo.Priority

	m.titleInput.Focus()
	m.notesInput.Blur()
}

func (m *model) cancelAddTodo() {
	m.addingTodo = false
	m.titleInput.Blur()
	m.notesInput.Blur()
}

func (m *model) cancelEditTodo() {
	m.editingTodo = false
	m.editingID = 0
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
	idx := priorityIndex(m.priority)
	if up {
		idx = (idx - 1 + len(priorityOptions)) % len(priorityOptions)
	} else {
		idx = (idx + 1) % len(priorityOptions)
	}
	m.priority = priorityAt(idx)
}

func (m *model) submitNewTodo() tea.Cmd {
	title := strings.TrimSpace(m.titleInput.Value())
	if title == "" {
		m.err = errors.New("Empty Title, cannot have empty title.")
		return nil
	}

	notes := strings.TrimSpace(m.notesInput.Value())

	err := m.todoList.Add(title, m.priority, notes)
	if err != nil {
		m.err = err
		return nil
	}

	m.cancelAddTodo()
	return m.saveTodosCmd()
}

func (m *model) updateTodo() tea.Cmd {
	title := strings.TrimSpace(m.titleInput.Value())
	if title == "" {
		m.err = errors.New("Empty Title, cannot have empty title.")
		return nil
	}

	notes := strings.TrimSpace(m.notesInput.Value())

	m.pushUndo()
	err := m.todoList.Update(m.editingID, title, m.priority, notes)
	if err != nil {
		m.err = err
		return nil
	}

	m.cancelEditTodo()
	return m.saveTodosCmd()
}

func (m *model) deleteTodo() tea.Cmd {
	m.pushUndo()
	err := m.todoList.Delete(m.deleteID)
	if err != nil {
		m.err = err
	}
	m.confirmDelete = false
	m.deleteID = 0

	if m.cursor >= len(m.todoList.Todos) && m.cursor > 0 {
		m.cursor--
	}
	return m.saveTodosCmd()
}

func (m *model) toggleTodo() tea.Cmd {
	todos := m.todoList.Todos
	if len(todos) == 0 || m.cursor >= len(todos) {
		return nil
	}
	selectedID := todos[m.cursor].ID
	m.pushUndo()
	if err := m.todoList.Toggle(selectedID); err != nil {
		m.err = err
		return nil
	}

	todos = m.todoList.Todos
	if idx := findTodoIndexByID(todos, selectedID); idx >= 0 {
		m.cursor = idx
	}
	return m.saveTodosCmd()
}

func (m *model) saveTodosCmd() tea.Cmd {
	// The save runs in its own goroutine, so snapshot the list to keep the
	// encoder from racing with event-loop mutations (toggle, edit, delete).
	todos := m.todoList.List()
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
	if m.confirmPurge {
		return
	}
	todos := m.todoList.Todos
	if len(todos) == 0 || m.cursor >= len(todos) {
		return
	}
	m.confirmDelete = true
	m.deleteID = todos[m.cursor].ID
}

func (m *model) startPurge() {
	if m.confirmDelete {
		return
	}

	for _, t := range m.todoList.Todos {
		if t.Completed {
			m.confirmPurge = true
			return
		}
	}

	m.err = errors.New("No completed todos to purge.")
}

func (m *model) purgeCompleted() tea.Cmd {
	kept := m.todoList.Todos[:0]
	for _, t := range m.todoList.Todos {
		if !t.Completed {
			kept = append(kept, t)
		}
	}
	m.todoList.Todos = kept
	m.confirmPurge = false

	if m.cursor >= len(m.todoList.Todos) && m.cursor > 0 {
		m.cursor--
	}
	return m.saveTodosCmd()
}

func (m *model) cancelPurge() {
	m.confirmPurge = false
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
