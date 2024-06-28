package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Define color styles
	successMsg = color.New(color.FgGreen).SprintFunc()
	errorMsg   = color.New(color.FgRed).SprintFunc()
	dataMsg    = color.New(color.FgCyan).SprintFunc()

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "nt",
		Short: "A general purpose network tool",
		Long:  `A general purpose network tool`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
