package tui

import (
	"fmt"
	"strings"
)

func formatDNSResultPlain(domain string, result DNSResult) string {
	var sections []string

	if len(result.ARecords) > 0 {
		section := "A Records:\n"
		for _, record := range result.ARecords {
			section += fmt.Sprintf("  %s\n", record)
		}
		sections = append(sections, strings.TrimSpace(section))
	} else {
		sections = append(sections, "A Records:\n  No A records found")
	}

	if len(result.MXRecords) > 0 {
		section := "MX Records:\n"
		for _, mx := range result.MXRecords {
			section += fmt.Sprintf("  %s (priority: %d)\n", mx.Host, mx.Pref)
		}
		sections = append(sections, strings.TrimSpace(section))
	} else {
		sections = append(sections, "MX Records:\n  No MX records found")
	}

	if len(result.NSRecords) > 0 {
		section := "NS Records:\n"
		for _, ns := range result.NSRecords {
			section += fmt.Sprintf("  %s\n", ns.Host)
		}
		sections = append(sections, strings.TrimSpace(section))
	} else {
		sections = append(sections, "NS Records:\n  No NS records found")
	}

	if result.CNAMERecord != "" {
		sections = append(sections, fmt.Sprintf("CNAME Record:\n  %s", result.CNAMERecord))
	} else {
		sections = append(sections, "CNAME Record:\n  No CNAME record found")
	}

	if len(result.TXTRecords) > 0 {
		section := "TXT Records:\n"
		for _, txt := range result.TXTRecords {
			section += fmt.Sprintf("  %s\n", txt)
		}
		sections = append(sections, strings.TrimSpace(section))
	} else {
		sections = append(sections, "TXT Records:\n  No TXT records found")
	}

	return strings.Join(sections, "\n\n")
}

func formatConnectionResultPlain(host, port string, success bool, err error) string {
	var lines []string

	target := fmt.Sprintf("%s:%s", host, port)

	if success {
		lines = append(lines, "✓ Connection successful")
		lines = append(lines, fmt.Sprintf("Target: %s", target))
		lines = append(lines, "Status: Port is open and accepting connections")
	} else {
		lines = append(lines, "✗ Connection failed")
		lines = append(lines, fmt.Sprintf("Target: %s", target))
		if err != nil {
			lines = append(lines, fmt.Sprintf("Error: %s", err.Error()))
		} else {
			lines = append(lines, "Status: Port is closed or not reachable")
		}
	}

	return strings.Join(lines, "\n")
}

func formatIPResultPlain(ipType string, privateIP, publicIP string, privateErr, publicErr error) string {
	var lines []string

	switch ipType {
	case "private":
		if privateErr != nil {
			lines = append(lines, fmt.Sprintf("Private IP: Error - %v", privateErr))
		} else {
			lines = append(lines, fmt.Sprintf("Private IP: %s", privateIP))
		}

	case "public":
		if publicErr != nil {
			lines = append(lines, fmt.Sprintf("Public IP: Error - %v", publicErr))
		} else {
			lines = append(lines, fmt.Sprintf("Public IP: %s", publicIP))
		}

	case "both":
		if privateErr != nil {
			lines = append(lines, fmt.Sprintf("Private IP: Error - %v", privateErr))
		} else {
			lines = append(lines, fmt.Sprintf("Private IP: %s", privateIP))
		}

		if publicErr != nil {
			lines = append(lines, fmt.Sprintf("Public IP: Error - %v", publicErr))
		} else {
			lines = append(lines, fmt.Sprintf("Public IP: %s", publicIP))
		}
	}

	return strings.Join(lines, "\n")
}

func formatNetstatResultPlain(connections []NetstatConnection, availableWidth int) string {
	if len(connections) == 0 {
		return "No active network connections found"
	}

	if availableWidth < 40 {
		availableWidth = 40
	}

	tableWidth := availableWidth - 8

	addrHeader := "Local Address"
	portHeader := "Port"
	pidHeader := "PID/Program"

	minAddrLen := len(addrHeader)
	minPortLen := len(portHeader)
	minPIDLen := len(pidHeader)

	maxAddrLen := minAddrLen
	maxPortLen := minPortLen
	maxPIDLen := minPIDLen

	for _, conn := range connections {
		if len(conn.LocalAddr) > maxAddrLen {
			maxAddrLen = len(conn.LocalAddr)
		}
		if len(conn.Port) > maxPortLen {
			maxPortLen = len(conn.Port)
		}
		if len(conn.PIDProgram) > maxPIDLen {
			maxPIDLen = len(conn.PIDProgram)
		}
	}

	totalMinWidth := minAddrLen + minPortLen + minPIDLen + 6
	totalMaxWidth := maxAddrLen + maxPortLen + maxPIDLen + 6

	var addrWidth, portWidth, pidWidth int

	if totalMaxWidth <= tableWidth {
		addrWidth = maxAddrLen
		portWidth = maxPortLen
		pidWidth = maxPIDLen
	} else if totalMinWidth <= tableWidth {
		remainingWidth := tableWidth - 6

		addrWidth = minAddrLen + (remainingWidth-minAddrLen-minPortLen-minPIDLen)*2/5
		portWidth = minPortLen + (remainingWidth-minAddrLen-minPortLen-minPIDLen)*1/5
		pidWidth = remainingWidth - addrWidth - portWidth

		if addrWidth < minAddrLen {
			addrWidth = minAddrLen
		}
		if portWidth < minPortLen {
			portWidth = minPortLen
		}
		if pidWidth < minPIDLen {
			pidWidth = minPIDLen
		}
	} else {
		addrWidth = minAddrLen
		portWidth = minPortLen
		pidWidth = minPIDLen
	}

	var lines []string

	headerLine := fmt.Sprintf("%-*s %-*s %-*s", addrWidth, addrHeader, portWidth, portHeader, pidWidth, pidHeader)
	lines = append(lines, headerLine)

	separatorWidth := addrWidth + portWidth + pidWidth + 2
	if separatorWidth > tableWidth {
		separatorWidth = tableWidth
	}
	separator := strings.Repeat("-", separatorWidth)
	lines = append(lines, separator)

	for _, conn := range connections {
		addr := conn.LocalAddr
		port := conn.Port
		pid := conn.PIDProgram

		if len(addr) > addrWidth {
			addr = addr[:addrWidth-3] + "..."
		}
		if len(port) > portWidth {
			port = port[:portWidth-3] + "..."
		}
		if len(pid) > pidWidth {
			pid = pid[:pidWidth-3] + "..."
		}

		dataLine := fmt.Sprintf("%-*s %-*s %-*s", addrWidth, addr, portWidth, port, pidWidth, pid)
		lines = append(lines, dataLine)
	}

	return strings.Join(lines, "\n")
}
