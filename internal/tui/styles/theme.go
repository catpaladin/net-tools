package styles

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

type Theme struct {
	Name   string
	IsDark bool

	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Tertiary  lipgloss.Color
	Accent    lipgloss.Color

	BgBase        lipgloss.Color
	BgBaseLighter lipgloss.Color
	BgSubtle      lipgloss.Color
	BgOverlay     lipgloss.Color

	FgBase      lipgloss.Color
	FgMuted     lipgloss.Color
	FgHalfMuted lipgloss.Color
	FgSubtle    lipgloss.Color
	FgSelected  lipgloss.Color

	Border      lipgloss.Color
	BorderFocus lipgloss.Color

	Success lipgloss.Color
	Error   lipgloss.Color
	Warning lipgloss.Color
	Info    lipgloss.Color

	White lipgloss.Color

	styles     *Styles
	stylesOnce sync.Once
}

type Styles struct {
	Base         lipgloss.Style
	SelectedBase lipgloss.Style

	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Text         lipgloss.Style
	TextSelected lipgloss.Style
	Muted        lipgloss.Style
	Subtle       lipgloss.Style

	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style

	// Help
	Help help.Styles
}

func (t *Theme) S() *Styles {
	t.stylesOnce.Do(func() {
		t.styles = t.buildStyles()
	})
	return t.styles
}

func (t *Theme) buildStyles() *Styles {
	base := lipgloss.NewStyle().
		Foreground(t.FgBase)

	return &Styles{
		Base: base,

		SelectedBase: base.Background(t.Primary),

		Title: base.
			Foreground(t.Accent).
			Bold(true),

		Subtitle: base.
			Foreground(t.Secondary).
			Bold(true),

		Text:         base,
		TextSelected: base.Background(t.Primary).Foreground(t.FgSelected),

		Muted: base.Foreground(t.FgMuted),

		Subtle: base.Foreground(t.FgSubtle),

		Success: base.Foreground(t.Success),

		Error: base.Foreground(t.Error),

		Warning: base.Foreground(t.Warning),

		Info: base.Foreground(t.Info),

		Help: help.Styles{
			ShortKey:       base.Foreground(t.FgMuted),
			ShortDesc:      base.Foreground(t.FgSubtle),
			ShortSeparator: base.Foreground(t.Border),
			Ellipsis:       base.Foreground(t.Border),
			FullKey:        base.Foreground(t.FgMuted),
			FullDesc:       base.Foreground(t.FgSubtle),
			FullSeparator:  base.Foreground(t.Border),
		},
	}
}

func NewCharmtoneTheme() *Theme {
	return &Theme{
		Name:   "charmtone",
		IsDark: true,

		Primary:   Charple,
		Secondary: Dolly,
		Tertiary:  Bok,
		Accent:    Zest,

		BgBase:        Pepper,
		BgBaseLighter: BBQ,
		BgSubtle:      Charcoal,
		BgOverlay:     Iron,

		FgBase:      Ash,
		FgMuted:     Squid,
		FgHalfMuted: Smoke,
		FgSubtle:    Oyster,
		FgSelected:  Salt,

		Border:      Charcoal,
		BorderFocus: Charple,

		Success: Guac,
		Error:   Sriracha,
		Warning: Zest,
		Info:    Malibu,

		White: Butter,
	}
}

type Manager struct {
	themes  map[string]*Theme
	current *Theme
}

var (
	defaultManager     *Manager
	defaultManagerOnce sync.Once
)

func DefaultManager() *Manager {
	defaultManagerOnce.Do(func() {
		defaultManager = &Manager{
			themes: make(map[string]*Theme),
		}
		t := NewCharmtoneTheme()
		defaultManager.Register(t)
		defaultManager.current = t
	})
	return defaultManager
}

func CurrentTheme() *Theme {
	return DefaultManager().Current()
}

func (m *Manager) Register(theme *Theme) {
	m.themes[theme.Name] = theme
}

func (m *Manager) Current() *Theme {
	return m.current
}

// Gradient Helpers

func ApplyForegroundGrad(input string, color1, color2 lipgloss.Color) string {
	if input == "" {
		return ""
	}
	var o strings.Builder
	clusters := ForegroundGrad(input, false, color1, color2)
	for _, c := range clusters {
		fmt.Fprint(&o, c)
	}
	return o.String()
}

func ApplyBoldForegroundGrad(input string, color1, color2 lipgloss.Color) string {
	if input == "" {
		return ""
	}
	var o strings.Builder
	clusters := ForegroundGrad(input, true, color1, color2)
	for _, c := range clusters {
		fmt.Fprint(&o, c)
	}
	return o.String()
}

func ForegroundGrad(input string, bold bool, color1, color2 lipgloss.Color) []string {
	if input == "" {
		return []string{""}
	}

	c1, _ := colorful.Hex(string(color1))
	c2, _ := colorful.Hex(string(color2))

	var clusters []string
	gr := uniseg.NewGraphemes(input)
	for gr.Next() {
		clusters = append(clusters, string(gr.Runes()))
	}

	if len(clusters) <= 1 {
		style := lipgloss.NewStyle().Foreground(color1)
		if bold {
			style = style.Bold(true)
		}
		return []string{style.Render(input)}
	}

	for i := 0; i < len(clusters); i++ {
		t := float64(i) / float64(len(clusters)-1)
		c := c1.BlendHcl(c2, t)
		hex := c.Hex()
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		if bold {
			style = style.Bold(true)
		}
		clusters[i] = style.Render(clusters[i])
	}
	return clusters
}
