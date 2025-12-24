package tui

import (
	"bytes"

	"github.com/catpaladin/net-tools/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
)

func (m MainModel) performDig(domain string) tea.Cmd {
	return func() tea.Msg {
		dnsResult := performDigLookup(domain)
		formattedResult := formatDNSResultPlain(domain, dnsResult)

		return digResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}

func (m MainModel) performNetcat(host, port string) tea.Cmd {
	return func() tea.Msg {
		if err := validatePort(port); err != nil {
			return netcatResultMsg{
				result: "",
				err:    err,
			}
		}

		err := network.Netcat(host, port)
		success := err == nil
		formattedResult := formatConnectionResultPlain(host, port, success, err)

		return netcatResultMsg{
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

func (m MainModel) performNetstat() tea.Cmd {
	contentWidth := m.contentWidth
	if contentWidth == 0 {
		contentWidth = 80
	}

	return func() tea.Msg {
		var buf bytes.Buffer
		network.Netstat()

		rawResult := buf.String()
		connections := parseNetstatOutput(rawResult)
		formattedResult := formatNetstatResultPlain(connections, contentWidth-8)

		return netstatResultMsg{
			result: formattedResult,
			err:    nil,
		}
	}
}
