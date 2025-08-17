package tui

import (
	"fmt"
	"net"
	"strings"

	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/lipgloss"
)

type DNSResult struct {
	ARecords    []string
	MXRecords   []*net.MX
	NSRecords   []*net.NS
	CNAMERecord string
	TXTRecords  []string
}

var (
	recordTypeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	recordValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				MarginLeft(2)

	recordEmptyStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				MarginLeft(2)

	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Underline(true)

	tableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	ipLabelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Width(12)

	ipValueStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	connectionSuccessStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)

	connectionFailStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)
)

func performDigLookup(domain string) DNSResult {
	nh := network.NetHostLookup{}

	var result DNSResult

	if addrs, err := nh.LookupHost(domain); err == nil {
		result.ARecords = addrs
	}

	if mxRecords, err := nh.LookupMX(domain); err == nil {
		result.MXRecords = mxRecords
	}

	if nsRecords, err := nh.LookupNS(domain); err == nil {
		result.NSRecords = nsRecords
	}

	if cname, err := nh.LookupCNAME(domain); err == nil {
		result.CNAMERecord = cname
	}

	if txtRecords, err := nh.LookupTXT(domain); err == nil {
		result.TXTRecords = txtRecords
	}

	return result
}

func formatDNSResult(domain string, result DNSResult) string {
	var sections []string

	if len(result.ARecords) > 0 {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("A Records"))
		section.WriteString("\n")
		for _, record := range result.ARecords {
			section.WriteString(recordValueStyle.Render("• " + record))
			section.WriteString("\n")
		}
		sections = append(sections, strings.TrimSpace(section.String()))
	} else {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("A Records"))
		section.WriteString("\n")
		section.WriteString(recordEmptyStyle.Render("No A records found"))
		sections = append(sections, strings.TrimSpace(section.String()))
	}

	if len(result.MXRecords) > 0 {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("MX Records"))
		section.WriteString("\n")
		for _, mx := range result.MXRecords {
			section.WriteString(recordValueStyle.Render(fmt.Sprintf("• %s (priority: %d)", mx.Host, mx.Pref)))
			section.WriteString("\n")
		}
		sections = append(sections, strings.TrimSpace(section.String()))
	} else {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("MX Records"))
		section.WriteString("\n")
		section.WriteString(recordEmptyStyle.Render("No MX records found"))
		sections = append(sections, strings.TrimSpace(section.String()))
	}

	if len(result.NSRecords) > 0 {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("NS Records"))
		section.WriteString("\n")
		for _, ns := range result.NSRecords {
			section.WriteString(recordValueStyle.Render("• " + ns.Host))
			section.WriteString("\n")
		}
		sections = append(sections, strings.TrimSpace(section.String()))
	} else {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("NS Records"))
		section.WriteString("\n")
		section.WriteString(recordEmptyStyle.Render("No NS records found"))
		sections = append(sections, strings.TrimSpace(section.String()))
	}

	if result.CNAMERecord != "" {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("CNAME Record"))
		section.WriteString("\n")
		section.WriteString(recordValueStyle.Render("• " + result.CNAMERecord))
		sections = append(sections, strings.TrimSpace(section.String()))
	} else {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("CNAME Record"))
		section.WriteString("\n")
		section.WriteString(recordEmptyStyle.Render("No CNAME record found"))
		sections = append(sections, strings.TrimSpace(section.String()))
	}

	if len(result.TXTRecords) > 0 {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("TXT Records"))
		section.WriteString("\n")
		for _, txt := range result.TXTRecords {
			section.WriteString(recordValueStyle.Render("• " + txt))
			section.WriteString("\n")
		}
		sections = append(sections, strings.TrimSpace(section.String()))
	} else {
		var section strings.Builder
		section.WriteString(recordTypeStyle.Render("TXT Records"))
		section.WriteString("\n")
		section.WriteString(recordEmptyStyle.Render("No TXT records found"))
		sections = append(sections, strings.TrimSpace(section.String()))
	}

	return strings.Join(sections, "\n\n")
}

func formatConnectionResult(host, port string, success bool, err error) string {
	var lines []string

	target := fmt.Sprintf("%s:%s", host, port)

	if success {
		lines = append(lines, connectionSuccessStyle.Render("✓ Connection successful"))
		lines = append(lines, recordValueStyle.Render(fmt.Sprintf("Target: %s", target)))
		lines = append(lines, recordValueStyle.Render("Status: Port is open and accepting connections"))
	} else {
		lines = append(lines, connectionFailStyle.Render("✗ Connection failed"))
		lines = append(lines, recordValueStyle.Render(fmt.Sprintf("Target: %s", target)))
		if err != nil {
			lines = append(lines, recordValueStyle.Render(fmt.Sprintf("Error: %s", err.Error())))
		} else {
			lines = append(lines, recordValueStyle.Render("Status: Port is closed or not reachable"))
		}
	}

	return strings.Join(lines, "\n")
}

func formatIPResult(ipType string, privateIP, publicIP string, privateErr, publicErr error) string {
	var lines []string

	switch ipType {
	case "private":
		line := ipLabelStyle.Render("Private IP:") + " "
		if privateErr != nil {
			line += recordEmptyStyle.Render(fmt.Sprintf("Error - %v", privateErr))
		} else {
			line += ipValueStyle.Render(privateIP)
		}
		lines = append(lines, line)

	case "public":
		line := ipLabelStyle.Render("Public IP:") + " "
		if publicErr != nil {
			line += recordEmptyStyle.Render(fmt.Sprintf("Error - %v", publicErr))
		} else {
			line += ipValueStyle.Render(publicIP)
		}
		lines = append(lines, line)

	case "both":
		privateLine := ipLabelStyle.Render("Private IP:") + " "
		if privateErr != nil {
			privateLine += recordEmptyStyle.Render(fmt.Sprintf("Error - %v", privateErr))
		} else {
			privateLine += ipValueStyle.Render(privateIP)
		}
		lines = append(lines, privateLine)

		publicLine := ipLabelStyle.Render("Public IP:") + " "
		if publicErr != nil {
			publicLine += recordEmptyStyle.Render(fmt.Sprintf("Error - %v", publicErr))
		} else {
			publicLine += ipValueStyle.Render(publicIP)
		}
		lines = append(lines, publicLine)
	}

	return strings.Join(lines, "\n")
}

func parseNetstatOutput(rawOutput string) []NetstatConnection {
	lines := strings.Split(rawOutput, "\n")
	var connections []NetstatConnection

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Local Address") || strings.Contains(line, "---") || strings.Contains(line, "─") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 3 {
			conn := NetstatConnection{
				LocalAddr:  parts[0],
				Port:       parts[1],
				PIDProgram: strings.Join(parts[2:], " "),
			}

			if len(conn.Port) > 6 {
				conn.Port = conn.Port[:6]
			}

			connections = append(connections, conn)
		}
	}

	return connections
}

func formatNetstatResult(connections []NetstatConnection, availableWidth int) string {
	if len(connections) == 0 {
		return recordEmptyStyle.Render("No active network connections found")
	}

	if availableWidth < 40 {
		availableWidth = 40
	}

	tableWidth := availableWidth - 8

	var content strings.Builder

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

	var tableContent strings.Builder

	headerLine := fmt.Sprintf("%-*s %-*s %-*s", addrWidth, addrHeader, portWidth, portHeader, pidWidth, pidHeader)
	tableContent.WriteString(headerLine)
	tableContent.WriteString("\n")

	separatorWidth := addrWidth + portWidth + pidWidth + 2
	if separatorWidth > tableWidth {
		separatorWidth = tableWidth
	}
	separator := strings.Repeat("─", separatorWidth)
	tableContent.WriteString(separator)
	tableContent.WriteString("\n")

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
		tableContent.WriteString(dataLine)
		tableContent.WriteString("\n")
	}

	styledTable := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(1).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(strings.TrimSpace(tableContent.String()))

	content.WriteString(styledTable)

	return content.String()
}
