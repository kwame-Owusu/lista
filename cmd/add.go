package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/spf13/cobra"
)

var priorityFlag string
var notesFlag string

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new todo",
	Long:  "Add a new todo with description and optional priority (high, medium, low) and notes",
	Args:  validateAddArgs,
	RunE:  addTodo,
}

func validateAddArgs(_ *cobra.Command, args []string) error {
	if len(args) > 0 {
		return nil
	}
	stat, err := os.Stdin.Stat()
	if err == nil && stat.Mode()&os.ModeCharDevice == 0 {
		return nil
	}
	return fmt.Errorf("requires a title argument or piped input, e.g. echo \"Buy milk\" | lista add")
}

func init() {
	addCmd.Flags().StringVarP(&priorityFlag, "priority", "p", "low", "Priority level (high/h, medium/m, low/l)")
	addCmd.Flags().StringVarP(&notesFlag, "notes", "n", "", "notes (lorem ipsum)")
}

func addTodo(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")
	if title == "" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		title = strings.TrimSpace(string(input))
	}
	// Parse the priority flag
	priority, err := models.ParsePriority(priorityFlag)
	if err != nil {
		return fmt.Errorf("adding todo: %w", err)
	}
	notes := notesFlag

	if err := todoList.Add(title, priority, notes); err != nil {
		return fmt.Errorf("adding todo: %w", err)
	}

	if err := saveTodos(); err != nil {
		return err
	}
	fmt.Printf("Added todo with ID: %d\n", todoList.NextID-1)
	return nil
}
