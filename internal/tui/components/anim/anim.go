// Package anim provides an animated spinner.
package anim

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	fps           = 20
	initialChar   = '.'
	labelGap      = " "
	labelGapWidth = 1

	// Periods of ellipsis animation speed in steps.
	ellipsisAnimSpeed = 8

	// The maximum amount of time that can pass before a character appears.
	maxBirthOffset = time.Second

	// Default number of cycling chars.
	defaultNumCyclingChars = 10
)

// Default colors for gradient.
var (
	defaultGradColorA = color.RGBA{R: 0x7D, G: 0x52, B: 0xFB, A: 0xff} // Charple
	defaultGradColorB = color.RGBA{R: 0x56, G: 0xCC, B: 0xF2, A: 0xff} // Malibu
	defaultLabelColor = color.RGBA{R: 0xBF, G: 0xBF, B: 0xBF, A: 0xff} // Oyster
)

var (
	availableRunes = []rune("0123456789abcdefABCDEF~!@#$£€%^&*()+=_")
	ellipsisFrames = []string{".", "..", "...", ""}
)

var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

// StepMsg is a message type used to trigger the next step in the animation.
type StepMsg struct{ id int }

// Settings defines settings for the animation.
type Settings struct {
	Size        int
	Label       string
	LabelColor  color.Color
	GradColorA  color.Color
	GradColorB  color.Color
	CycleColors bool
}

// Anim is a Bubble for an animated spinner.
type Anim struct {
	mu               sync.RWMutex
	width            int
	cyclingCharWidth int
	label            []string
	labelWidth       int
	labelColor       color.Color
	startTime        time.Time
	birthOffsets     []time.Duration
	initialFrames    [][]string
	initialized      atomic.Bool
	cyclingFrames    [][]string
	step             atomic.Int64
	ellipsisStep     atomic.Int64
	ellipsisFrames   []string
	id               int
}

// New creates a new Anim instance.
func New(opts Settings) *Anim {
	if opts.Size < 1 {
		opts.Size = defaultNumCyclingChars
	}
	if colorIsUnset(opts.GradColorA) {
		opts.GradColorA = defaultGradColorA
	}
	if colorIsUnset(opts.GradColorB) {
		opts.GradColorB = defaultGradColorB
	}
	if colorIsUnset(opts.LabelColor) {
		opts.LabelColor = defaultLabelColor
	}

	a := &Anim{}
	a.id = nextID()
	a.startTime = time.Now()
	a.cyclingCharWidth = opts.Size
	a.labelColor = opts.LabelColor
	a.labelWidth = lipgloss.Width(opts.Label)

	// Total width of anim
	a.width = opts.Size
	if opts.Label != "" {
		a.width += labelGapWidth + a.labelWidth
	}

	a.renderLabel(opts.Label)

	// Pre-generate gradient ramp
	numFrames := 20
	var ramp []color.Color
	if opts.CycleColors {
		ramp = makeGradientRamp(a.width*3, opts.GradColorA, opts.GradColorB, opts.GradColorA, opts.GradColorB)
		numFrames = a.width * 2
	} else {
		ramp = makeGradientRamp(a.width, opts.GradColorA, opts.GradColorB)
	}

	// Pre-render initial frames
	a.initialFrames = make([][]string, numFrames)
	offset := 0
	for i := range a.initialFrames {
		a.initialFrames[i] = make([]string, a.width)
		for j := range a.initialFrames[i] {
			if j+offset >= len(ramp) {
				continue
			}
			var c color.Color
			if j < a.cyclingCharWidth {
				c = ramp[j+offset]
			} else {
				c = opts.LabelColor
			}
			a.initialFrames[i][j] = lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(c))).Render(string(initialChar))
		}
		if opts.CycleColors {
			offset++
		}
	}

	// Pre-render cycling frames
	a.cyclingFrames = make([][]string, numFrames)
	offset = 0
	for i := range a.cyclingFrames {
		a.cyclingFrames[i] = make([]string, a.cyclingCharWidth)
		for j := range a.cyclingFrames[i] {
			if j+offset >= len(ramp) {
				continue
			}
			r := availableRunes[rand.IntN(len(availableRunes))]
			a.cyclingFrames[i][j] = lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(ramp[j+offset]))).Render(string(r))
		}
		if opts.CycleColors {
			offset++
		}
	}

	// Birth offsets
	a.birthOffsets = make([]time.Duration, a.width)
	for i := range a.birthOffsets {
		a.birthOffsets[i] = time.Duration(rand.N(int64(maxBirthOffset))) * time.Nanosecond
	}

	return a
}

func (a *Anim) renderLabel(label string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.labelWidth > 0 {
		labelRunes := []rune(label)
		a.label = make([]string, len(labelRunes))
		for i := range labelRunes {
			a.label[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(a.labelColor))).Render(string(labelRunes[i]))
		}

		a.ellipsisFrames = make([]string, len(ellipsisFrames))
		for i, frame := range ellipsisFrames {
			a.ellipsisFrames[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(toHex(a.labelColor))).Render(frame)
		}
	} else {
		a.label = nil
		a.ellipsisFrames = nil
	}
}

func (a *Anim) Init() tea.Cmd {
	return a.Step()
}

func (a *Anim) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StepMsg:
		if msg.id != a.id {
			return a, nil
		}

		step := a.step.Add(1)
		if int(step) >= len(a.cyclingFrames) {
			a.step.Store(0)
		}

		if a.initialized.Load() && a.labelWidth > 0 {
			ellipsisStep := a.ellipsisStep.Add(1)
			if int(ellipsisStep) >= ellipsisAnimSpeed*len(ellipsisFrames) {
				a.ellipsisStep.Store(0)
			}
		} else if !a.initialized.Load() && time.Since(a.startTime) >= maxBirthOffset {
			a.initialized.Store(true)
		}
		return a, a.Step()
	}
	return a, nil
}

func (a *Anim) View() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var b strings.Builder
	step := int(a.step.Load())
	for i := range a.width {
		switch {
		case !a.initialized.Load() && i < len(a.birthOffsets) && time.Since(a.startTime) < a.birthOffsets[i]:
			b.WriteString(a.initialFrames[step][i])
		case i < a.cyclingCharWidth:
			b.WriteString(a.cyclingFrames[step][i])
		case i == a.cyclingCharWidth:
			b.WriteString(labelGap)
		case i > a.cyclingCharWidth:
			idx := i - a.cyclingCharWidth - labelGapWidth
			if idx < len(a.label) {
				b.WriteString(a.label[idx])
			}
		}
	}

	if a.initialized.Load() && a.labelWidth > 0 {
		ellipsisStep := int(a.ellipsisStep.Load())
		frameIdx := ellipsisStep / ellipsisAnimSpeed
		if frameIdx < len(a.ellipsisFrames) {
			b.WriteString(a.ellipsisFrames[frameIdx])
		}
	}

	return b.String()
}

// SetLabel updates the label for the animation.
func (a *Anim) SetLabel(label string) {
	a.renderLabel(label)
	a.mu.Lock()
	a.labelWidth = lipgloss.Width(label)
	// Recalculate width
	a.width = a.cyclingCharWidth
	if label != "" {
		a.width += labelGapWidth + a.labelWidth
	}
	a.mu.Unlock()
}

func (a *Anim) Step() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return StepMsg{id: a.id}
	})
}

func makeGradientRamp(size int, stops ...color.Color) []color.Color {
	if len(stops) < 2 {
		return nil
	}
	points := make([]colorful.Color, len(stops))
	for i, k := range stops {
		points[i], _ = colorful.MakeColor(k)
	}
	numSegments := len(stops) - 1
	blended := make([]color.Color, 0, size)
	segmentSize := size / numSegments
	for i := 0; i < numSegments; i++ {
		c1 := points[i]
		c2 := points[i+1]
		for j := 0; j < segmentSize; j++ {
			t := float64(j) / float64(segmentSize)
			blended = append(blended, c1.BlendHcl(c2, t))
		}
	}
	return blended
}

func toHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

func colorIsUnset(c color.Color) bool {
	if c == nil {
		return true
	}
	_, _, _, a := c.RGBA()
	return a == 0
}
