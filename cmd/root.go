package cmd

import (
	"fmt"

	"github.com/kwame-Owusu/todo-cli/internal/models"
	"github.com/kwame-Owusu/todo-cli/internal/storage"
	"github.com/spf13/cobra"
)

var todoList *models.TodoList
var dataFile string = "todos.json" //will change to $HOME path and make it not se with .file.json

func loadTodos() {
	todos, err := storage.LoadTodos(dataFile)
	if err != nil {
		// File doesn't exist or error reading - create new TodoList
		todoList = models.NewTodoList()
		return
	}

	// File exists - create TodoList and populate it
	todoList = models.NewTodoList()
	for _, todo := range todos {
		todoList.Todos = append(todoList.Todos, todo)
		// Update NextID to be higher than highest existing ID
		if todo.ID >= todoList.NextID {
			todoList.NextID = todo.ID + 1
		}
	}
}

func saveTodos() {
	err := storage.SaveTodos(todoList, dataFile)
	if err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "lista",
	Short: "A minimal todo CLI in Go",
	Long:  `lista is a simple and aesthetic CLI app to manage your todos on the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to todo-cli! Use 'todo-cli help' to see available commands.")
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		// Exit with error code 1 if something goes wrong
		// os.Exit(1)  // optional
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	loadTodos()
}
