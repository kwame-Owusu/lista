package cmd

import (
	"fmt"
	"github.com/kwame-Owusu/lista/internal/tui"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos",
	Long:  "List all todo tasks in the todo list",
	Args:  cobra.MinimumNArgs(0),
	Run:   listTodos,
}

func listTodos(cmd *cobra.Command, args []string) {
	todos := todoList.List()

	sort.Slice(todos, func(i, j int) bool {
		if todos[i].Completed != todos[j].Completed {
			return !todos[i].Completed
		}
		return todos[i].Priority > todos[j].Priority
	})

	// Header
	fmt.Printf("%-4s %-10s %-10s %-14s %s\n",
		tui.RenderHeader("ID"),
		tui.RenderHeader("STATUS"),
		tui.RenderHeader("PRIORITY"),
		tui.RenderHeader("TIME"),
		tui.RenderHeader("TITLE"),
	)
	fmt.Println(tui.RenderMuted(strings.Repeat("-", 80)))

	for _, todo := range todos {
		title := todo.Title
		if len(todo.Notes) > 0 {
			title += " 󰈙"
		}

		fmt.Println(
			fmt.Sprint(todo.ID),
			tui.RenderStatus(todo.Completed),
			tui.RenderPriority(todo.Priority.String()),
			tui.RenderTimestamp(todo.TimeAgo()),
			tui.RenderTodoTitle(title, todo.Completed),
		)
	}
}
