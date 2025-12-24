package network

import (
	"fmt"
	"net"
)

// DNSResult represents the result of a DNS lookup
type DNSResult struct {
	Domain      string
	ARecords    []string
	MXRecords   []string
	NSRecords   []string
	CNAMERecord string
	TXTRecords  []string
}

// DNS takes a domain and performs DNS queries, returning raw results
func DNS(domain string) DNSResult {
	return DNSWithLookup(domain, NetHostLookup{})
}

// DNSWithLookup performs DNS queries using the provided HostLookup
func DNSWithLookup(domain string, lookup HostLookup) DNSResult {
	// Perform all DNS lookups
	result := DNSResult{
		Domain:      domain,
		ARecords:    lookupARecords(lookup, domain),
		MXRecords:   lookupMXRecords(lookup, domain),
		NSRecords:   lookupNSRecords(lookup, domain),
		CNAMERecord: lookupCNAMERecord(lookup, domain),
		TXTRecords:  lookupTXTRecords(lookup, domain),
	}

	return result
}

// HostLookup defines an interface for looking up hostnames.
type HostLookup interface {
	LookupHost(domain string) ([]string, error)
	LookupMX(domain string) ([]*net.MX, error)
	LookupNS(domain string) ([]*net.NS, error)
	LookupCNAME(domain string) (string, error)
	LookupTXT(domain string) ([]string, error)
}

// NetHostLookup is a concrete implementation of HostLookup using the net package.
type NetHostLookup struct{}

// LookupHost looks up the hostnames using net.LookupHost.
func (n NetHostLookup) LookupHost(domain string) ([]string, error) {
	return net.LookupHost(domain)
}

func lookupARecords(lookup HostLookup, domain string) []string {
	var output []string
	addrs, err := lookup.LookupHost(domain)
	if err != nil {
		return []string{}
	}
	return append(output, addrs...)
}

// LookupMX looks up the MX records using net.LookupMX.
func (n NetHostLookup) LookupMX(domain string) ([]*net.MX, error) {
	return net.LookupMX(domain)
}

func lookupMXRecords(lookup HostLookup, domain string) []string {
	var output []string
	mxRecords, err := lookup.LookupMX(domain)
	if err != nil {
		return []string{}
	}
	for _, mx := range mxRecords {
		output = append(output, fmt.Sprintf("%s %d\n", mx.Host, mx.Pref))
	}
	return output
}

// LookupNS looks up the NS records using net.LookupNS.
func (n NetHostLookup) LookupNS(domain string) ([]*net.NS, error) {
	return net.LookupNS(domain)
}

func lookupNSRecords(lookup HostLookup, domain string) []string {
	var output []string
	nsRecords, err := lookup.LookupNS(domain)
	if err != nil {
		return []string{}
	}
	for _, ns := range nsRecords {
		output = append(output, ns.Host)
	}
	return output
}

// LookupCNAME looks up the CNAME record using net.LookupCNAME.
func (n NetHostLookup) LookupCNAME(domain string) (string, error) {
	return net.LookupCNAME(domain)
}

func lookupCNAMERecord(lookup HostLookup, domain string) string {
	cname, err := lookup.LookupCNAME(domain)
	if err != nil {
		return ""
	}
	return cname
}

// LookupTXT looks up the TXT records using net.LookupTXT.
func (n NetHostLookup) LookupTXT(domain string) ([]string, error) {
	return net.LookupTXT(domain)
}

func lookupTXTRecords(lookup HostLookup, domain string) []string {
	var output []string
	txtRecords, err := lookup.LookupTXT(domain)
	if err != nil {
		return []string{}
	}
	return append(output, txtRecords...)
}
