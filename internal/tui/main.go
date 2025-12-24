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
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	// Initialize all tab models
	dnsModel := &DNSModel{
		domainInput: newTextInput("example.com", "Enter domain to lookup..."),
		viewport:    viewport.New(80, 20),
	}
	// Focus the domain input for the initial tab
	dnsModel.domainInput.Focus()
	dnsModel.domainInput.PromptStyle = focusedPromptStyle
	dnsModel.domainInput.TextStyle = focusedTextStyle

	portModel := &PortModel{
		hostInput: newTextInput("localhost", "Enter host..."),
		portInput: newTextInput("80", "Enter port..."),
		focused:   0, // Start with host input focused
		viewport:  viewport.New(80, 20),
	}
	// Set up initial focus states for port model
	portModel.hostInput.Focus()
	portModel.hostInput.PromptStyle = focusedPromptStyle
	portModel.hostInput.TextStyle = focusedTextStyle
	portModel.portInput.Blur()
	portModel.portInput.PromptStyle = unfocusedPromptStyle
	portModel.portInput.TextStyle = unfocusedTextStyle

	ipModel := &IPModel{
		ipType:   "both",
		selector: 2, // Default to "Both" option
		viewport: viewport.New(80, 20),
	}

	processesModel := &ProcessesModel{
		viewport: viewport.New(80, 20),
	}

	return MainModel{
		tabs: []Tab{
			{Name: "DNS Lookup", Active: true},
			{Name: "Port Test", Active: false},
			{Name: "IP Address", Active: false},
			{Name: "Processes", Active: false},
		},
		activeTab:      0,
		dnsModel:       dnsModel,
		portModel:      portModel,
		ipModel:        ipModel,
		processesModel: processesModel,
		spinner:        s,
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
	// Use backgrounds that match the content container (234)
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

		// Content width is slightly less than window width for margins
		m.contentWidth = msg.Width - 6
		if m.contentWidth < 60 {
			m.contentWidth = 60
		}

		// Header area: Title (2) + Tabs (2) + Spacing (2) = 6
		// Footer area: Help (2) = 2
		// Padding/Margins: Content Padding (2) + Margins (2) = 4
		overhead := 12

		availableHeight := m.height - overhead
		if availableHeight < 10 {
			availableHeight = 10
		}

		// Adjust viewport sizes
		viewWidth := m.contentWidth - 4 // Account for container padding

		m.dnsModel.viewport.Width = viewWidth
		m.dnsModel.viewport.Height = availableHeight
		m.portModel.viewport.Width = viewWidth
		m.portModel.viewport.Height = availableHeight
		m.ipModel.viewport.Width = viewWidth
		m.ipModel.viewport.Height = availableHeight
		m.processesModel.viewport.Width = viewWidth
		m.processesModel.viewport.Height = availableHeight

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
				m.setTabFocus()
			}

		case key.Matches(msg, keys.Right):
			if m.activeTab < len(m.tabs)-1 {
				m.tabs[m.activeTab].Active = false
				m.activeTab++
				m.tabs[m.activeTab].Active = true
				m.setTabFocus()
			}

		case key.Matches(msg, keys.Tab):
			// Tab behavior depends on the active tab
			// For IP tab (2), it cycles through IP types
			// For other tabs, it handles field navigation

		case key.Matches(msg, keys.Clear):
			m.clearActiveTabContent()

		case key.Matches(msg, keys.Enter):
			return m.handleEnter()
		}

	case dnsResultMsg:
		m.dnsModel.loading = false
		if msg.err != nil {
			m.dnsModel.error = msg.err.Error()
		} else {
			m.dnsModel.result = msg.result
			m.dnsModel.error = ""
			m.dnsModel.viewport.SetContent(msg.result)
		}

	case portResultMsg:
		m.portModel.loading = false
		if msg.err != nil {
			m.portModel.error = msg.err.Error()
		} else {
			m.portModel.result = msg.result
			m.portModel.error = ""
			m.portModel.viewport.SetContent(msg.result)
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

	case processesResultMsg:
		m.processesModel.loading = false
		if msg.err != nil {
			m.processesModel.error = msg.err.Error()
		} else {
			m.processesModel.result = msg.result
			m.processesModel.error = ""
			m.processesModel.viewport.SetContent(msg.result)
		}

	}

	// Update spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update active tab model
	switch m.activeTab {
	case 0: // DNS Lookup
		cmd = m.updateDNSModel(msg)
	case 1: // Port Test
		cmd = m.updatePortModel(msg)
	case 2: // IP Address
		cmd = m.updateIPModel(msg)
	case 3: // Processes
		cmd = m.updateProcessesModel(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m MainModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Render Title
	title := titleStyle.Width(m.contentWidth).Render("NET-TOOLS")

	// Render Tabs
	tabs := make([]string, len(m.tabs))
	for i, tab := range m.tabs {
		if tab.Active {
			tabs[i] = activeTabStyle.Render(tab.Name)
		} else {
			tabs[i] = inactiveTabStyle.Render(tab.Name)
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Render Content
	var content string
	switch m.activeTab {
	case 0:
		content = m.renderDNSTab()
	case 1:
		content = m.renderPortTab()
	case 2:
		content = m.renderIPTab()
	case 3:
		content = m.renderProcessesTab()
	}

	contentBox := tabOuterStyle.
		Width(m.contentWidth).
		Height(m.height - 14). // More cushion for scaling
		Render(content)

	// Render Help
	help := m.renderHelp()

	// Join all parts vertically and center them
	ui := lipgloss.JoinVertical(lipgloss.Center,
		title,
		tabBar,
		contentBox,
		help,
	)

	// Fill the screen and center the UI
	view := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		ui,
		lipgloss.WithWhitespaceBackground(bgColor),
	)

	return appStyle.Render(view)
}

func (m MainModel) renderHelp() string {
	keys := []struct {
		key  string
		desc string
	}{
		{"←/→", "Tabs"},
		{"Tab", "Nav"},
		{"Enter", "Run"},
		{"Esc", "Clear"},
		{"q", "Quit"},
	}

	var helpParts []string
	for _, k := range keys {
		part := fmt.Sprintf("%s %s", helpKeyStyle.Render(k.key), k.desc)
		helpParts = append(helpParts, part)
	}

	return helpStyle.Width(m.width).Render(strings.Join(helpParts, "  •  "))
}

func (m MainModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case 0: // DNS Lookup
		if m.dnsModel.domainInput.Value() != "" {
			m.dnsModel.loading = true
			return m, m.performDNS(m.dnsModel.domainInput.Value())
		}
	case 1: // Port Test
		if m.portModel.hostInput.Value() != "" && m.portModel.portInput.Value() != "" {
			m.portModel.loading = true
			return m, m.performPort(m.portModel.hostInput.Value(), m.portModel.portInput.Value())
		}
	case 2: // IP Address
		m.ipModel.loading = true
		return m, m.performIPLookup()
	case 3: // Processes
		m.processesModel.loading = true
		return m, m.performProcesses()
	}
	return m, nil
}

func (m MainModel) renderDNSTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Render("🔍 DNS LOOKUP"))
	content.WriteString("\n\n")

	// Domain Input Row - Wrap in a style that fills background
	inputRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Domain: "),
		m.dnsModel.domainInput.View(),
	)
	content.WriteString(inputRow + "\n\n")

	if m.dnsModel.loading {
		content.WriteString(loadingBlockStyle.Render(" BUSY ") + " Looking up records...")
	} else if m.dnsModel.error != "" {
		content.WriteString(errorBlockStyle.Render(" ERROR ") + " " + m.dnsModel.error)
	} else if m.dnsModel.result != "" {
		content.WriteString(successBlockStyle.Render(" DONE ") + " DNS Records Found:\n\n")
		content.WriteString(m.dnsModel.viewport.View())
	} else {
		content.WriteString(infoStyle.Render("Enter a domain name and press Enter"))
	}

	return content.String()
}

func (m MainModel) renderPortTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Render("🔌 PORT CONNECTIVITY"))
	content.WriteString("\n\n")

	// Host Input Row
	content.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Host: "),
		m.portModel.hostInput.View(),
	))
	content.WriteString("\n")

	// Port Input Row
	content.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Port: "),
		m.portModel.portInput.View(),
	))
	content.WriteString("\n\n")

	if m.portModel.loading {
		content.WriteString(loadingBlockStyle.Render(" BUSY ") + " Testing connection...")
	} else if m.portModel.error != "" {
		content.WriteString(errorBlockStyle.Render(" FAIL ") + " " + m.portModel.error)
	} else if m.portModel.result != "" {
		content.WriteString(successBlockStyle.Render(" OPEN ") + " Connection successful:\n\n")
		content.WriteString(m.portModel.viewport.View())
	} else {
		content.WriteString(infoStyle.Render("Enter host and port, then press Enter"))
	}

	return content.String()
}

func (m MainModel) renderIPTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Render("🌐 IP INFORMATION"))
	content.WriteString("\n\n")

	// IP Type Selection
	content.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("IP Range: "),
		m.getIPTypeDisplayName(),
	))
	content.WriteString("\n\n")

	if m.ipModel.loading {
		content.WriteString(loadingBlockStyle.Render(" BUSY ") + " Fetching IP data...")
	} else if m.ipModel.error != "" {
		content.WriteString(errorBlockStyle.Render(" ERROR ") + " " + m.ipModel.error)
	} else if m.ipModel.result != "" {
		content.WriteString(successBlockStyle.Render(" OK ") + " IP Details:\n\n")
		content.WriteString(m.ipModel.viewport.View())
	} else {
		content.WriteString(infoStyle.Render("Use Up/Down to change range, Enter to fetch"))
	}

	return content.String()
}

func (m MainModel) renderProcessesTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Render("📊 NETWORK PROCESSES"))
	content.WriteString("\n\n")

	if m.processesModel.loading {
		content.WriteString(loadingBlockStyle.Render(" SCANNING ") + " Gathering connection list...")
	} else if m.processesModel.error != "" {
		content.WriteString(errorBlockStyle.Render(" ERROR ") + " " + m.processesModel.error)
	} else if m.processesModel.result != "" {
		content.WriteString(successBlockStyle.Render(" FOUND ") + " Recent Connections:\n\n")
		content.WriteString(m.processesModel.viewport.View())
	} else {
		content.WriteString(infoStyle.Render("Press Enter to scan active connections"))
	}

	return content.String()
}

func (m *MainModel) clearActiveTabContent() {
	switch m.activeTab {
	case 0: // DNS Lookup
		m.dnsModel.result = ""
		m.dnsModel.error = ""
		m.dnsModel.loading = false
		m.dnsModel.domainInput.SetValue("")
		m.dnsModel.viewport.SetContent("")
	case 1: // Port Test
		m.portModel.result = ""
		m.portModel.error = ""
		m.portModel.loading = false
		m.portModel.hostInput.SetValue("")
		m.portModel.portInput.SetValue("")
		m.portModel.viewport.SetContent("")
	case 2: // IP Address
		m.ipModel.result = ""
		m.ipModel.error = ""
		m.ipModel.loading = false
		m.ipModel.viewport.SetContent("")
	case 3: // Processes
		m.processesModel.result = ""
		m.processesModel.error = ""
		m.processesModel.loading = false
		m.processesModel.viewport.SetContent("")
	}
}

func (m *MainModel) initializeTabState() {
	m.dnsModel.result = ""
	m.dnsModel.error = ""
	m.dnsModel.loading = false
	m.dnsModel.domainInput.SetValue("example.com")

	m.portModel.result = ""
	m.portModel.error = ""
	m.portModel.loading = false
	m.portModel.hostInput.SetValue("localhost")
	m.portModel.portInput.SetValue("80")
	m.portModel.focused = 0

	m.ipModel.result = ""
	m.ipModel.error = ""
	m.ipModel.loading = false
	m.ipModel.selector = 2
	m.ipModel.ipType = "both"

	m.processesModel.result = ""
	m.processesModel.error = ""
	m.processesModel.loading = false

	m.setTabFocus()
}

func (m *MainModel) setTabFocus() {
	m.dnsModel.domainInput.Blur()
	m.dnsModel.domainInput.PromptStyle = unfocusedPromptStyle
	m.dnsModel.domainInput.TextStyle = unfocusedTextStyle

	m.portModel.hostInput.Blur()
	m.portModel.hostInput.PromptStyle = unfocusedPromptStyle
	m.portModel.hostInput.TextStyle = unfocusedTextStyle
	m.portModel.portInput.Blur()
	m.portModel.portInput.PromptStyle = unfocusedPromptStyle
	m.portModel.portInput.TextStyle = unfocusedTextStyle

	switch m.activeTab {
	case 0: // DNS Lookup
		m.dnsModel.domainInput.Focus()
		m.dnsModel.domainInput.PromptStyle = focusedPromptStyle
		m.dnsModel.domainInput.TextStyle = focusedTextStyle
	case 1: // Port Test
		m.portModel.hostInput.Focus()
		m.portModel.hostInput.PromptStyle = focusedPromptStyle
		m.portModel.hostInput.TextStyle = focusedTextStyle
		m.portModel.focused = 0
		// IP Address and Processes tabs don't have text inputs to focus
	}
}

func (m *MainModel) clearTabContent() {
	m.dnsModel.result = ""
	m.dnsModel.error = ""
	m.dnsModel.loading = false

	m.portModel.result = ""
	m.portModel.error = ""
	m.portModel.loading = false

	m.ipModel.result = ""
	m.ipModel.error = ""
	m.ipModel.loading = false
	m.ipModel.selector = 2
	m.ipModel.ipType = "both"

	m.processesModel.result = ""
	m.processesModel.error = ""
	m.processesModel.loading = false
}

func (m MainModel) getIPTypeDisplayName() string {
	typeName := ""
	switch m.ipModel.selector {
	case 0:
		typeName = "Private"
	case 1:
		typeName = "Public"
	case 2:
		typeName = "Both"
	default:
		typeName = "Both"
	}

	return activeTabStyle.
		Background(secondaryColor).
		Render(typeName)
}
