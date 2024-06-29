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
		privateIP, err := getPrivateIP()
		if err != nil {
			return "", fmt.Errorf("Error getting private IP: %v\n", err)
		} else {
			return privateIP, nil
		}
	case "public":
		// Find external IP
		externalIP, err := getExternalIP()
		if err != nil {
			return "", fmt.Errorf("Error getting external IP: %v\n", err)
		} else {
			return externalIP, nil
		}
	default:
		return "", fmt.Errorf("invalid IP type: %s", ipType)
	}
}

func getPrivateIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // interface down or loopback interface
		}
		addrs, err := iface.Addrs()
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

func getExternalIP() (string, error) {
	response, err := http.Get("http://checkip.amazonaws.com/")
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
