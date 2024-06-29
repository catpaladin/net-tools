package network

import (
	"fmt"
	"net"
)

// Dig takes a domain and performs DNS queries
func Dig(domain string) {
	// Perform a DNS lookup for the A records
	addrs, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("No A records found for: %s\n", domain)
	} else {
		fmt.Printf("A records for %s:\n", domain)
		for _, addr := range addrs {
			fmt.Println(addr)
		}
		fmt.Println()
	}

	// Perform a DNS lookup for the MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Printf("No MX records found for: %s\n", domain)
	} else {
		fmt.Printf("\nMX records for %s:\n", domain)
		for _, mx := range mxRecords {
			fmt.Printf("%s %d\n", mx.Host, mx.Pref)
		}
	}

	// Perform a DNS lookup for the NS records
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		fmt.Printf("No NS records found for: %s\n", domain)
	} else {
		fmt.Printf("\nNS records for %s:\n", domain)
		for _, ns := range nsRecords {
			fmt.Println(ns.Host)
		}
	}

	// Perform a DNS lookup for the CNAME record
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		fmt.Printf("Error looking up CNAME: %v\n", err)
	} else {
		fmt.Printf("\nCNAME for %s: %s\n", domain, cname)
	}

	// Perform a DNS lookup for the TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		fmt.Printf("No TXT records for: %s\n", domain)
	} else {
		fmt.Printf("\nTXT records for %s:\n", domain)
		for _, txt := range txtRecords {
			fmt.Println(txt)
		}
	}
}
