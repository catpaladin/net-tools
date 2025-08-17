package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateDigModel handles updates for the DNS lookup tab
func (m MainModel) updateDigModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Tab), key.Matches(msg, keys.Down), key.Matches(msg, keys.Up):
			// Allow viewport scrolling with up/down when there's content
			if m.digModel.result != "" {
				m.digModel.viewport, cmd = m.digModel.viewport.Update(msg)
				return cmd
			}
			return nil
		}
	}

	// Update text input if it's focused
	m.digModel.domainInput, cmd = m.digModel.domainInput.Update(msg)

	// Also update viewport for any other messages
	var viewportCmd tea.Cmd
	m.digModel.viewport, viewportCmd = m.digModel.viewport.Update(msg)

	return tea.Batch(cmd, viewportCmd)
}

// updateNetcatModel handles updates for the port test tab
func (m MainModel) updateNetcatModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Tab):
			m.ncModel.focused = (m.ncModel.focused + 1) % 2
			if m.ncModel.focused == 0 {
				m.ncModel.hostInput.Focus()
				m.ncModel.hostInput.PromptStyle = focusedPromptStyle
				m.ncModel.hostInput.TextStyle = focusedTextStyle
				m.ncModel.portInput.Blur()
				m.ncModel.portInput.PromptStyle = unfocusedPromptStyle
				m.ncModel.portInput.TextStyle = unfocusedTextStyle
			} else {
				m.ncModel.hostInput.Blur()
				m.ncModel.hostInput.PromptStyle = unfocusedPromptStyle
				m.ncModel.hostInput.TextStyle = unfocusedTextStyle
				m.ncModel.portInput.Focus()
				m.ncModel.portInput.PromptStyle = focusedPromptStyle
				m.ncModel.portInput.TextStyle = focusedTextStyle
			}
			return nil
		case key.Matches(msg, keys.Up), key.Matches(msg, keys.Down):
			// Allow viewport scrolling when there's content, otherwise handle input focus
			if m.ncModel.result != "" {
				m.ncModel.viewport, cmd = m.ncModel.viewport.Update(msg)
				return cmd
			}
			// Handle input focus navigation
			if key.Matches(msg, keys.Down) {
				m.ncModel.focused = (m.ncModel.focused + 1) % 2
			} else {
				m.ncModel.focused = (m.ncModel.focused - 1 + 2) % 2
			}
			if m.ncModel.focused == 0 {
				m.ncModel.hostInput.Focus()
				m.ncModel.hostInput.PromptStyle = focusedPromptStyle
				m.ncModel.hostInput.TextStyle = focusedTextStyle
				m.ncModel.portInput.Blur()
				m.ncModel.portInput.PromptStyle = unfocusedPromptStyle
				m.ncModel.portInput.TextStyle = unfocusedTextStyle
			} else {
				m.ncModel.hostInput.Blur()
				m.ncModel.hostInput.PromptStyle = unfocusedPromptStyle
				m.ncModel.hostInput.TextStyle = unfocusedTextStyle
				m.ncModel.portInput.Focus()
				m.ncModel.portInput.PromptStyle = focusedPromptStyle
				m.ncModel.portInput.TextStyle = focusedTextStyle
			}
			return nil
		}
	}

	var inputCmd tea.Cmd
	if m.ncModel.focused == 0 {
		m.ncModel.hostInput, inputCmd = m.ncModel.hostInput.Update(msg)
	} else {
		m.ncModel.portInput, inputCmd = m.ncModel.portInput.Update(msg)
	}

	// Also update viewport
	var viewportCmd tea.Cmd
	m.ncModel.viewport, viewportCmd = m.ncModel.viewport.Update(msg)

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

func (m MainModel) updateNetstatModel(msg tea.Msg) tea.Cmd {
	// Update viewport for scrolling
	var cmd tea.Cmd
	m.netstatModel.viewport, cmd = m.netstatModel.viewport.Update(msg)
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
