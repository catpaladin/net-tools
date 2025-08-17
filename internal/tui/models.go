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
	tabs         []Tab
	activeTab    int
	width        int
	height       int
	contentWidth int
	ready        bool
	digModel     *DigModel
	ncModel      *NetcatModel
	ipModel      *IPModel
	netstatModel *NetstatModel
	spinner      spinner.Model
}

// DigModel represents the DNS lookup tab
type DigModel struct {
	domainInput textinput.Model
	result      string
	loading     bool
	error       string
	viewport    viewport.Model
}

// NetcatModel represents the port testing tab
type NetcatModel struct {
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

// NetstatConnection represents a single network connection
type NetstatConnection struct {
	LocalAddr  string
	Port       string
	PIDProgram string
}

// NetstatModel represents the netstat tab
type NetstatModel struct {
	result      string
	loading     bool
	error       string
	connections []NetstatConnection
	viewport    viewport.Model
}

// Key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Tab    key.Binding
	Escape key.Binding
	Quit   key.Binding
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
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// Color palette for consistent theming
var (
	// Primary colors
	primaryColor   = lipgloss.Color("39")  // Bright blue
	secondaryColor = lipgloss.Color("205") // Bright magenta
	accentColor    = lipgloss.Color("46")  // Bright green
	errorColor     = lipgloss.Color("196") // Bright red
	warningColor   = lipgloss.Color("214") // Orange
	mutedColor     = lipgloss.Color("240") // Gray
	borderColor    = lipgloss.Color("62")  // Dark blue
)

// Unified styles for consistent appearance
var (
	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(0, 3).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			BorderTop(true).
			BorderRight(true).
			BorderLeft(true).
			BorderBottom(false)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 3).
				Margin(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				BorderTop(true).
				BorderRight(true).
				BorderLeft(true).
				BorderBottom(false)

	// Content styles
	tabContentStyle = lipgloss.NewStyle().
			Padding(2, 3).
			Margin(1, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(primaryColor)

	// Input container styles (for layout only)
	inputContainerStyle = lipgloss.NewStyle().
				MarginBottom(1)

	// Textinput prompt styles (applied directly to textinput)
	focusedPromptStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	unfocusedPromptStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Textinput text styles
	focusedTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")) // White text

	unfocusedTextStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Textinput cursor style
	cursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Label styles
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginRight(1)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)

	loadingStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)

	// Result styles
	resultHeaderStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	resultContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(1).
				Margin(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Background(lipgloss.Color("235"))

	// IP type display style
	ipTypeDisplayStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(2).
			Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor)
)

// Messages for async operations
type digResultMsg struct {
	result string
	err    error
}

type netcatResultMsg struct {
	result string
	err    error
}

type ipResultMsg struct {
	result string
	err    error
}

type netstatResultMsg struct {
	result string
	err    error
}
