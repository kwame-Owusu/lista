package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes,omitempty"`
	Completed bool      `json:"completed"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Todo) TimeAgo() string {
	if t.CreatedAt.IsZero() {
		return ""
	}
	d := time.Since(t.CreatedAt)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("added %ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("added %dmins ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("added %dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("added %dd ago", int(d.Hours()/24))
	}
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
		CreatedAt: time.Now(),
	}
	tl.Todos = append(tl.Todos, todo)
	tl.NextID++
	return nil
}

func (tl *TodoList) Update(id int, title string, priority Priority, notes string) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos[i].Title = title
			tl.Todos[i].Priority = priority
			tl.Todos[i].Notes = notes
			return nil
		}
	}

	return errors.New("todo not found")
}

func (tl *TodoList) Complete(id int) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			tl.Todos[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("todo with ID %d not found", id)
}

func (tl *TodoList) List() []Todo {
	result := make([]Todo, len(tl.Todos))
	copy(result, tl.Todos)
	return result
}

func (tl *TodoList) GetByID(id int) (*Todo, error) {
	if len(tl.Todos) == 0 {
		return nil, fmt.Errorf("no entries in todo list")
	}

	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			return &tl.Todos[i], nil
		}
	}

	return nil, fmt.Errorf("todo item with id %d not found", id)
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
