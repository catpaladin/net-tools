package commands

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/catpaladin/net-tools/internal/printers"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// DNSCommand creates and returns the dns cobra command
func DNSCommand() *cobra.Command {
	var domain string

	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "Performs DNS lookups for a given domain",
		Long:  "Performs basic DNS lookups including A, MX, NS, CNAME, and TXT records.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				interactiveDNS(&domain)
			} else {
				domain = args[0]
			}
			// Get raw results from pkg/network
			result := network.DNS(domain)
			// Format and print using internal printers
			printers.PrintDNSResult(os.Stdout, printers.DNSResult(result))
		},
	}

	return dnsCmd
}

func interactiveDNS(domain *string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Domain/Subdomain:").
				Prompt("? ").
				Validate(func(str string) error {
					if !strings.Contains(str, ".") {
						return errors.New("domains should have a '.' in them")
					}
					return nil
				}).
				Value(domain),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
