package cmd

import (
	"github.com/catpaladin/net-tools/pkg/network"

	"github.com/spf13/cobra"
)

// netstatCmd represents the netstat command
var netstatCmd = &cobra.Command{
	Use:   "netstat",
	Short: "Performs a netstat on the host",
	Long:  "Performs a netstat on the host",
	Run: func(cmd *cobra.Command, args []string) {
		network.Netstat()
	},
}

func init() {
	rootCmd.AddCommand(netstatCmd)
}
