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
	// Primary colors - modern dark theme
	primaryColor   = lipgloss.Color("39")  // Bright blue accent
	secondaryColor = lipgloss.Color("69")  // Soft cyan
	accentColor    = lipgloss.Color("207") // Modern Purple/Pink accent
	errorColor     = lipgloss.Color("196") // Bright red
	warningColor   = lipgloss.Color("208") // Modern orange
	successColor   = lipgloss.Color("40")  // Success green
	mutedColor     = lipgloss.Color("242") // Light gray
	borderColor    = lipgloss.Color("63")  // Deep indigo border
	subtleColor    = lipgloss.Color("236") // Dark surface color
	highlightColor = lipgloss.Color("255") // Pure white for emphasis
	bgColor        = lipgloss.Color("233") // Deep background
)

// Unified styles for consistent appearance
var (
	// Main container that fills the window
	appStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(lipgloss.Color("252"))

	// Modern rounded tab styles (pill-shaped) - ensure background matches app
	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 2).
			MarginRight(1).
			Height(1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(subtleColor).
				Padding(0, 2).
				MarginRight(1).
				Height(1)

	// Modern content styles with subtle background layering
	tabOuterStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Background(lipgloss.Color("234"))

	// Modern header styles - ensure background matches container
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("234")).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Modern label styles - ensure background matches container
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Background(lipgloss.Color("234")).
			Bold(true).
			Width(10).
			MarginRight(1)

	// Text input background should also match container
	focusedPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Background(lipgloss.Color("234")).
				Bold(true)

	unfocusedPromptStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(lipgloss.Color("234"))

	focusedTextStyle = lipgloss.NewStyle().
				Foreground(highlightColor).
				Background(lipgloss.Color("234"))

	unfocusedTextStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(lipgloss.Color("234"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("234"))

	// Enhanced status styles with block backgrounds
	successBlockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(successColor).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	errorBlockStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(errorColor).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	loadingBlockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(warningColor).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	// Enhanced result styles with specialized background
	resultContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(subtleColor).
				Background(lipgloss.Color("235"))

	// Modern help styles with clean bar appearance
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(bgColor).
			Padding(0, 2).
			Height(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Background(bgColor).
			Bold(true)

	// Enhanced title style with gradient effect and modern framing
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(bgColor).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(primaryColor)

	// Styling for info text
	infoStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(lipgloss.Color("234"))
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
