package printers

import "fmt"

// PrintNetstatResults prints network connections with color formatting
func PrintNetstatResults(results []NetstatResult) {
	if len(results) == 0 {
		fmt.Println("No network connections found")
		return
	}

	// Print header
	fmt.Printf("%s%-45s %-15s %s\n", SuccessMsg("[Success] "), "Local Address", "Port", "PID/Program")
	fmt.Println(SuccessMsg("--------------------------------------------------------------------------------"))

	// Print each connection
	for _, result := range results {
		fmt.Printf("%s%-45s %-15s %d/%s\n",
			SuccessMsg("[Success] "), result.LocalAddr, result.LocalPort, result.PID, result.Program)
	}
}

// NetstatResult represents a network connection (mirrors pkg/network.NetstatResult)
type NetstatResult struct {
	LocalAddr string
	LocalPort string
	PID       int32
	Program   string
}
