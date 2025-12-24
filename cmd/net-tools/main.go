package main

import (
	"os"

	"github.com/catpaladin/net-tools/internal/commands"
	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "net-tools",
		Short: "A general purpose network tool",
		Long:  `A general purpose network tool`,
	}
)

func init() {
	// Add all commands to root command
	dnsCmd := commands.DNSCommand()
	rootCmd.AddCommand(dnsCmd)

	ipCmd := commands.IPCmd()
	rootCmd.AddCommand(ipCmd)

	portCmd := commands.PortCmd()
	rootCmd.AddCommand(portCmd)

	processesCmd := commands.ProcessesCmd()
	rootCmd.AddCommand(processesCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
