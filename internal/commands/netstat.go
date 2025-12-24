package commands

import (
	"github.com/catpaladin/net-tools/internal/printers"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/spf13/cobra"
)

// NetstatCmd creates and returns the netstat cobra command
func NetstatCmd() *cobra.Command {
	netstatCmd := &cobra.Command{
		Use:   "netstat",
		Short: "Performs a netstat on the host",
		Long:  "Performs a netstat on the host",
		Run: func(cmd *cobra.Command, args []string) {
			// Get raw results from pkg/network
			results := network.Netstat()
			// Format and print using internal printers
			printers.PrintNetstatResults(convertNetstatResults(results))
		},
	}

	return netstatCmd
}

// convertNetstatResults converts pkg/network.NetstatResult to printers.NetstatResult
func convertNetstatResults(results []network.NetstatResult) []printers.NetstatResult {
	var converted []printers.NetstatResult
	for _, result := range results {
		converted = append(converted, printers.NetstatResult{
			LocalAddr: result.LocalAddr,
			LocalPort: result.LocalPort,
			PID:       result.PID,
			Program:   result.Program,
		})
	}
	return converted
}
