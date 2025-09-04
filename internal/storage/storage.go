package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kwame-Owusu/todo-cli/internal/models"
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
		return fmt.Errorf("Error enconding todos: %s", err)
	}
	fmt.Printf("Todos saved in %s\n", filename)
	return nil
}

func LoadTodos(filename string) error {
	return nil
}
