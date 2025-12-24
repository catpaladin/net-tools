package commands

import (
	"fmt"
	"os"

	"github.com/catpaladin/net-tools/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// TUICmd creates and returns the tui cobra command
func TUICmd() *cobra.Command {
	tuiCmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive TUI interface",
		Long:  `Launch the interactive Terminal User Interface for net-tools with tabbed navigation`,
		Run: func(cmd *cobra.Command, args []string) {
			m := tui.NewModel()
			p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

			if _, err := p.Run(); err != nil {
				fmt.Printf("Error running TUI: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return tuiCmd
}
