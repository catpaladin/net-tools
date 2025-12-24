package printers

import (
	"fmt"
	"io"
)

// PrintProcessesResults prints network connections with color formatting to the given writer
func PrintProcessesResults(w io.Writer, results []ProcessesResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No active network processes found")
		return
	}

	// Print header
	fmt.Fprintf(w, "%s%-45s %-15s %s\n", SuccessMsg("[Success] "), "Local Address", "Port", "PID/Program")
	fmt.Fprintln(w, SuccessMsg("--------------------------------------------------------------------------------"))

	// Print each connection
	for _, result := range results {
		fmt.Fprintf(w, "%s%-45s %-15s %d/%s\n",
			SuccessMsg("[Success] "), result.LocalAddr, result.LocalPort, result.PID, result.Program)
	}
}

// ProcessesResult represents a network connection (mirrors pkg/network.ProcessesResult)
type ProcessesResult struct {
	LocalAddr string
	LocalPort string
	PID       int32
	Program   string
}
