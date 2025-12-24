package main

import (
	"fmt"
	"os"

	"github.com/catpaladin/net-tools/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
