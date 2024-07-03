package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"

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
				interactiveNetcat()
			} else {
				host = args[0]
				port = args[1]
			}

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

func interactiveNetcat() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("IP to check:").
				Prompt("? ").
				Validate(func(str string) error {
					// autopass localhost
					if str == "localhost" {
						return nil
					}

					// test ipv4 and ipv6
					parsedIP := net.ParseIP(str)
					if parsedIP == nil {
						// secondary test for dns
						_, err := net.LookupHost(str)
						if err != nil {
							return errors.New("not a valid IP address")
						}
					}
					return nil
				}).
				Value(&host),
			huh.NewInput().
				Title("Port to check:").
				Prompt("? ").
				Validate(func(str string) error {
					portNum, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("not a valid port")
					}
					isValidPort := portNum > 0 && portNum <= 65535
					if !isValidPort {
						return errors.New("port provided is outside of valid port range")
					}
					return nil
				}).
				Value(&port),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
