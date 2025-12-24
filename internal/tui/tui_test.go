package tui

import (
	"fmt"
	"net"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
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

	plainFail := formatConnectionResultPlain("localhost", "9999", false, fmt.Errorf("connection refused"))
	if !strings.Contains(plainFail, "Connection failed") {
		t.Error("Plain fail result should indicate failed connection")
	}

	plainSuccess := formatConnectionResultPlain("google.com", "80", true, nil)
	if !strings.Contains(plainSuccess, "Connection successful") {
		t.Error("Plain success result should indicate success")
	}
}

func TestIPResultFormattingExhaustive(t *testing.T) {
	// Already have TestIPResultFormatting, adding more cases here if needed
	// or I can just update the existing one.
}

func TestProcessesTableFormattingExhaustive(t *testing.T) {
	connections := []ProcessesConnection{
		{LocalAddr: "127.0.0.1", Port: "80", PIDProgram: "nginx"},
	}

	// Test with data
	res := formatProcessesResult(connections, 100)
	assert.Contains(t, res, "Local Address")

	// Test empty
	emptyRes := formatProcessesResult([]ProcessesConnection{}, 100)
	assert.Contains(t, emptyRes, "No active network processes found")

	// Test narrow
	narrow := formatProcessesResult(connections, 40)
	assert.NotEmpty(t, narrow)
}

func TestDNSResultFormattingPlain(t *testing.T) {
	result := DNSResult{
		ARecords:    []string{"1.2.3.4"},
		MXRecords:   []*net.MX{{Host: "mail", Pref: 10}},
		NSRecords:   []*net.NS{{Host: "ns"}},
		CNAMERecord: "alias",
		TXTRecords:  []string{"txt"},
	}

	formatted := formatDNSResultPlain("example.com", result)
	assert.Contains(t, formatted, "A Records:")
	assert.Contains(t, formatted, "MX Records:")

	empty := formatDNSResultPlain("example.com", DNSResult{})
	assert.Contains(t, empty, "No A records found")
}

func TestProcessesTableFormatting(t *testing.T) {
	connections := []ProcessesConnection{
		{LocalAddr: "127.0.0.1", Port: "8080", PIDProgram: "chrome 1234"},
		{LocalAddr: "0.0.0.0", Port: "22", PIDProgram: "sshd 5678"},
		{LocalAddr: "192.168.1.100", Port: "443", PIDProgram: "nginx 9012"},
	}

	narrowResult := formatProcessesResult(connections, 50)
	wideResult := formatProcessesResult(connections, 120)

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

	connections := []ProcessesConnection{
		{LocalAddr: "127.0.0.1", Port: "8080", PIDProgram: "chrome 1234"},
		{LocalAddr: "0.0.0.0", Port: "22", PIDProgram: "sshd 5678"},
	}

	plainProcesses := formatProcessesResultPlain(connections, 80)

	if !strings.Contains(plainProcesses, "Local Address") {
		t.Error("Plain processes result should contain header")
	}

	if !strings.Contains(plainProcesses, "127.0.0.1") {
		t.Error("Plain processes result should contain connection data")
	}

	if strings.Contains(plainProcesses, "\x1b[") {
		t.Error("Plain processes result should not contain ANSI escape codes")
	}
}

func TestIPResultFormatting(t *testing.T) {
	tests := []struct {
		name       string
		ipType     string
		privateIP  string
		publicIP   string
		privateErr error
		publicErr  error
		contains   []string
	}{
		{
			name:      "private success",
			ipType:    "private",
			privateIP: "192.168.1.1",
			contains:  []string{"Private IP:", "192.168.1.1"},
		},
		{
			name:     "public success",
			ipType:   "public",
			publicIP: "1.2.3.4",
			contains: []string{"Public IP:", "1.2.3.4"},
		},
		{
			name:      "both success",
			ipType:    "both",
			privateIP: "192.168.1.1",
			publicIP:  "1.2.3.4",
			contains:  []string{"Private IP:", "192.168.1.1", "Public IP:", "1.2.3.4"},
		},
		{
			name:       "private error",
			ipType:     "private",
			privateErr: fmt.Errorf("failed"),
			contains:   []string{"Private IP:", "Error"},
		},
		{
			name:      "public error",
			ipType:    "public",
			publicErr: fmt.Errorf("failed"),
			contains:  []string{"Public IP:", "Error"},
		},
		{
			name:       "both error",
			ipType:     "both",
			privateErr: fmt.Errorf("p_fail"),
			publicErr:  fmt.Errorf("pub_fail"),
			contains:   []string{"Private IP:", "p_fail", "Public IP:", "pub_fail"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := formatIPResult(tt.ipType, tt.privateIP, tt.publicIP, tt.privateErr, tt.publicErr)
			for _, s := range tt.contains {
				if !strings.Contains(formatted, s) {
					t.Errorf("Expected formatted result to contain %q", s)
				}
			}

			plain := formatIPResultPlain(tt.ipType, tt.privateIP, tt.publicIP, tt.privateErr, tt.publicErr)
			for _, s := range tt.contains {
				if !strings.Contains(plain, s) {
					t.Errorf("Expected plain result to contain %q", s)
				}
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port    string
		wantErr bool
	}{
		{"80", false},
		{"1", false},
		{"65535", false},
		{"0", true},
		{"65536", true},
		{"abc", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.port, func(t *testing.T) {
			err := validatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePort(%q) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestDNSResultFormattingExhaustive(t *testing.T) {
	result := DNSResult{
		ARecords:    []string{"1.2.3.4"},
		MXRecords:   []*net.MX{{Host: "mail", Pref: 10}},
		NSRecords:   []*net.NS{{Host: "ns"}},
		CNAMERecord: "alias",
		TXTRecords:  []string{"txt"},
	}

	formatted := formatDNSResult("example.com", result)
	assert.Contains(t, formatted, "A Records")
	assert.Contains(t, formatted, "MX Records")
	assert.Contains(t, formatted, "NS Records")
	assert.Contains(t, formatted, "CNAME Record")
	assert.Contains(t, formatted, "TXT Records")

	empty := formatDNSResult("example.com", DNSResult{})
	assert.Contains(t, empty, "No A records found")
	assert.Contains(t, empty, "No MX records found")
}

func TestParseProcessesOutput(t *testing.T) {
	raw := "127.0.0.1 80 nginx 1234\n::1 443 chrome 5678\n"
	conns := parseProcessesOutput(raw)
	assert.Equal(t, 2, len(conns))
	assert.Equal(t, "127.0.0.1", conns[0].LocalAddr)
	assert.Equal(t, "80", conns[0].Port)
	assert.Equal(t, "nginx 1234", conns[0].PIDProgram)
}

func TestTUIResponsivenessSizing(t *testing.T) {
	m := NewModel()

	// Test small terminal size
	smallModel, _ := m.Update(tea.WindowSizeMsg{Width: 50, Height: 15})
	m = smallModel.(MainModel)

	if m.width != 50 || m.height != 15 {
		t.Error("Model should update to match terminal size")
	}

	if m.dnsModel.viewport.Width <= 0 || m.dnsModel.viewport.Height <= 0 {
		t.Error("Viewports should have positive dimensions even in small terminals")
	}

	// Test large terminal size
	largeModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = largeModel.(MainModel)

	if m.width != 120 || m.height != 40 {
		t.Error("Model should update to match larger terminal size")
	}

	if m.dnsModel.viewport.Width <= 0 || m.dnsModel.viewport.Height <= 0 {
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
