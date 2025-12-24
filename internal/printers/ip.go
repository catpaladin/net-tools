package printers

import (
	"fmt"
	"io"
)

// PrintIPResult prints the IP address information to the given writer
func PrintIPResult(w io.Writer, ipType string, ip string, err error) {
	if err != nil {
		fmt.Fprintf(w, "%s Error getting %s IP: %v\n", ErrorMsg("[Error]"), ipType, err)
	} else {
		fmt.Fprintf(w, "%s %s IP: %s\n", SuccessMsg("[Success]"), ipType, DataMsg(ip))
	}
}
