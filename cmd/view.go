package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/tui"

	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [id]",
	Short: "View a todo",
	Long:  "View a todo in its entirety, with notes if present",
	Args:  cobra.MinimumNArgs(1),
	Run:   viewTodo,
}

func viewTodo(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(tui.RenderError("Invalid todo ID"))
		return
	}

	todo, err := todoList.GetByID(id)
	if err != nil {
		fmt.Println(
			tui.RenderError(fmt.Sprintf("No todo with ID %d exists", id)),
		)
		return
	}

	renderViewText(todo)
}

func renderViewText(todo *models.Todo) {
	fmt.Println(tui.RenderSectionTitle("Todo Details"))

	fmt.Println(
		tui.RenderLabel("ID:"),
		tui.RenderValue(fmt.Sprint(todo.ID)),
	)

	fmt.Println(
		tui.RenderLabel("Title:"),
		tui.RenderTodoTitle(todo.Title, todo.Completed),
	)

	fmt.Println(
		tui.RenderLabel("Priority:"),
		tui.RenderPriority(todo.Priority.String()),
	)

	fmt.Println(
		tui.RenderLabel("Created:"),
		tui.RenderTimestamp(todo.TimeAgo()),
	)

	fmt.Println(
		tui.RenderLabel("Status:"),
		tui.RenderStatus(todo.Completed),
	)

	fmt.Println(
		tui.RenderLabel("Notes:"),
	)

	if strings.TrimSpace(todo.Notes) == "" {
		fmt.Println(tui.RenderMuted("  (none)"))
	} else {
		// Indent notes block for readability
		for _, line := range strings.Split(todo.Notes, "\n") {
			fmt.Println("  " + tui.RenderValue(line))
		}
	}
}
