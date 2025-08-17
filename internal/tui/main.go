package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NewModel creates a new TUI main model
func NewModel() MainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Initialize all tab models
	digModel := &DigModel{
		domainInput: newTextInput("example.com", "Enter domain to lookup..."),
		viewport:    viewport.New(80, 20),
	}
	// Focus the domain input for the initial tab
	digModel.domainInput.Focus()
	digModel.domainInput.PromptStyle = focusedPromptStyle
	digModel.domainInput.TextStyle = focusedTextStyle

	ncModel := &NetcatModel{
		hostInput: newTextInput("localhost", "Enter host..."),
		portInput: newTextInput("80", "Enter port..."),
		focused:   0, // Start with host input focused
		viewport:  viewport.New(80, 20),
	}
	// Set up initial focus states for netcat model
	ncModel.hostInput.Focus()
	ncModel.hostInput.PromptStyle = focusedPromptStyle
	ncModel.hostInput.TextStyle = focusedTextStyle
	ncModel.portInput.Blur()
	ncModel.portInput.PromptStyle = unfocusedPromptStyle
	ncModel.portInput.TextStyle = unfocusedTextStyle

	ipModel := &IPModel{
		ipType:   "both",
		selector: 2, // Default to "Both" option
		viewport: viewport.New(80, 20),
	}

	netstatModel := &NetstatModel{
		viewport: viewport.New(80, 20),
	}

	return MainModel{
		tabs: []Tab{
			{Name: "DNS Lookup", Active: true},
			{Name: "Port Test", Active: false},
			{Name: "IP Address", Active: false},
			{Name: "Netstat", Active: false},
		},
		activeTab:    0,
		digModel:     digModel,
		ncModel:      ncModel,
		ipModel:      ipModel,
		netstatModel: netstatModel,
		spinner:      s,
	}
}

func newTextInput(placeholder, prompt string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = "> "
	ti.Blur()
	ti.CharLimit = 256
	ti.Width = 40

	// Configure proper text input styling for BubbleTea v1.x
	ti.PromptStyle = unfocusedPromptStyle
	ti.TextStyle = unfocusedTextStyle
	ti.PlaceholderStyle = unfocusedTextStyle
	ti.Cursor.Style = cursorStyle

	return ti
}

// Init initializes the TUI
func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
	)
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.contentWidth = msg.Width - 6 // Match View() calculation
		if m.contentWidth < 40 {
			m.contentWidth = 40
		}

		// Calculate space allocation more precisely
		titleHeight := 3   // Title + 2 newlines
		tabBarHeight := 2  // Tab bar + 1 newline
		helpHeight := 1    // Help line
		paddingHeight := 6 // Content box padding (3 top + 3 bottom)
		headerHeight := 5  // Space for tab header and inputs

		totalOverhead := titleHeight + tabBarHeight + helpHeight + paddingHeight + headerHeight
		availableContentHeight := m.height - totalOverhead

		if availableContentHeight < 3 {
			availableContentHeight = 3
		}

		m.digModel.viewport.Width = m.contentWidth - 8 // Account for content padding
		m.digModel.viewport.Height = availableContentHeight
		m.ncModel.viewport.Width = m.contentWidth - 8
		m.ncModel.viewport.Height = availableContentHeight
		m.ipModel.viewport.Width = m.contentWidth - 8
		m.ipModel.viewport.Height = availableContentHeight
		m.netstatModel.viewport.Width = m.contentWidth - 8
		m.netstatModel.viewport.Height = availableContentHeight

		m.ready = true

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Left):
			if m.activeTab > 0 {
				m.tabs[m.activeTab].Active = false
				m.activeTab--
				m.tabs[m.activeTab].Active = true
				// Clear previous content when switching tabs
				m.clearTabContent()
			}

		case key.Matches(msg, keys.Right):
			if m.activeTab < len(m.tabs)-1 {
				m.tabs[m.activeTab].Active = false
				m.activeTab++
				m.tabs[m.activeTab].Active = true
				// Clear previous content when switching tabs
				m.clearTabContent()
			}

		case key.Matches(msg, keys.Tab):
			// Tab behavior depends on the active tab
			// For IP tab (2), it cycles through IP types
			// For other tabs, it handles field navigation

		case key.Matches(msg, keys.Escape):
			// Escape can be used for additional functionality if needed
			// Currently not used for tab switching

		case key.Matches(msg, keys.Enter):
			return m.handleEnter()
		}

	case digResultMsg:
		m.digModel.loading = false
		if msg.err != nil {
			m.digModel.error = msg.err.Error()
		} else {
			m.digModel.result = msg.result
			m.digModel.error = ""
			m.digModel.viewport.SetContent(msg.result)
		}

	case netcatResultMsg:
		m.ncModel.loading = false
		if msg.err != nil {
			m.ncModel.error = msg.err.Error()
		} else {
			m.ncModel.result = msg.result
			m.ncModel.error = ""
			m.ncModel.viewport.SetContent(msg.result)
		}

	case ipResultMsg:
		m.ipModel.loading = false
		if msg.err != nil {
			m.ipModel.error = msg.err.Error()
		} else {
			m.ipModel.result = msg.result
			m.ipModel.error = ""
			m.ipModel.viewport.SetContent(msg.result)
		}

	case netstatResultMsg:
		m.netstatModel.loading = false
		if msg.err != nil {
			m.netstatModel.error = msg.err.Error()
		} else {
			m.netstatModel.result = msg.result
			m.netstatModel.error = ""
			m.netstatModel.viewport.SetContent(msg.result)
		}

	}

	// Update spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update active tab model
	switch m.activeTab {
	case 0: // DNS Lookup
		cmd = m.updateDigModel(msg)
	case 1: // Port Test
		cmd = m.updateNetcatModel(msg)
	case 2: // IP Address
		cmd = m.updateIPModel(msg)
	case 3: // Netstat
		cmd = m.updateNetstatModel(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m MainModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Title section - keep styled
	title := titleStyle.Width(m.width - 4).Render("NET-TOOLS - Network Diagnostics Suite")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Tab bar section - keep styled
	tabs := make([]string, len(m.tabs))
	for i, tab := range m.tabs {
		if tab.Active {
			tabs[i] = activeTabStyle.Render(tab.Name)
		} else {
			tabs[i] = inactiveTabStyle.Render(tab.Name)
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	b.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tabBar))
	b.WriteString("\n")

	// Content section - use simple container with minimal styling
	var content string
	switch m.activeTab {
	case 0:
		content = m.renderDigTab(m.contentWidth)
	case 1:
		content = m.renderNetcatTab(m.contentWidth)
	case 2:
		content = m.renderIPTab(m.contentWidth)
	case 3:
		content = m.renderNetstatTab(m.contentWidth)
	}

	// Simple content container that fits the available space
	contentContainer := lipgloss.NewStyle().
		Width(m.contentWidth).
		Padding(2, 3).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(content)

	b.WriteString(contentContainer)
	b.WriteString("\n")

	// Help section - keep styled
	var help string
	if m.activeTab == 2 {
		help = "Use ← → to switch tabs • ↑ ↓ or Tab to cycle IP types • Enter to execute • q to quit"
	} else {
		help = "Use ← → to switch tabs • Enter to execute • Tab to navigate fields • q to quit"
	}
	b.WriteString(helpStyle.Width(m.width).Render(help))

	return b.String()
}

func (m MainModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case 0: // DNS Lookup
		if m.digModel.domainInput.Value() != "" {
			m.digModel.loading = true
			return m, m.performDig(m.digModel.domainInput.Value())
		}
	case 1: // Port Test
		if m.ncModel.hostInput.Value() != "" && m.ncModel.portInput.Value() != "" {
			m.ncModel.loading = true
			return m, m.performNetcat(m.ncModel.hostInput.Value(), m.ncModel.portInput.Value())
		}
	case 2: // IP Address
		m.ipModel.loading = true
		return m, m.performIPLookup()
	case 3: // Netstat
		m.netstatModel.loading = true
		return m, m.performNetstat()
	}
	return m, nil
}

func (m MainModel) renderDigTab(contentWidth int) string {
	var content strings.Builder

	content.WriteString("🔍 DNS Lookup\n\n")

	content.WriteString("Domain:\n")
	content.WriteString(m.digModel.domainInput.View() + "\n\n")

	if m.digModel.loading {
		content.WriteString(fmt.Sprintf("%s Looking up DNS records...", m.spinner.View()))
	} else if m.digModel.error != "" {
		content.WriteString("❌ Error: " + m.digModel.error)
	} else if m.digModel.result != "" {
		content.WriteString("✅ DNS Records:\n\n")
		content.WriteString(m.digModel.viewport.View())
	} else {
		content.WriteString("Enter a domain name and press Enter to perform DNS lookup")
	}

	return content.String()
}

func (m MainModel) renderNetcatTab(contentWidth int) string {
	var content strings.Builder

	content.WriteString("🔌 Port Connectivity Test\n\n")

	content.WriteString("Host:\n")
	content.WriteString(m.ncModel.hostInput.View() + "\n\n")

	content.WriteString("Port:\n")
	content.WriteString(m.ncModel.portInput.View() + "\n\n")

	if m.ncModel.loading {
		content.WriteString(fmt.Sprintf("%s Testing connection...", m.spinner.View()))
	} else if m.ncModel.error != "" {
		content.WriteString("❌ Error: " + m.ncModel.error)
	} else if m.ncModel.result != "" {
		content.WriteString("✅ Connection Test:\n\n")
		content.WriteString(m.ncModel.viewport.View())
	} else {
		content.WriteString("Enter host and port, then press Enter to test connectivity")
	}

	return content.String()
}

func (m MainModel) renderIPTab(contentWidth int) string {
	var content strings.Builder

	content.WriteString("🌐 IP Address Information\n\n")

	content.WriteString("IP Type:\n")
	currentType := m.getIPTypeDisplayName()
	content.WriteString(fmt.Sprintf("[%s]\n\n", currentType))

	if m.ipModel.loading {
		content.WriteString(fmt.Sprintf("%s Getting IP address...", m.spinner.View()))
	} else if m.ipModel.error != "" {
		content.WriteString("❌ Error: " + m.ipModel.error)
	} else if m.ipModel.result != "" {
		content.WriteString("✅ IP Information:\n\n")
		content.WriteString(m.ipModel.viewport.View())
	} else {
		content.WriteString("Press Enter to get IP address")
	}

	return content.String()
}

func (m MainModel) renderNetstatTab(contentWidth int) string {
	var content strings.Builder

	content.WriteString("📊 Network Connections\n\n")

	if m.netstatModel.loading {
		content.WriteString(fmt.Sprintf("%s Gathering network information...", m.spinner.View()))
	} else if m.netstatModel.error != "" {
		content.WriteString("❌ Error: " + m.netstatModel.error)
	} else if m.netstatModel.result != "" {
		content.WriteString("✅ Active Connections:\n\n")
		content.WriteString(m.netstatModel.viewport.View())
	} else {
		content.WriteString("Press Enter to scan for active network connections")
	}

	return content.String()
}

func (m *MainModel) initializeTabState() {
	m.digModel.result = ""
	m.digModel.error = ""
	m.digModel.loading = false
	m.digModel.domainInput.SetValue("example.com")

	m.ncModel.result = ""
	m.ncModel.error = ""
	m.ncModel.loading = false
	m.ncModel.hostInput.SetValue("localhost")
	m.ncModel.portInput.SetValue("80")
	m.ncModel.focused = 0

	m.ipModel.result = ""
	m.ipModel.error = ""
	m.ipModel.loading = false
	m.ipModel.selector = 2
	m.ipModel.ipType = "both"

	m.netstatModel.result = ""
	m.netstatModel.error = ""
	m.netstatModel.loading = false

	m.setTabFocus()
}

func (m *MainModel) setTabFocus() {
	m.digModel.domainInput.Blur()
	m.digModel.domainInput.PromptStyle = unfocusedPromptStyle
	m.digModel.domainInput.TextStyle = unfocusedTextStyle

	m.ncModel.hostInput.Blur()
	m.ncModel.hostInput.PromptStyle = unfocusedPromptStyle
	m.ncModel.hostInput.TextStyle = unfocusedTextStyle
	m.ncModel.portInput.Blur()
	m.ncModel.portInput.PromptStyle = unfocusedPromptStyle
	m.ncModel.portInput.TextStyle = unfocusedTextStyle

	switch m.activeTab {
	case 0: // DNS Lookup
		m.digModel.domainInput.Focus()
		m.digModel.domainInput.PromptStyle = focusedPromptStyle
		m.digModel.domainInput.TextStyle = focusedTextStyle
	case 1: // Port Test
		m.ncModel.hostInput.Focus()
		m.ncModel.hostInput.PromptStyle = focusedPromptStyle
		m.ncModel.hostInput.TextStyle = focusedTextStyle
		m.ncModel.focused = 0
		// IP Address and Netstat tabs don't have text inputs to focus
	}
}

func (m *MainModel) clearTabContent() {
	m.digModel.result = ""
	m.digModel.error = ""
	m.digModel.loading = false

	m.ncModel.result = ""
	m.ncModel.error = ""
	m.ncModel.loading = false

	m.ipModel.result = ""
	m.ipModel.error = ""
	m.ipModel.loading = false
	m.ipModel.selector = 2
	m.ipModel.ipType = "both"

	m.netstatModel.result = ""
	m.netstatModel.error = ""
	m.netstatModel.loading = false
}

func (m MainModel) getIPTypeDisplayName() string {
	switch m.ipModel.selector {
	case 0:
		return "Private"
	case 1:
		return "Public"
	case 2:
		return "Both"
	default:
		return "Both"
	}
}
