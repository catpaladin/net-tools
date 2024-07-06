package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// GetIP returns the private or public IP address
func GetIP(ipType string) (string, error) {
	switch ipType {
	case "private":
		// Find private IP
		ni := RealNetworkInterface{}
		privateIP, err := getPrivateIP(ni)
		if err != nil {
			return "", fmt.Errorf("error getting private IP: %v", err)
		} else {
			return privateIP, nil
		}
	case "public":
		// Find external IP
		hc := RealHTTPClient{}
		externalIP, err := getExternalIP(hc)
		if err != nil {
			return "", fmt.Errorf("error getting external IP: %v", err)
		} else {
			return externalIP, nil
		}
	default:
		return "", fmt.Errorf("invalid IP type: %s", ipType)
	}
}

// NetworkInterface represents a network interface.
type NetworkInterface interface {
	Interfaces() ([]net.Interface, error)
	Addrs(iface net.Interface) ([]net.Addr, error)
}

// RealNetworkInterface is a concrete implementation of NetworkInterface using the net package.
type RealNetworkInterface struct{}

// Interfaces returns the list of system network interfaces.
func (r RealNetworkInterface) Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

// Addrs returns the list of addresses for a given network interface.
func (r RealNetworkInterface) Addrs(iface net.Interface) ([]net.Addr, error) {
	return iface.Addrs()
}

// getPrivateIP retrieves the first non-loopback, IPv4 address.
func getPrivateIP(netIf NetworkInterface) (string, error) {
	interfaces, err := netIf.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // interface down or loopback interface
		}
		addrs, err := netIf.Addrs(iface)
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no private IP address found")
}

// HTTPClient is an interface that defines the method for making HTTP GET requests.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// RealHTTPClient is a concrete implementation of HTTPClient using the net/http package.
type RealHTTPClient struct{}

// Get makes an HTTP GET request.
func (r RealHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// getExternalIP retrieves the external IP address by making an HTTP request.
func getExternalIP(client HTTPClient) (string, error) {
	response, err := client.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}
