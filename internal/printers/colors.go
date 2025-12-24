package printers

import "github.com/charmbracelet/lipgloss"

// Color functions for consistent output formatting
var (
	SuccessMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render  // Bright green
	ErrorMsg   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render // Bright red
	DataMsg    = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render  // Bright blue
	WarningMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render // Orange
)
