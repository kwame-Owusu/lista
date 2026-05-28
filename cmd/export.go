package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kwame-Owusu/lista/internal/storage"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports todos",
	Long:  "Exports todos to stdout in .md syntax",
	Args:  cobra.MinimumNArgs(0),
	Run:   exportTodos,
}

func exportTodos(cmd *cobra.Command, args []string) {
	todos, err := storage.LoadTodos(dataFile)
	if err != nil {
		fmt.Printf("Error reading todos %v", err)
	}

	var sb strings.Builder

	sb.WriteString("# Lista\n\n")
	sb.WriteString("## Todos\n")
	for _, todo := range todos {
		checked := " "
		if todo.Completed {
			checked = "x"
		}
		fmt.Fprintf(&sb, "- [%s] %s\n", checked, todo.Title)
	}

	sb.WriteString("\n")
	sb.WriteString("## Notes\n")
	for _, todo := range todos {
		fmt.Fprintf(&sb, "- %s: %s\n", todo.Title, todo.Notes)
	}

	fmt.Fprint(os.Stdout, sb.String())
}
