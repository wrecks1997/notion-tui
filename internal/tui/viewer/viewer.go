package viewer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6C5CE7")).
			Padding(0, 1).
			Width(80)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

type Model struct {
	Viewport viewport.Model
	title    string
	ready    bool
	width    int
	height   int
}

func New() Model {
	return Model{}
}

func (m *Model) SetContent(title, markdown string, width int) {
	if width < 20 {
		width = 80
	}
	renderWidth := width - 6
	if renderWidth < 20 {
		renderWidth = 20
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(renderWidth),
	)

	var rendered string
	if err != nil {
		rendered = markdown
	} else {
		rendered, err = renderer.Render(markdown)
		if err != nil {
			rendered = markdown
		}
	}

	m.title = title
	m.Viewport.SetContent(rendered)
	m.Viewport.GotoTop()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3
		footerHeight := 2

		if !m.ready {
			m.Viewport = viewport.New(msg.Width-4, msg.Height-headerHeight-footerHeight)
			m.Viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.Viewport.Width = msg.Width - 4
			m.Viewport.Height = msg.Height - headerHeight - footerHeight
		}
	}

	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	titleWidth := m.width - 4
	if titleWidth < 10 {
		titleWidth = 80
	}
	header := titleBarStyle.Width(titleWidth).Render(m.title)

	pct := fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100)
	info := infoStyle.Render(fmt.Sprintf("  esc: back · j/k: scroll · %s", pct))

	return strings.Join([]string{header, "", m.Viewport.View(), info}, "\n")
}
