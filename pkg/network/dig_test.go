package network

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockHostLookup is a mock implementation of the HostLookup interface.
type MockHostLookup struct {
	LookupHostFunc  func(domain string) ([]string, error)
	LookupMXFunc    func(domain string) ([]*net.MX, error)
	LookupNSFunc    func(domain string) ([]*net.NS, error)
	LookupCNAMEFunc func(domain string) (string, error)
	LookupTXTFunc   func(domain string) ([]string, error)
}

func (m MockHostLookup) LookupHost(domain string) ([]string, error) {
	return m.LookupHostFunc(domain)
}

func (m MockHostLookup) LookupMX(domain string) ([]*net.MX, error) {
	return m.LookupMXFunc(domain)
}

func (m MockHostLookup) LookupNS(domain string) ([]*net.NS, error) {
	return m.LookupNSFunc(domain)
}

func (m MockHostLookup) LookupCNAME(domain string) (string, error) {
	return m.LookupCNAMEFunc(domain)
}

func (m MockHostLookup) LookupTXT(domain string) ([]string, error) {
	return m.LookupTXTFunc(domain)
}

func TestLookupARecords(t *testing.T) {
	tests := []struct {
		domain   string
		addrs    []string
		err      error
		expected []string
	}{
		{"example.com", []string{"93.184.216.34"}, nil, []string{"93.184.216.34"}},
		{"invalid-domain", nil, fmt.Errorf("no such host"), []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			mockLookup := MockHostLookup{
				LookupHostFunc: func(domain string) ([]string, error) {
					assert.Equal(t, tt.domain, domain)
					return tt.addrs, tt.err
				},
			}

			result := lookupARecords(mockLookup, tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupMXRecords(t *testing.T) {
	tests := []struct {
		domain    string
		mxRecords []*net.MX
		err       error
		expected  []string
	}{
		{"example.com", []*net.MX{{Host: "mail.example.com.", Pref: 10}}, nil, []string{"mail.example.com. 10\n"}},
		{"invalid-domain", nil, fmt.Errorf("no such host"), []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			mockLookup := MockHostLookup{
				LookupMXFunc: func(domain string) ([]*net.MX, error) {
					assert.Equal(t, tt.domain, domain)
					return tt.mxRecords, tt.err
				},
			}

			result := lookupMXRecords(mockLookup, tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupNSRecords(t *testing.T) {
	tests := []struct {
		domain    string
		nsRecords []*net.NS
		err       error
		expected  []string
	}{
		{"example.com", []*net.NS{{Host: "ns1.example.com."}}, nil, []string{"ns1.example.com."}},
		{"invalid-domain", nil, fmt.Errorf("no such host"), []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			mockLookup := MockHostLookup{
				LookupNSFunc: func(domain string) ([]*net.NS, error) {
					assert.Equal(t, tt.domain, domain)
					return tt.nsRecords, tt.err
				},
			}

			result := lookupNSRecords(mockLookup, tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupCNAMERecord(t *testing.T) {
	tests := []struct {
		domain   string
		cname    string
		err      error
		expected string
	}{
		{"www.example.com", "example.com.", nil, "example.com."},
		{"invalid-domain", "", fmt.Errorf("no such host"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			mockLookup := MockHostLookup{
				LookupCNAMEFunc: func(domain string) (string, error) {
					assert.Equal(t, tt.domain, domain)
					return tt.cname, tt.err
				},
			}

			result := lookupCNAMERecord(mockLookup, tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLookupTXTRecords(t *testing.T) {
	tests := []struct {
		domain     string
		txtRecords []string
		err        error
		expected   []string
	}{
		{"example.com", []string{"v=spf1 include:_spf.example.com ~all"}, nil, []string{"v=spf1 include:_spf.example.com ~all"}},
		{"invalid-domain", nil, fmt.Errorf("no such host"), []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			mockLookup := MockHostLookup{
				LookupTXTFunc: func(domain string) ([]string, error) {
					assert.Equal(t, tt.domain, domain)
					return tt.txtRecords, tt.err
				},
			}

			result := lookupTXTRecords(mockLookup, tt.domain)
			assert.Equal(t, tt.expected, result)
		})
	}
}
