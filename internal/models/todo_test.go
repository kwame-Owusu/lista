package models

import (
	"testing"
)

func TestNewTodoList(t *testing.T) {
	tl := NewTodoList()

	if tl == nil {
		t.Fatal("NewTodoList() returned nil")
	}

	if len(tl.Todos) != 0 {
		t.Errorf("Expected empty todo list, got %d todos", len(tl.Todos))
	}

	if tl.NextID != 1 {
		t.Errorf("Expected NextID to be 1, got %d", tl.NextID)
	}
}

func TestTodoList_Add(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		expectError bool
	}{
		{"Valid todo", "Buy groceries", false},
		{"Empty string", "", true},
		{"Only spaces", "   ", true},
		{"Valid with spaces", "  Buy milk  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := NewTodoList()
			err := tl.Add(tt.text, Low)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.text)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(tl.Todos) != 1 {
				t.Errorf("Expected 1 todo, got %d", len(tl.Todos))
			}

			todo := tl.Todos[0]
			if todo.ID != 1 {
				t.Errorf("Expected ID 1, got %d", todo.ID)
			}

			if todo.Title != tt.text {
				t.Errorf("Expected text '%s', got '%s'", tt.text, todo.Title)
			}

			if todo.Completed {
				t.Error("Expected todo to be incomplete")
			}
		})
	}
}

func TestTodoList_Add_Multiple(t *testing.T) {
	tl := NewTodoList()

	tl.Add("First todo", Low)
	tl.Add("Second todo", Medium)
	tl.Add("Third todo", High)

	if len(tl.Todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(tl.Todos))
	}

	if tl.NextID != 4 {
		t.Errorf("Expected NextID to be 4, got %d", tl.NextID)
	}

	// Check IDs are sequential
	expectedIDs := []int{1, 2, 3}
	for i, todo := range tl.Todos {
		if todo.ID != expectedIDs[i] {
			t.Errorf("Expected ID %d, got %d", expectedIDs[i], todo.ID)
		}
	}
}

func TestTodoList_Complete(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Test todo", Low)

	// Test completing existing todo
	err := tl.Complete(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !tl.Todos[0].Completed {
		t.Error("Expected todo to be completed")
	}

	// Test completing non-existent todo
	err = tl.Complete(999)
	if err == nil {
		t.Error("Expected error when completing non-existent todo")
	}
}

func TestTodoList_List(t *testing.T) {
	tl := NewTodoList()
	tl.Add("First todo", Low)
	tl.Add("Second todo", Low)

	todos := tl.List()

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	if todos[0].Title != "First todo" {
		t.Errorf("Expected 'First todo', got '%s'", todos[0].Title)
	}

	if todos[1].Title != "Second todo" {
		t.Errorf("Expected 'Second todo', got '%s'", todos[1].Title)
	}
}

func TestTodoList_GetByID(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Test todo", Low)

	// Test getting existing todo
	todo, err := tl.GetByID(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}

	if todo.Title != "Test todo" {
		t.Errorf("Expected 'Test todo', got '%s'", todo.Title)
	}

	// Test getting non-existent todo
	_, err = tl.GetByID(999)
	if err == nil {
		t.Error("Expected error when getting non-existent todo")
	}

	// Test empty list
	emptyTl := NewTodoList()
	_, err = emptyTl.GetByID(1)
	if err == nil {
		t.Error("Expected error when getting from empty list")
	}
}

func TestTodoList_Delete(t *testing.T) {
	tl := NewTodoList()
	tl.Add("First todo", Low)
	tl.Add("Second todo", Low)
	tl.Add("Third todo", Low)

	// Delete middle item
	err := tl.Delete(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(tl.Todos) != 2 {
		t.Errorf("Expected 2 todos after deletion, got %d", len(tl.Todos))
	}

	// Verify correct items remain
	todos := tl.List()
	if todos[0].ID != 1 || todos[0].Title != "First todo" {
		t.Error("First todo should remain")
	}

	if todos[1].ID != 3 || todos[1].Title != "Third todo" {
		t.Error("Third todo should remain")
	}

	// Test deleting non-existent todo
	err = tl.Delete(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent todo")
	}
}

func TestTodoList_Count(t *testing.T) {
	tl := NewTodoList()

	if tl.Count() != 0 {
		t.Errorf("Expected count 0, got %d", tl.Count())
	}

	tl.Add("First todo", Low)
	if tl.Count() != 1 {
		t.Errorf("Expected count 1, got %d", tl.Count())
	}

	tl.Add("Second todo", Low)
	if tl.Count() != 2 {
		t.Errorf("Expected count 2, got %d", tl.Count())
	}

	tl.Delete(1)
	if tl.Count() != 1 {
		t.Errorf("Expected count 1 after deletion, got %d", tl.Count())
	}
}

func TestTodoList_GetPending(t *testing.T) {
	tl := NewTodoList()

	// Test empty list
	pending := tl.GetPending()
	if len(pending) != 0 {
		t.Errorf("Expected 0 pending todos, got %d", len(pending))
	}

	// Add todos
	tl.Add("Pending todo 1", Low)
	tl.Add("Pending todo 2", Low)
	tl.Add("To be completed", Low)

	// Complete one todo
	tl.Complete(3)

	pending = tl.GetPending()
	if len(pending) != 2 {
		t.Errorf("Expected 2 pending todos, got %d", len(pending))
	}

	// Verify only incomplete todos are returned
	for _, todo := range pending {
		if todo.Completed {
			t.Error("GetPending returned a completed todo")
		}
	}
}

func TestTodoList_CountPending(t *testing.T) {
	tl := NewTodoList()

	// Test empty list
	if tl.CountPending() != 0 {
		t.Errorf("Expected 0 pending todos, got %d", tl.CountPending())
	}

	// Add todos
	tl.Add("Todo 1", Low)
	tl.Add("Todo 2", Low)
	tl.Add("Todo 3", Low)

	if tl.CountPending() != 3 {
		t.Errorf("Expected 3 pending todos, got %d", tl.CountPending())
	}

	// Complete some todos
	tl.Complete(1)
	tl.Complete(2)

	if tl.CountPending() != 1 {
		t.Errorf("Expected 1 pending todo, got %d", tl.CountPending())
	}
}

func TestTodoList_GetCompleted(t *testing.T) {
	tl := NewTodoList()

	// Test empty list
	completed := tl.GetCompleted()
	if len(completed) != 0 {
		t.Errorf("Expected 0 completed todos, got %d", len(completed))
	}

	// Add and complete todos
	tl.Add("Todo 1", Low)
	tl.Add("Todo 2", Low)
	tl.Add("Todo 3", Low)

	tl.Complete(1)
	tl.Complete(3)

	completed = tl.GetCompleted()
	if len(completed) != 2 {
		t.Errorf("Expected 2 completed todos, got %d", len(completed))
	}

	// Verify only completed todos are returned
	for _, todo := range completed {
		if !todo.Completed {
			t.Error("GetCompleted returned an incomplete todo")
		}
	}

	// Check specific IDs
	expectedIDs := map[int]bool{1: true, 3: true}
	for _, todo := range completed {
		if !expectedIDs[todo.ID] {
			t.Errorf("Unexpected completed todo ID: %d", todo.ID)
		}
	}
}

func TestTodoList_Edit(t *testing.T) {
	tests := []struct {
		name        string
		initial     []Todo
		id          int
		newText     string
		expectError bool
		expected    []Todo
	}{
		{
			name: "Edit existing todo",
			initial: []Todo{
				{ID: 1, Title: "Old text", Completed: false},
				{ID: 2, Title: "Another todo", Completed: false},
			},
			id:          1,
			newText:     "Updated text",
			expectError: false,
			expected: []Todo{
				{ID: 1, Title: "Updated text", Completed: false},
				{ID: 2, Title: "Another todo", Completed: false},
			},
		},
		{
			name: "Edit non-existent todo",
			initial: []Todo{
				{ID: 1, Title: "Old text", Completed: false},
			},
			id:          99,
			newText:     "Should not work",
			expectError: true,
			expected: []Todo{
				{ID: 1, Title: "Old text", Completed: false},
			},
		},
		{
			name:        "Edit in empty list",
			initial:     []Todo{},
			id:          1,
			newText:     "No todos here",
			expectError: true,
			expected:    []Todo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := &TodoList{
				Todos:  append([]Todo(nil), tt.initial...), // copy to avoid mutation issues
				NextID: len(tt.initial) + 1,
			}

			err := tl.Edit(tt.id, tt.newText)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error editing ID %d, got nil", tt.id)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if len(tl.Todos) != len(tt.expected) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expected), len(tl.Todos))
			}

			for i := range tl.Todos {
				if tl.Todos[i] != tt.expected[i] {
					t.Errorf("Expected todo %+v, got %+v", tt.expected[i], tl.Todos[i])
				}
			}
		})
	}
}
