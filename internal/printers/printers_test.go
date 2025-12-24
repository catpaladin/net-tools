package printers

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintDNSResult(t *testing.T) {
	tests := []struct {
		name     string
		result   DNSResult
		contains []string
	}{
		{
			name: "full results",
			result: DNSResult{
				Domain:      "example.com",
				ARecords:    []string{"1.2.3.4"},
				MXRecords:   []string{"mail.example.com 10"},
				NSRecords:   []string{"ns1.example.com"},
				CNAMERecord: "alias.example.com",
				TXTRecords:  []string{"v=spf1"},
			},
			contains: []string{"A records for example.com", "1.2.3.4", "MX records", "NS records", "CNAME record", "TXT records", "v=spf1"},
		},
		{
			name: "empty results",
			result: DNSResult{
				Domain: "example.com",
			},
			contains: []string{"No A records found", "No MX records found", "No NS records found", "No CNAME found", "No TXT records found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintDNSResult(&buf, tt.result)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
		})
	}
}

func TestPrintIPResult(t *testing.T) {
	tests := []struct {
		name     string
		ipType   string
		ip       string
		err      error
		contains []string
	}{
		{
			name:     "success",
			ipType:   "Private",
			ip:       "192.168.1.1",
			err:      nil,
			contains: []string{"[Success]", "Private IP: 192.168.1.1"},
		},
		{
			name:     "error",
			ipType:   "Public",
			ip:       "",
			err:      errors.New("network unreachable"),
			contains: []string{"[Error]", "Error getting Public IP: network unreachable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintIPResult(&buf, tt.ipType, tt.ip, tt.err)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
		})
	}
}

func TestPrintProcessesResults(t *testing.T) {
	tests := []struct {
		name     string
		results  []ProcessesResult
		contains []string
	}{
		{
			name: "with results",
			results: []ProcessesResult{
				{LocalAddr: "127.0.0.1", LocalPort: "80", PID: 1234, Program: "nginx"},
			},
			contains: []string{"Local Address", "Port", "PID/Program", "127.0.0.1", "80", "1234/nginx"},
		},
		{
			name:     "empty results",
			results:  []ProcessesResult{},
			contains: []string{"No active network processes found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintProcessesResults(&buf, tt.results)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
			}
		})
	}
}
