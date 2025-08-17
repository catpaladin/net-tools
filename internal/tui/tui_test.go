package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if len(m.tabs) != 4 {
		t.Errorf("Expected 4 tabs, got %d", len(m.tabs))
	}

	if m.activeTab != 0 {
		t.Errorf("Expected active tab to be 0, got %d", m.activeTab)
	}

	if !m.tabs[0].Active {
		t.Error("Expected first tab to be active")
	}
}

func TestTUIRendering(t *testing.T) {
	m := NewModel()

	m.Init()

	model, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = model.(MainModel)

	view := m.View()

	if view == "" {
		t.Error("TUI view is empty")
	}

	if !strings.Contains(view, "NET-TOOLS") {
		t.Error("TUI view does not contain expected title")
	}

	if !strings.Contains(view, "DNS Lookup") {
		t.Error("TUI view does not contain DNS Lookup tab")
	}
}

func TestTabSwitching(t *testing.T) {
	m := NewModel()
	model, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = model.(MainModel)

	keyMsg := tea.KeyMsg{Type: tea.KeyRight}
	model, _ = m.Update(keyMsg)
	m = model.(MainModel)

	if m.activeTab != 1 {
		t.Errorf("Expected active tab to be 1 after right key, got %d", m.activeTab)
	}

	if !m.tabs[1].Active {
		t.Error("Expected second tab to be active")
	}

	if m.tabs[0].Active {
		t.Error("Expected first tab to be inactive")
	}
}

func TestDNSResultFormatting(t *testing.T) {
	result := DNSResult{
		ARecords:   []string{"192.168.1.1", "10.0.0.1"},
		TXTRecords: []string{"v=spf1 -all"},
	}

	formatted := formatDNSResult("example.com", result)

	if !strings.Contains(formatted, "A Records") {
		t.Error("Formatted result should contain A Records section")
	}

	if !strings.Contains(formatted, "192.168.1.1") {
		t.Error("Formatted result should contain A record value")
	}

	if !strings.Contains(formatted, "TXT Records") {
		t.Error("Formatted result should contain TXT Records section")
	}
}

func TestConnectionResultFormatting(t *testing.T) {
	successResult := formatConnectionResult("google.com", "80", true, nil)

	if !strings.Contains(successResult, "Connection successful") {
		t.Error("Success result should indicate successful connection")
	}

	if !strings.Contains(successResult, "google.com:80") {
		t.Error("Success result should contain target address")
	}

	failResult := formatConnectionResult("localhost", "9999", false, fmt.Errorf("connection refused"))

	if !strings.Contains(failResult, "Connection failed") {
		t.Error("Fail result should indicate failed connection")
	}
}

func TestNetstatTableFormatting(t *testing.T) {
	connections := []NetstatConnection{
		{LocalAddr: "127.0.0.1", Port: "8080", PIDProgram: "chrome 1234"},
		{LocalAddr: "0.0.0.0", Port: "22", PIDProgram: "sshd 5678"},
		{LocalAddr: "192.168.1.100", Port: "443", PIDProgram: "nginx 9012"},
	}

	narrowResult := formatNetstatResult(connections, 50)
	wideResult := formatNetstatResult(connections, 120)

	if !strings.Contains(narrowResult, "Local Address") {
		t.Error("Narrow table should contain header")
	}

	if !strings.Contains(wideResult, "127.0.0.1") {
		t.Error("Wide table should contain connection data")
	}

	narrowLines := strings.Split(narrowResult, "\n")
	wideLines := strings.Split(wideResult, "\n")

	if len(narrowLines) < 3 {
		t.Error("Table should have header, separator, and data rows")
	}

	if len(wideLines) < 3 {
		t.Error("Table should have header, separator, and data rows")
	}

	if len(narrowResult) == 0 {
		t.Error("Narrow table should not be empty")
	}

	if len(wideResult) == 0 {
		t.Error("Wide table should not be empty")
	}
}

func TestNewlineFormatting(t *testing.T) {
	result := DNSResult{
		ARecords:   []string{"192.168.1.1", "10.0.0.1"},
		TXTRecords: []string{"v=spf1 -all"},
	}

	formatted := formatDNSResult("example.com", result)

	sections := strings.Split(formatted, "\n\n")
	if len(sections) < 2 {
		t.Error("DNS result should have properly separated sections")
	}

	for _, section := range sections {
		if strings.HasPrefix(section, "\n") || strings.HasSuffix(section, "\n") {
			t.Error("Sections should not have leading or trailing newlines")
		}
	}

	connectionResult := formatConnectionResult("google.com", "80", true, nil)
	connectionLines := strings.Split(connectionResult, "\n")

	if len(connectionLines) != 3 {
		t.Errorf("Connection result should have exactly 3 lines, got %d", len(connectionLines))
	}

	for i, line := range connectionLines {
		if strings.TrimSpace(line) == "" {
			t.Errorf("Line %d should not be empty", i)
		}
	}

	ipResult := formatIPResult("both", "192.168.1.1", "203.0.113.1", nil, nil)
	ipLines := strings.Split(ipResult, "\n")

	if len(ipLines) != 2 {
		t.Errorf("IP result for 'both' should have exactly 2 lines, got %d", len(ipLines))
	}
}

func TestPlainFormatters(t *testing.T) {
	result := DNSResult{
		ARecords:   []string{"192.168.1.1", "10.0.0.1"},
		TXTRecords: []string{"v=spf1 -all"},
	}

	plainFormatted := formatDNSResultPlain("example.com", result)

	if !strings.Contains(plainFormatted, "A Records:") {
		t.Error("Plain DNS result should contain A Records section")
	}

	if !strings.Contains(plainFormatted, "192.168.1.1") {
		t.Error("Plain DNS result should contain A record value")
	}

	if strings.Contains(plainFormatted, "\x1b[") {
		t.Error("Plain DNS result should not contain ANSI escape codes")
	}

	plainConnection := formatConnectionResultPlain("google.com", "80", true, nil)

	if !strings.Contains(plainConnection, "Connection successful") {
		t.Error("Plain connection result should indicate success")
	}

	if strings.Contains(plainConnection, "\x1b[") {
		t.Error("Plain connection result should not contain ANSI escape codes")
	}

	plainIP := formatIPResultPlain("both", "192.168.1.1", "203.0.113.1", nil, nil)

	if !strings.Contains(plainIP, "Private IP: 192.168.1.1") {
		t.Error("Plain IP result should contain private IP")
	}

	if !strings.Contains(plainIP, "Public IP: 203.0.113.1") {
		t.Error("Plain IP result should contain public IP")
	}

	if strings.Contains(plainIP, "\x1b[") {
		t.Error("Plain IP result should not contain ANSI escape codes")
	}

	connections := []NetstatConnection{
		{LocalAddr: "127.0.0.1", Port: "8080", PIDProgram: "chrome 1234"},
		{LocalAddr: "0.0.0.0", Port: "22", PIDProgram: "sshd 5678"},
	}

	plainNetstat := formatNetstatResultPlain(connections, 80)

	if !strings.Contains(plainNetstat, "Local Address") {
		t.Error("Plain netstat result should contain header")
	}

	if !strings.Contains(plainNetstat, "127.0.0.1") {
		t.Error("Plain netstat result should contain connection data")
	}

	if strings.Contains(plainNetstat, "\x1b[") {
		t.Error("Plain netstat result should not contain ANSI escape codes")
	}
}

func TestTUIResponsivenessSizing(t *testing.T) {
	m := NewModel()

	// Test small terminal size
	smallModel, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 15})
	m = smallModel.(MainModel)

	if m.width != 50 || m.height != 15 {
		t.Error("Model should update to match terminal size")
	}

	if m.digModel.viewport.Width <= 0 || m.digModel.viewport.Height <= 0 {
		t.Error("Viewports should have positive dimensions even in small terminals")
	}

	// Test large terminal size
	largeModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = largeModel.(MainModel)

	if m.width != 120 || m.height != 40 {
		t.Error("Model should update to match larger terminal size")
	}

	if m.digModel.viewport.Width <= 0 || m.digModel.viewport.Height <= 0 {
		t.Error("Viewports should have positive dimensions in large terminals")
	}

	// Verify the view renders without panic
	view := m.View()
	if view == "" {
		t.Error("View should render content at any reasonable terminal size")
	}

	// Test minimum size constraints
	tinyModel, _ := m.Update(tea.WindowSizeMsg{Width: 10, Height: 5})
	m = tinyModel.(MainModel)

	view = m.View()
	if view == "" {
		t.Error("View should still render at minimum size")
	}
}
