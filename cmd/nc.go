package cmd

import (
	"fmt"
	"os"

	"github.com/catpaladin/net-tools/pkg/network"

	"github.com/spf13/cobra"
)

var (
	host string
	port string

	// ncCmd represents the nc command
	ncCmd = &cobra.Command{
		Use:   "nc",
		Short: "Netcat subcommand to test host and port to see if open",
		Long:  "Netcat subcommand to test host and port to see if open",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("Usage: nt nc [host] [port]")
				os.Exit(1)
			}

			host = args[0]
			port = args[1]

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
}
