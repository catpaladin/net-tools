package cmd

import (
	"fmt"
	"os"

	"nt/pkg/network"

	"github.com/spf13/cobra"
)

// digCmd represents the dig command
var digCmd = &cobra.Command{
	Use:   "dig",
	Short: "Performs DNS lookups like dig",
	Long:  "Performs DNS lookups like dig",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: nt dig [domain]")
			os.Exit(1)
		}
		domain := args[0]
		network.Dig(domain)
	},
}

func init() {
	rootCmd.AddCommand(digCmd)
}
