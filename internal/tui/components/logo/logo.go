package logo

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/catpaladin/net-tools/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

const diag = `╱`

type Opts struct {
	FieldColor  color.Color
	TitleColorA color.Color
	TitleColorB color.Color
	Version     string
	Width       int
}

func Render(o Opts) string {
	title := "NET-TOOLS"

	// Apply gradient to title
	gradTitle := styles.ApplyBoldForegroundGrad(title, lipgloss.Color(toHex(o.TitleColorA)), lipgloss.Color(toHex(o.TitleColorB)))

	// Add version if provided
	if o.Version != "" {
		version := lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(o.FieldColor))).Render(" v" + o.Version)
		gradTitle += version
	}

	// Decorative fields (diagonal lines)
	fieldWidth := 10
	if o.Width > 0 {
		fieldWidth = (o.Width - lipgloss.Width(gradTitle)) / 2
		if fieldWidth < 2 {
			fieldWidth = 2
		}
	}

	field := lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(o.FieldColor))).Render(strings.Repeat(diag, fieldWidth))

	return lipgloss.JoinHorizontal(lipgloss.Center, field, " ", gradTitle, " ", field)
}

func toHex(c color.Color) string {
	if c == nil {
		return "#FFFFFF"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
