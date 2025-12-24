package commands

import (
	"log"
	"os"

	"github.com/catpaladin/net-tools/internal/printers"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// IPCmd creates and returns the ip cobra command
func IPCmd() *cobra.Command {
	var ipType string

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "Used to get the public or private IP address of the host",
		Long:  "Used to get the public or private IP address of the host",
		Run: func(cmd *cobra.Command, args []string) {
			if ipType == "" {
				interactiveIP(&ipType)
			}
			// Find IP
			switch ipType {
			case "both":
				privateIP, err := network.GetIP("private")
				printers.PrintIPResult(os.Stdout, "Private", privateIP, err)

				publicIP, err := network.GetIP("public")
				printers.PrintIPResult(os.Stdout, "Public", publicIP, err)
			default:
				ip, err := network.GetIP(ipType)
				printers.PrintIPResult(os.Stdout, ipType, ip, err)
			}
		},
	}

	ipCmd.PersistentFlags().StringVarP(&ipType, "type", "t", "", "public|private|both")

	return ipCmd
}

func interactiveIP(ipType *string) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("IP Type:").
				Options(
					huh.NewOption("Both", "both"),
					huh.NewOption("Private", "private"),
					huh.NewOption("Public", "public"),
				).
				Value(ipType),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
