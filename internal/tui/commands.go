package tui

import (
	"fmt"

	"github.com/catpaladin/net-tools/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
)

func (m MainModel) performDNS(domain string) tea.Cmd {
	return func() tea.Msg {
		dnsResult := performDNSLookup(domain)
		formattedResult := formatDNSResultPlain(domain, dnsResult)

		return dnsResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}

func (m MainModel) performPort(host, port string) tea.Cmd {
	return func() tea.Msg {
		if err := validatePort(port); err != nil {
			return portResultMsg{
				result: "",
				err:    err,
			}
		}

		err := network.Netcat(host, port)
		success := err == nil
		formattedResult := formatConnectionResultPlain(host, port, success, err)

		return portResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}

func (m MainModel) performIPLookup() tea.Cmd {
	return func() tea.Msg {
		var privateIP, publicIP string
		var privateErr, publicErr error

		switch m.ipModel.ipType {
		case "private":
			privateIP, privateErr = network.GetIP("private")
		case "public":
			publicIP, publicErr = network.GetIP("public")
		case "both":
			privateIP, privateErr = network.GetIP("private")
			publicIP, publicErr = network.GetIP("public")
		}

		formattedResult := formatIPResultPlain(m.ipModel.ipType, privateIP, publicIP, privateErr, publicErr)

		return ipResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}

func (m MainModel) performProcesses() tea.Cmd {
	contentWidth := m.contentWidth
	if contentWidth == 0 {
		contentWidth = 80
	}

	return func() tea.Msg {
		// Get results from network package
		results := network.Processes()

		// Convert results to our connection format
		var connections []ProcessesConnection
		for _, result := range results {
			connections = append(connections, ProcessesConnection{
				LocalAddr:  result.LocalAddr,
				Port:       result.LocalPort,
				PIDProgram: fmt.Sprintf("%s %d", result.Program, result.PID),
			})
		}

		// Format the results
		formattedResult := formatProcessesResult(connections, contentWidth)

		return processesResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}
