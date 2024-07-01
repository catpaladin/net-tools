package cmd

import (
	"fmt"

	"github.com/catpaladin/net-tools/pkg/network"

	"github.com/spf13/cobra"
)

var (
	ipType string

	// ipCmd represents the ip command
	ipCmd = &cobra.Command{
		Use:   "ip",
		Short: "Used to get the public or private IP address of the host",
		Long:  "Used to get the public or private IP address of the host",
		Run: func(cmd *cobra.Command, args []string) {
			// Find IP
			ip, err := network.GetIP(ipType)
			if err != nil {
				fmt.Printf("%s Error getting %s IP: %v\n", errorMsg("[Error]"), ipType, err)
			} else {
				fmt.Printf("%s %s IP: %s\n", successMsg("[Success]"), ipType, dataMsg(ip))
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(ipCmd)

	ipCmd.PersistentFlags().StringVarP(&ipType, "type", "t", "public", "public or private")
}
