package commands

import (
	"errors"
	"log"
	"strings"

	"github.com/catpaladin/net-tools/internal/printers"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// DigCommand creates and returns the dig cobra command
func DigCommand() *cobra.Command {
	var domain string

	digCmd := &cobra.Command{
		Use:   "dig",
		Short: "Performs DNS lookups like dig",
		Long:  "Performs DNS lookups like dig",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				interactiveDig(&domain)
			} else {
				domain = args[0]
			}
			// Get raw results from pkg/network
			result := network.Dig(domain)
			// Format and print using internal printers
			printers.PrintDNSResult(printers.DNSResult(result))
		},
	}

	return digCmd
}

func interactiveDig(domain *string) {
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
