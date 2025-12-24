package commands

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/catpaladin/net-tools/internal/utils"
	"github.com/catpaladin/net-tools/pkg/network"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// PortCmd creates and returns the port cobra command
func PortCmd() *cobra.Command {
	var host string
	var port string

	portCmd := &cobra.Command{
		Use:   "port",
		Short: "Tests if a host and port are open",
		Long:  "Tests network connectivity to a specific host and port to see if it's accepting connections.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				interactivePort(&host, &port)
			} else {
				host = args[0]
				port = args[1]
			}

			fmt.Printf("Testing %s:%s\n", utils.DataMsg(host), utils.DataMsg(port))
			err := network.Netcat(host, port)
			if err != nil {
				fmt.Printf("%s Error connecting to %s:%s - %v\n", utils.ErrorMsg("[Error]"), host, port, err)
			} else {
				fmt.Printf("%s Connection to %s:%s successful\n", utils.SuccessMsg("[Success]"), host, port)
			}
		},
	}

	return portCmd
}

func interactivePort(host, port *string) {
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
				Value(host),
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
				Value(port),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
}
