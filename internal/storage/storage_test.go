package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kwame-Owusu/todo-cli/internal/models"
)

func createTestTodoList() *models.TodoList {
	tl := models.NewTodoList()
	tl.Add("Buy groceries")
	tl.Add("Walk the dog")
	tl.Complete(2)
	return tl
}

func TestSaveTodos_Success(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "todo-test-*")
	defer os.RemoveAll(tempDir)

	tl := createTestTodoList()
	filename := filepath.Join(tempDir, "test.json")

	err := SaveTodos(tl, filename)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check file was created and has content
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("File not created: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Buy groceries") {
		t.Error("File doesn't contain expected content")
	}
}

func TestSaveTodos_InvalidFilename(t *testing.T) {
	tl := createTestTodoList()

	err := SaveTodos(tl, "test.txt")
	if err == nil {
		t.Error("Expected error for invalid filename")
	}

	if !strings.Contains(err.Error(), "Invalid filename") {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestLoadTodos_Success(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "todo-test-*")
	defer os.RemoveAll(tempDir)

	// Save first
	tl := createTestTodoList()
	filename := filepath.Join(tempDir, "test.json")
	SaveTodos(tl, filename)

	// Then load
	todos, err := LoadTodos(filename)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	if todos[0].Text != "Buy groceries" {
		t.Errorf("Expected 'Buy groceries', got '%s'", todos[0].Text)
	}

	if !todos[1].Completed {
		t.Error("Expected second todo to be completed")
	}
}

func TestLoadTodos_InvalidFilename(t *testing.T) {
	_, err := LoadTodos("test.txt")
	if err == nil {
		t.Error("Expected error for invalid filename")
	}
}

func TestLoadTodos_NonExistentFile(t *testing.T) {
	_, err := LoadTodos("missing.json")
	if err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "todo-test-*")
	defer os.RemoveAll(tempDir)

	original := createTestTodoList()
	filename := filepath.Join(tempDir, "roundtrip.json")

	// Save and load
	SaveTodos(original, filename)
	loaded, err := LoadTodos(filename)
	if err != nil {
		t.Fatalf("Round trip failed: %v", err)
	}

	// Compare
	originalTodos := original.List()
	if len(loaded) != len(originalTodos) {
		t.Errorf("Length mismatch: expected %d, got %d", len(originalTodos), len(loaded))
	}

	for i := range originalTodos {
		if originalTodos[i].Text != loaded[i].Text {
			t.Errorf("Text mismatch: expected '%s', got '%s'", originalTodos[i].Text, loaded[i].Text)
		}
	}
}
