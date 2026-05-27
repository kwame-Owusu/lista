package cmd

import "github.com/spf13/cobra"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports todos",
	Long:  "Exports todos to stdout in .md syntax",
	Args:  cobra.MinimumNArgs(0),
	Run:   exportTodos,
}

func exportTodos(cmd *cobra.Command, args []string) {}
