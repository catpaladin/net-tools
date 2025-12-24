//go:build darwin
// +build darwin

package network

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// ProcessesResult represents a network connection
type ProcessesResult struct {
	LocalAddr string
	LocalPort string
	PID       int32
	Program   string
}

// Processes retrieves TCP connections using pure Go and returns raw results
func Processes() []ProcessesResult {
	// Get all TCP connections
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil
	}

	var results []ProcessesResult

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

		results = append(results, ProcessesResult{
			LocalAddr: localAddr,
			LocalPort: localPort,
			PID:       pid,
			Program:   programName,
		})
	}

	return results
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
