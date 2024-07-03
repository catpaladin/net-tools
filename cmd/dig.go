package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"

	"github.com/spf13/cobra"
)

var (
	domain string

	// digCmd represents the dig command
	digCmd = &cobra.Command{
		Use:   "dig",
		Short: "Performs DNS lookups like dig",
		Long:  "Performs DNS lookups like dig",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				interactiveDig()
			} else {
				domain = args[0]
			}
			network.Dig(domain)
		},
	}
)

func init() {
	rootCmd.AddCommand(digCmd)
}

func interactiveDig() {
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
				Value(&domain),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
