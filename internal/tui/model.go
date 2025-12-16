package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/models"
)

type model struct {
	todoList      *models.TodoList
	cursor        int // which todo is selected
	width         int // terminal width
	height        int // terminal height
	filename      string
	err           error
	confirmDelete bool
	deleteID      int
}

func NewModel(todoList *models.TodoList, filename string) model {
	return model{
		todoList: todoList,
		cursor:   0,
		filename: filename,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func findTodoIndexByID(todos []models.Todo, id int) int {
	for i, t := range todos {
		if t.ID == id {
			return i
		}
	}
	return -1
}
