package models

import (
	"fmt"
	"strings"
)

type Todo struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Notes     string   `json:"notes,omitempty"`
	Completed bool     `json:"completed"`
	Priority  Priority `json:"priority"`
}

type TodoList struct {
	Todos  []Todo `json:"todos"`
	NextID int    `json:"nextID"`
}

func NewTodoList() *TodoList {
	return &TodoList{
		Todos:  make([]Todo, 0),
		NextID: 1,
	}
}

func (tl *TodoList) Add(title string, priority Priority, notes string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("todo title cannot be empty")
	}
	todo := Todo{
		ID:        tl.NextID,
		Title:     title,
		Notes:     notes,
		Completed: false,
		Priority:  priority,
	}
	tl.Todos = append(tl.Todos, todo)
	tl.NextID++
	return nil
}

func (tl *TodoList) Complete(id int) error {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (tl *TodoList) List() []Todo {
	return tl.Todos
}

func (tl *TodoList) GetByID(id int) (*Todo, error) {
	if len(tl.Todos) < 1 {
		return nil, fmt.Errorf("No entries in todo list")
	}

	for _, todo := range tl.Todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return &Todo{}, fmt.Errorf("Todo item with id %d not found", id)
}

func (tl *TodoList) Delete(id int) error {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (tl *TodoList) Edit(id int, title string) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos[i].Title = title
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (tl *TodoList) Count() int {
	return len(tl.Todos)
}

func (tl *TodoList) GetPending() []Todo {
	if tl.Count() < 1 {
		return []Todo{}
	}
	result := []Todo{}

	for _, todo := range tl.Todos {
		if !todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}

func (tl *TodoList) CountPending() int {
	if tl.Count() < 1 {
		return 0
	}

	count := 0
	for _, todo := range tl.Todos {
		if !todo.Completed {
			count++
		}
	}
	return count
}

func (tl *TodoList) GetCompleted() []Todo {
	if tl.Count() < 1 {
		return []Todo{}
	}
	result := []Todo{}

	for _, todo := range tl.Todos {
		if todo.Completed {
			result = append(result, todo)
		}
	}
	return result
}

// Toggle completion status
func (tl *TodoList) Toggle(id int) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos[i].Completed = !tl.Todos[i].Completed
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}
