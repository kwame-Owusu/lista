package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/config"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/storage"
	"github.com/kwame-Owusu/lista/internal/tui"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var todoList *models.TodoList
var dataFile string //$HOME/.config/lista, where our json configs live

func loadTodos() {
	path, err := config.DataFilePath()
	if err != nil {
		fmt.Printf("Error resolving config path: %v\n", err)
		os.Exit(1)
	}
	dataFile = path

	// 0755 = rwx for owner, rx for group and others.
	// This is the standard permission set for config directories:
	// - Owner can read/write config files
	// - Others can traverse the directory but not modify its contents
	permissions := 0755

	if err := os.MkdirAll(filepath.Dir(dataFile), os.FileMode(permissions)); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

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
	err := storage.SaveTodos(todoList.Todos, dataFile)
	if err != nil {
		fmt.Printf("Error saving todos: %v\n", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "lista",
	Short: "A minimal todo CLI program",
	Long:  `Lista is a simple and aesthetic CLI app to manage your todos on the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.NewModel(todoList, dataFile)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(addNotesCmd)
	rootCmd.AddCommand(exportCmd)
	loadTodos()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v, using defaults\n", err)
		cfg = &config.Config{Theme: config.DefaultTheme()}
	}

	tui.InitStyles(cfg.Theme)
}
