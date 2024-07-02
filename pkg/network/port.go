package network

import (
	"fmt"
	"net"
)

// Netcat tests a host and port to see if it is open
func Netcat(host string, port string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("error connecting: %v", err)
	}
	defer conn.Close()

	return nil
}
