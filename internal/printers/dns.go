package printers

import (
	"fmt"
	"io"
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

// PrintDNSResult prints a DNS result with color formatting to the given writer
func PrintDNSResult(w io.Writer, result DNSResult) {
	// Print A Records
	if len(result.ARecords) > 0 {
		fmt.Fprintf(w, "%sA records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, ar := range result.ARecords {
			fmt.Fprintln(w, DataMsg(ar))
		}
	} else {
		fmt.Fprintf(w, "%sNo A records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Fprintln(w)

	// Print MX Records
	if len(result.MXRecords) > 0 {
		fmt.Fprintf(w, "%sMX records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, mx := range result.MXRecords {
			fmt.Fprintln(w, DataMsg(mx))
		}
	} else {
		fmt.Fprintf(w, "%sNo MX records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Fprintln(w)

	// Print NS Records
	if len(result.NSRecords) > 0 {
		fmt.Fprintf(w, "%sNS records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, ns := range result.NSRecords {
			fmt.Fprintln(w, DataMsg(ns))
		}
	} else {
		fmt.Fprintf(w, "%sNo NS records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Fprintln(w)

	// Print CNAME Record
	if result.CNAMERecord != "" {
		fmt.Fprintf(w, "%sCNAME record for %s:\n", SuccessMsg("[Success] "), result.Domain)
		fmt.Fprintln(w, DataMsg(result.CNAMERecord))
	} else {
		fmt.Fprintf(w, "%sNo CNAME found for %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Fprintln(w)

	// Print TXT Records
	if len(result.TXTRecords) > 0 {
		fmt.Fprintf(w, "%sTXT records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, tx := range result.TXTRecords {
			fmt.Fprintln(w, DataMsg(tx))
		}
	} else {
		fmt.Fprintf(w, "%sNo TXT records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
}
