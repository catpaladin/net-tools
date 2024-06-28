package cmd

import (
	"fmt"

	"nt/pkg/network"

	"github.com/spf13/cobra"
)

var (
	host string
	port string

	// ncCmd represents the nc command
	ncCmd = &cobra.Command{
		Use:   "nc",
		Short: "Test host and port to see if open",
		Long:  "Test host and port to see if open",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Testing %s:%s\n", dataMsg(host), dataMsg(port))
			err := network.Netcat(host, port)
			if err != nil {
				fmt.Printf("%s Error connecting to %s:%s - %v\n", errorMsg("[Error]"), host, port, err)
			} else {
				fmt.Printf("%s Connection to %s:%s successful\n", successMsg("[Success]"), host, port)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(ncCmd)

	ncCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "host to test")
	ncCmd.PersistentFlags().StringVarP(&port, "port", "p", "80", "port to test")
}
