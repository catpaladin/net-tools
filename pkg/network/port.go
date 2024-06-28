package network

import (
	"fmt"
	"net"
)

func Netcat(host string, port string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("Error connecting: %v\n", err)
	}
	defer conn.Close()

	return nil
}
