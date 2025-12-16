package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kwame-Owusu/lista/internal/tui"
	"github.com/spf13/cobra"
	"os"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Run lista in interactive TUI mode",
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.NewModel(todoList, dataFile)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}
