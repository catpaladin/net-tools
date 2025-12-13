//go:build darwin
// +build darwin

package network

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// Netstat retrieves and prints TCP connections using pure Go
func Netstat() {
	// Get all TCP connections
	connections, err := net.Connections("tcp")
	if err != nil {
		color.Red("Failed to get network connections: %v\n", err)
		return
	}

	// Print header matching Linux format
	color.Green("%-45s %-15s %s\n", "Local Address", "Port", "PID/Program")
	color.Green(strings.Repeat("-", 80))

	// Filter and display only LISTEN state connections
	for _, conn := range connections {
		if conn.Status != "LISTEN" {
			continue
		}

		// Format local address
		localAddr := formatAddress(conn.Laddr.IP)
		localPort := fmt.Sprintf("%d", conn.Laddr.Port)

		// Get process information
		pid := conn.Pid
		programName := getProgramName(pid)

		// Truncate program name if too long
		if len(programName) > 12 {
			programName = programName[:12]
		}

		// Print in format matching Linux implementation
		color.Green("%-45s %-15s %d/%s\n",
			localAddr, localPort, pid, programName)
	}
}

// formatAddress formats IP address for display
func formatAddress(ip string) string {
	if ip == "" || ip == "0.0.0.0" {
		return "*"
	}
	if ip == "::" {
		return "[::]"
	}
	// IPv6 addresses
	if strings.Contains(ip, ":") {
		return "[" + ip + "]"
	}
	return ip
}

// getProgramName retrieves the program name for a given PID
func getProgramName(pid int32) string {
	if pid == 0 {
		return ""
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return ""
	}

	name, err := proc.Name()
	if err != nil {
		return ""
	}

	return name
}
