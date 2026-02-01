package tui

import (
	"fmt"
	"net"
	"strings"

	"github.com/catpaladin/net-tools/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"

	"github.com/catpaladin/net-tools/pkg/network"
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
			Foreground(styles.CurrentTheme().Primary).
			Bold(true).
			Padding(0, 1).
			Background(styles.CurrentTheme().BgSubtle).
			MarginBottom(1)

	recordValueStyle = lipgloss.NewStyle().
				Foreground(styles.CurrentTheme().FgBase).
				MarginLeft(2)

	recordEmptyStyle = lipgloss.NewStyle().
				Foreground(styles.CurrentTheme().FgMuted).
				Italic(true).
				MarginLeft(2)

	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(styles.CurrentTheme().Primary).
				Bold(true)

	tableRowStyle = lipgloss.NewStyle().
			Foreground(styles.CurrentTheme().FgBase)

	ipLabelStyle = lipgloss.NewStyle().
			Foreground(styles.CurrentTheme().Secondary).
			Bold(true).
			Width(12)

	ipValueStyle = lipgloss.NewStyle().
			Foreground(styles.CurrentTheme().Accent).
			Bold(true)

	connectionSuccessStyle = lipgloss.NewStyle().
				Foreground(styles.CurrentTheme().Success).
				Bold(true)

	connectionFailStyle = lipgloss.NewStyle().
				Foreground(styles.CurrentTheme().Error).
				Bold(true)
)

func performDNSLookup(domain string) DNSResult {
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

func parseProcessesOutput(rawOutput string) []ProcessesConnection {
	lines := strings.Split(rawOutput, "\n")
	var connections []ProcessesConnection

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Local Address") || strings.Contains(line, "---") || strings.Contains(line, "─") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 3 {
			conn := ProcessesConnection{
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

func formatProcessesResult(connections []ProcessesConnection, availableWidth int) string {
	if len(connections) == 0 {
		return recordEmptyStyle.Render("No active network processes found")
	}

	var tableContent strings.Builder

	addrHeader := "Local Address"
	portHeader := "Port"
	pidHeader := "PID/Program"

	maxAddrLen := len(addrHeader)
	maxPortLen := len(portHeader)
	maxPIDLen := len(pidHeader)

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

	addrWidth := maxAddrLen
	portWidth := maxPortLen
	pidWidth := maxPIDLen

	headerLine := fmt.Sprintf("%-*s %-*s %-*s", addrWidth, tableHeaderStyle.Render(addrHeader), portWidth, tableHeaderStyle.Render(portHeader), pidWidth, tableHeaderStyle.Render(pidHeader))
	tableContent.WriteString(headerLine)
	tableContent.WriteString("\n")

	for _, conn := range connections {
		addr := conn.LocalAddr
		port := conn.Port
		pid := conn.PIDProgram

		dataLine := fmt.Sprintf("%-*s %-*s %-*s", addrWidth, tableRowStyle.Render(addr), portWidth, tableRowStyle.Render(port), pidWidth, tableRowStyle.Render(pid))
		tableContent.WriteString(dataLine)
		tableContent.WriteString("\n")
	}

	return lipgloss.NewStyle().Padding(1, 0).Width(availableWidth).Render(tableContent.String())
}
