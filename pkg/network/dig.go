package network

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

// Dig takes a domain and performs DNS queries
func Dig(domain string) {
	nh := NetHostLookup{}

	// Perform a DNS lookup for the A records
	ars := lookupARecords(nh, domain)
	if len(ars) > 0 {
		color.Green("A records for %s:\n", domain)
		for _, ar := range ars {
			color.Cyan(ar)
		}
	} else {
		color.Yellow("No A records found for: %s\n", domain)
	}
	fmt.Println()

	// Perform a DNS lookup for the MX records
	mxs := lookupMXRecords(nh, domain)
	if len(mxs) > 0 {
		color.Green("MX records for %s:\n", domain)
		for _, mx := range mxs {
			color.Cyan(mx)
		}
	} else {
		color.Yellow("No MX records found for: %s\n", domain)
	}
	fmt.Println()

	// Perform a DNS lookup for the NS records
	nss := lookupNSRecords(nh, domain)
	if len(nss) > 0 {
		color.Green("NS records for %s:\n", domain)
		for _, ns := range nss {
			color.Cyan(ns)
		}
	} else {
		color.Yellow("No NS records found for: %s\n", domain)
	}
	fmt.Println()

	// Perform a DNS lookup for the CNAME record
	cn := lookupCNAMERecord(nh, domain)
	if cn != "" {
		color.Green("CNAME record for %s:\n", domain)
		color.Cyan(cn)
	} else {
		color.Yellow("No CNAME found for %s\n", domain)
	}
	fmt.Println()

	// Perform a DNS lookup for the TXT records
	txt := lookupTXTRecords(nh, domain)
	if len(txt) > 0 {
		color.Green("TXT records for %s:\n", domain)
		for _, tx := range txt {
			color.Cyan(tx)
		}
	} else {
		color.Yellow("No TXT records found for: %s\n", domain)
	}
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
	for _, txt := range txtRecords {
		output = append(output, txt)
	}
	return output
}
