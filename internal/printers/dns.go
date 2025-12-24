package printers

import "fmt"

// DNSResult represents the result of a DNS lookup
type DNSResult struct {
	Domain      string
	ARecords    []string
	MXRecords   []string
	NSRecords   []string
	CNAMERecord string
	TXTRecords  []string
}

// PrintDNSResult prints a DNS result with color formatting
func PrintDNSResult(result DNSResult) {
	// Print A Records
	if len(result.ARecords) > 0 {
		fmt.Printf("%sA records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, ar := range result.ARecords {
			fmt.Println(DataMsg(ar))
		}
	} else {
		fmt.Printf("%sNo A records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Println()

	// Print MX Records
	if len(result.MXRecords) > 0 {
		fmt.Printf("%sMX records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, mx := range result.MXRecords {
			fmt.Println(DataMsg(mx))
		}
	} else {
		fmt.Printf("%sNo MX records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Println()

	// Print NS Records
	if len(result.NSRecords) > 0 {
		fmt.Printf("%sNS records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, ns := range result.NSRecords {
			fmt.Println(DataMsg(ns))
		}
	} else {
		fmt.Printf("%sNo NS records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Println()

	// Print CNAME Record
	if result.CNAMERecord != "" {
		fmt.Printf("%sCNAME record for %s:\n", SuccessMsg("[Success] "), result.Domain)
		fmt.Println(DataMsg(result.CNAMERecord))
	} else {
		fmt.Printf("%sNo CNAME found for %s\n", WarningMsg("[Warning] "), result.Domain)
	}
	fmt.Println()

	// Print TXT Records
	if len(result.TXTRecords) > 0 {
		fmt.Printf("%sTXT records for %s:\n", SuccessMsg("[Success] "), result.Domain)
		for _, tx := range result.TXTRecords {
			fmt.Println(DataMsg(tx))
		}
	} else {
		fmt.Printf("%sNo TXT records found for: %s\n", WarningMsg("[Warning] "), result.Domain)
	}
}
