package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab in the TUI
type Tab struct {
	Name    string
	Content string
	Active  bool
	Loading bool
	Error   string
}

// MainModel represents the main TUI model with tabs
type MainModel struct {
	tabs           []Tab
	activeTab      int
	width          int
	height         int
	contentWidth   int
	ready          bool
	dnsModel       *DNSModel
	portModel      *PortModel
	ipModel        *IPModel
	processesModel *ProcessesModel
	spinner        spinner.Model
}

// DNSModel represents the DNS lookup tab
type DNSModel struct {
	domainInput textinput.Model
	result      string
	loading     bool
	error       string
	viewport    viewport.Model
}

// PortModel represents the port testing tab
type PortModel struct {
	hostInput textinput.Model
	portInput textinput.Model
	result    string
	loading   bool
	error     string
	focused   int // 0 for host, 1 for port
	viewport  viewport.Model
}

// IPModel represents the IP address tab
type IPModel struct {
	ipType   string
	result   string
	loading  bool
	error    string
	selector int // 0=private, 1=public, 2=both
	viewport viewport.Model
}

// ProcessesConnection represents a single network connection
type ProcessesConnection struct {
	LocalAddr  string
	Port       string
	PIDProgram string
}

// ProcessesModel represents the processes tab
type ProcessesModel struct {
	result      string
	loading     bool
	error       string
	connections []ProcessesConnection
	viewport    viewport.Model
}

// Key bindings
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Tab   key.Binding
	Clear key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "execute"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Clear: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Modern color palette for consistent theming
var (
	// Primary colors - modern dark theme with enhanced contrast
	primaryColor   = lipgloss.Color("85")  // Bright cyan
	secondaryColor = lipgloss.Color("86")  // Aqua?
	accentColor    = lipgloss.Color("226") // Warm yellow/orange
	errorColor     = lipgloss.Color("196") // Bright red
	warningColor   = lipgloss.Color("220") // Orange/coral
	successColor   = lipgloss.Color("46")  // Success green
	mutedColor     = lipgloss.Color("240") // Light gray
	borderColor    = lipgloss.Color("240") // Subtle gray border
	subtleColor    = lipgloss.Color("238") // Dark surface color
	highlightColor = lipgloss.Color("255") // Pure white for emphasis
	bgColor        = lipgloss.NoColor{}    // Transparent/Default
	overlayColor   = lipgloss.Color("231") // Light overlay
)

// Unified styles for consistent appearance
var (
	// Main container that fills the window
	appStyle = lipgloss.NewStyle().
			Margin(1, 2)

	// Modern tab styles with enhanced visual separation
	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 1).
				MarginRight(1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor)

	// Modern header styles - ensure background matches container
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder(), false, false, true, false).
			BorderForeground(primaryColor)

	// Modern label styles - ensure background matches container
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Background(bgColor).
			Bold(true).
			Width(10).
			MarginRight(1)

	// Input row background style
	inputRowStyle = lipgloss.NewStyle().
			Padding(1).
			MarginTop(1).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor)

	// Info box style
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1).
			MarginTop(1)

	// Info style for helper text
	infoStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(bgColor)

	// Help text styles
	helpKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(bgColor)

	helpValueStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Background(bgColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(bgColor).
			Padding(0, 1).
			MarginTop(1)

	// Status block styles
	loadingBlockStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Background(bgColor)

	errorBlockStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Background(bgColor)

	successBlockStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Background(bgColor)

	// Cursor style
	cursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(bgColor)

	// Text input styles for focus states
	focusedPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(primaryColor).
				Bold(true)

	focusedTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(bgColor)

	unfocusedPromptStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Background(bgColor)

	unfocusedTextStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(bgColor)
)

// Messages for async operations
type dnsResultMsg struct {
	result string
	err    error
}

type portResultMsg struct {
	result string
	err    error
}

type ipResultMsg struct {
	result string
	err    error
}

type processesResultMsg struct {
	result string
	err    error
}
