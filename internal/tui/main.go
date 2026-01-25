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
		m.contentWidth = msg.Width - 4

		inputWidth := m.contentWidth - 8
		m.dnsModel.viewport.Width = m.contentWidth - 10
		m.portModel.viewport.Width = m.contentWidth - 10
		m.ipModel.viewport.Width = m.contentWidth - 10
		m.processesModel.viewport.Width = m.contentWidth - 10

		m.dnsModel.domainInput.Width = inputWidth
		m.portModel.hostInput.Width = inputWidth
		m.portModel.portInput.Width = 10

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

	title := titleStyle.Render("NET-TOOLS")

	tabs := make([]string, len(m.tabs))
	for i, tab := range m.tabs {
		if tab.Active {
			tabs[i] = activeTabStyle.Render(tab.Name)
		} else {
			tabs[i] = inactiveTabStyle.Render(tab.Name)
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

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

	help := m.renderHelp()

	ui := lipgloss.JoinVertical(lipgloss.Left, title, tabBar, content, help)

	return appStyle.Render(ui)
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
		helpKey := helpKeyStyle.Background(bgColor).Render(k.key)
		helpValue := helpValueStyle.Background(bgColor).Render(k.desc)
		part := fmt.Sprintf("%s %s", helpKey, helpValue)
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

	content.WriteString(headerStyle.Width(m.contentWidth - 2).Render("🔍 DNS LOOKUP"))

	inputRow := inputRowStyle.Width(m.contentWidth - 6).Render(lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render("Domain: "),
		m.dnsModel.domainInput.View(),
	))
	content.WriteString(inputRow)

	if m.dnsModel.loading {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(loadingBlockStyle.Render(" ⟳ ") + " Looking up records..."))
	} else if m.dnsModel.error != "" {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(errorBlockStyle.Render(" ❌ ") + " " + m.dnsModel.error))
	} else if m.dnsModel.result != "" {
		header := successBlockStyle.Render(" ✓ ") + " DNS Records Found:\n\n"
		resultContent := header + m.dnsModel.viewport.View()
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(resultContent))
	}

	return content.String()
}

func (m MainModel) renderPortTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Width(m.contentWidth - 2).Render("🔌 PORT CONNECTIVITY"))

	content.WriteString(inputRowStyle.Width(m.contentWidth - 6).Render(lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render("Host: "),
		m.portModel.hostInput.View(),
	)))

	content.WriteString(inputRowStyle.Width(m.contentWidth - 6).Render(lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render("Port: "),
		m.portModel.portInput.View(),
	)))

	if m.portModel.loading {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(loadingBlockStyle.Render(" ⟳ ") + " Testing connection..."))
	} else if m.portModel.error != "" {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(errorBlockStyle.Render(" ❌ ") + " " + m.portModel.error))
	} else if m.portModel.result != "" {
		header := successBlockStyle.Render(" ✓ ") + " Connection successful:\n\n"
		resultContent := header + m.portModel.viewport.View()
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(resultContent))
	}

	return content.String()
}

func (m MainModel) renderIPTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Width(m.contentWidth - 2).Render("🌐 IP INFORMATION"))

	content.WriteString(inputRowStyle.Width(m.contentWidth - 6).Render(lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render("IP Range: "),
		m.getIPTypeDisplayName(),
	)))

	if m.ipModel.loading {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(loadingBlockStyle.Render(" ⟳ ") + " Fetching IP data..."))
	} else if m.ipModel.error != "" {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(errorBlockStyle.Render(" ❌ ") + " " + m.ipModel.error))
	} else if m.ipModel.result != "" {
		header := successBlockStyle.Render(" ✓ ") + " IP Details:\n\n"
		resultContent := header + m.ipModel.viewport.View()
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(resultContent))
	}

	return content.String()
}

func (m MainModel) renderProcessesTab() string {
	var content strings.Builder

	content.WriteString(headerStyle.Width(m.contentWidth - 2).Render("📊 NETWORK PROCESSES"))

	if m.processesModel.loading {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(loadingBlockStyle.Render(" ⟳ ") + " Gathering connection list..."))
	} else if m.processesModel.error != "" {
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(errorBlockStyle.Render(" ❌ ") + " " + m.processesModel.error))
	} else if m.processesModel.result != "" {
		header := successBlockStyle.Render(" ✓ ") + " Recent Connections:\n\n"
		resultContent := header + m.processesModel.viewport.View()
		content.WriteString(infoBoxStyle.Width(m.contentWidth - 6).Render(resultContent))
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
