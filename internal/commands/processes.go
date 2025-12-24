package commands

import (
	"os"

	"github.com/catpaladin/net-tools/internal/printers"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/spf13/cobra"
)

// ProcessesCmd creates and returns the processes cobra command
func ProcessesCmd() *cobra.Command {
	processesCmd := &cobra.Command{
		Use:   "processes",
		Short: "Lists active network processes and connections",
		Long:  "Lists all active network connections along with their local addresses, ports, and associated PIDs/programs.",
		Run: func(cmd *cobra.Command, args []string) {
			// Get raw results from pkg/network
			results := network.Processes()
			// Format and print using internal printers
			printers.PrintProcessesResults(os.Stdout, convertProcessesResults(results))
		},
	}

	return processesCmd
}

// convertProcessesResults converts pkg/network.ProcessesResult to printers.ProcessesResult
func convertProcessesResults(results []network.ProcessesResult) []printers.ProcessesResult {
	var converted []printers.ProcessesResult
	for _, result := range results {
		converted = append(converted, printers.ProcessesResult{
			LocalAddr: result.LocalAddr,
			LocalPort: result.LocalPort,
			PID:       result.PID,
			Program:   result.Program,
		})
	}
	return converted
}
