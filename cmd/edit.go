package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "edit the text of a todo",
	Long:  "edit the text of a todo, given the correct ID",
	Args:  cobra.MinimumNArgs(2),
	Run:   editTodo,
}

func editTodo(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error occurred converting id to int: %s\n", err)
	}
	editText := strings.Join(args[1:], " ")
	err = todoList.Edit(id, editText)
	if err != nil {
		fmt.Printf("Error editing todo with id: %d, and string: %s, %s\n", id, editText, err)
	}
	saveTodos()
}
