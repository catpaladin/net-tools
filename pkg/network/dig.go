package network

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

// Dig takes a domain and performs DNS queries
func Dig(domain string) {
	// Perform a DNS lookup for the A records
	addrs, err := net.LookupHost(domain)
	if err != nil {
		color.Yellow("No A records found for: %s\n", domain)
	} else {
		color.Green("A records for %s:\n", domain)
		for _, addr := range addrs {
			fmt.Println(addr)
		}
		fmt.Println()
	}

	// Perform a DNS lookup for the MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		color.Yellow("No MX records found for: %s\n", domain)
	} else {
		color.Green("\nMX records for %s:\n", domain)
		for _, mx := range mxRecords {
			color.Cyan("%s %d\n", mx.Host, mx.Pref)
		}
	}

	// Perform a DNS lookup for the NS records
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		color.Yellow("No NS records found for: %s\n", domain)
	} else {
		color.Green("\nNS records for %s:\n", domain)
		for _, ns := range nsRecords {
			color.Cyan(ns.Host)
		}
	}

	// Perform a DNS lookup for the CNAME record
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		color.Yellow("Error looking up CNAME: %v\n", err)
	} else {
		color.Green("\nCNAME for %s: %s\n", domain, cname)
	}

	// Perform a DNS lookup for the TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		color.Yellow("No TXT records for: %s\n", domain)
	} else {
		color.Green("\nTXT records for %s:\n", domain)
		for _, txt := range txtRecords {
			color.Cyan(txt)
		}
	}
}
