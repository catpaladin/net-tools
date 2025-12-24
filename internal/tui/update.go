package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateDNSModel handles updates for the DNS lookup tab
func (m MainModel) updateDNSModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Tab), key.Matches(msg, keys.Down), key.Matches(msg, keys.Up):
			// Allow viewport scrolling with up/down when there's content
			if m.dnsModel.result != "" {
				m.dnsModel.viewport, cmd = m.dnsModel.viewport.Update(msg)
				return cmd
			}
			return nil
		}
	}

	// Update text input if it's focused
	m.dnsModel.domainInput, cmd = m.dnsModel.domainInput.Update(msg)

	// Also update viewport for any other messages
	var viewportCmd tea.Cmd
	m.dnsModel.viewport, viewportCmd = m.dnsModel.viewport.Update(msg)

	return tea.Batch(cmd, viewportCmd)
}

// updatePortModel handles updates for the port test tab
func (m MainModel) updatePortModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Tab):
			m.portModel.focused = (m.portModel.focused + 1) % 2
			if m.portModel.focused == 0 {
				m.portModel.hostInput.Focus()
				m.portModel.hostInput.PromptStyle = focusedPromptStyle
				m.portModel.hostInput.TextStyle = focusedTextStyle
				m.portModel.portInput.Blur()
				m.portModel.portInput.PromptStyle = unfocusedPromptStyle
				m.portModel.portInput.TextStyle = unfocusedTextStyle
			} else {
				m.portModel.hostInput.Blur()
				m.portModel.hostInput.PromptStyle = unfocusedPromptStyle
				m.portModel.hostInput.TextStyle = unfocusedTextStyle
				m.portModel.portInput.Focus()
				m.portModel.portInput.PromptStyle = focusedPromptStyle
				m.portModel.portInput.TextStyle = focusedTextStyle
			}
			return nil
		case key.Matches(msg, keys.Up), key.Matches(msg, keys.Down):
			// Allow viewport scrolling when there's content, otherwise handle input focus
			if m.portModel.result != "" {
				m.portModel.viewport, cmd = m.portModel.viewport.Update(msg)
				return cmd
			}
			// Handle input focus navigation
			if key.Matches(msg, keys.Down) {
				m.portModel.focused = (m.portModel.focused + 1) % 2
			} else {
				m.portModel.focused = (m.portModel.focused - 1 + 2) % 2
			}
			if m.portModel.focused == 0 {
				m.portModel.hostInput.Focus()
				m.portModel.hostInput.PromptStyle = focusedPromptStyle
				m.portModel.hostInput.TextStyle = focusedTextStyle
				m.portModel.portInput.Blur()
				m.portModel.portInput.PromptStyle = unfocusedPromptStyle
				m.portModel.portInput.TextStyle = unfocusedTextStyle
			} else {
				m.portModel.hostInput.Blur()
				m.portModel.hostInput.PromptStyle = unfocusedPromptStyle
				m.portModel.hostInput.TextStyle = unfocusedTextStyle
				m.portModel.portInput.Focus()
				m.portModel.portInput.PromptStyle = focusedPromptStyle
				m.portModel.portInput.TextStyle = focusedTextStyle
			}
			return nil
		}
	}

	var inputCmd tea.Cmd
	if m.portModel.focused == 0 {
		m.portModel.hostInput, inputCmd = m.portModel.hostInput.Update(msg)
	} else {
		m.portModel.portInput, inputCmd = m.portModel.portInput.Update(msg)
	}

	// Also update viewport
	var viewportCmd tea.Cmd
	m.portModel.viewport, viewportCmd = m.portModel.viewport.Update(msg)

	return tea.Batch(inputCmd, viewportCmd)
}

// updateIPModel handles updates for the IP address tab
func (m MainModel) updateIPModel(msg tea.Msg) tea.Cmd {
	if m.activeTab != 2 {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			// If there's content, scroll viewport, otherwise change IP type
			if m.ipModel.result != "" {
				var cmd tea.Cmd
				m.ipModel.viewport, cmd = m.ipModel.viewport.Update(msg)
				return cmd
			}
			if m.ipModel.selector > 0 {
				m.ipModel.selector--
				m.updateIPType()
			}
		case key.Matches(msg, keys.Down):
			// If there's content, scroll viewport, otherwise change IP type
			if m.ipModel.result != "" {
				var cmd tea.Cmd
				m.ipModel.viewport, cmd = m.ipModel.viewport.Update(msg)
				return cmd
			}
			if m.ipModel.selector < 2 {
				m.ipModel.selector++
				m.updateIPType()
			}
		case key.Matches(msg, keys.Tab):
			if m.ipModel.selector < 2 {
				m.ipModel.selector++
			} else {
				m.ipModel.selector = 0
			}
			m.updateIPType()
		}
	}

	// Update viewport for other messages
	var cmd tea.Cmd
	m.ipModel.viewport, cmd = m.ipModel.viewport.Update(msg)
	return cmd
}

func (m MainModel) updateProcessesModel(msg tea.Msg) tea.Cmd {
	// Update viewport for scrolling
	var cmd tea.Cmd
	m.processesModel.viewport, cmd = m.processesModel.viewport.Update(msg)
	return cmd
}

func (m *MainModel) updateIPType() {
	switch m.ipModel.selector {
	case 0:
		m.ipModel.ipType = "private"
	case 1:
		m.ipModel.ipType = "public"
	case 2:
		m.ipModel.ipType = "both"
	}
}

func validatePort(portStr string) error {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port number")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}
