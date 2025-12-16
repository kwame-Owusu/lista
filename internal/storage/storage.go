package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kwame-Owusu/lista/internal/models"
)

func SaveTodos(todos *models.TodoList, filename string) error {
	if !strings.HasSuffix(filename, ".json") {
		return fmt.Errorf("Invalid filename, file has to end in .json")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error occurred creating file: %s", err)
	}
	defer file.Close()

	// write to json
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ") //make json pretty
	err = encoder.Encode(todos.Todos)
	if err != nil {
		return fmt.Errorf("Error encoding todos: %s", err)
	}
	return nil
}

func LoadTodos(filename string) ([]models.Todo, error) {
	if !strings.HasSuffix(filename, ".json") {
		return nil, fmt.Errorf("Invalid filename, file has to end in .json")
	}

	var todos []models.Todo
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file, %s", err)
	}
	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling file, %s", err)
	}

	return todos, nil
}
