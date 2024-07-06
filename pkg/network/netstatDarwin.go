//go:build darwin
// +build darwin

package network

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// Netstat retrieves and prints TCP connections
func Netstat() {
	cmd := exec.Command("netstat", "-an", "|", "grep", "LISTEN")
	stdout, err := cmd.Output()
	if err != nil {
		color.Red("Failed to run netstat command: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tcp") || strings.HasPrefix(line, "tcp4") || strings.HasPrefix(line, "tcp6") {
			color.Green(line)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error reading netstat output: %v\n", err)
	}
}
