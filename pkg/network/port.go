package network

import (
	"fmt"
	"net"
)

// Dialer is an interface that defines the Dial method.
type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

// NetDialer is a concrete implementation of Dialer using the net package.
type NetDialer struct{}

// Dial connects to the address on the named network.
func (d NetDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

// Netcat connects to the specified host and port using the provided Dialer.
func Netcat(host, port string) error {
	d := NetDialer{}
	err := netcatDialer(d, host, port)
	if err != nil {
		return err
	}

	return nil
}

func netcatDialer(dialer Dialer, host, port string) error {
	conn, err := dialer.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("error connecting: %v", err)
	}
	defer conn.Close()

	return nil
}
