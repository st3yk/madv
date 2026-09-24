package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"madv/internal/render"
)

// Model represents the Bubble Tea TUI state.
type Model struct {
	Name       string
	RawContent string
	Style      string
	Viewport   viewport.Model
	KeyMap     KeyMap
	Width      int
	Height     int
	Ready      bool
	Quitting   bool
}

// NewModel initializes a new TUI model. style must already be resolved (see
// render.ResolveStyle); the model never probes the terminal itself.
func NewModel(name, rawContent, style string) Model {
	return Model{
		Name:       name,
		RawContent: rawContent,
		Style:      style,
		KeyMap:     DefaultKeyMap(),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.KeyMap.Top):
			m.Viewport.GotoTop()
			return m, nil

		case key.Matches(msg, m.KeyMap.Bottom):
			m.Viewport.GotoBottom()
			return m, nil

		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.Viewport.HalfViewUp()
			return m, nil

		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.Viewport.HalfViewDown()
			return m, nil

		case key.Matches(msg, m.KeyMap.PageUp):
			m.Viewport.ViewUp()
			return m, nil

		case key.Matches(msg, m.KeyMap.PageDown):
			m.Viewport.ViewDown()
			return m, nil

		case key.Matches(msg, m.KeyMap.Up):
			m.Viewport.LineUp(1)
			return m, nil

		case key.Matches(msg, m.KeyMap.Down):
			m.Viewport.LineDown(1)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		footerHeight := 1
		vpHeight := msg.Height - footerHeight
		if vpHeight < 1 {
			vpHeight = 1
		}

		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, vpHeight)
			// Disable default viewport key bindings to avoid conflicts with our custom KeyMap
			m.Viewport.KeyMap = viewport.KeyMap{}
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = vpHeight
		}

		// Re-render content to fit the current window width
		rendered := m.renderContent(msg.Width)
		m.Viewport.SetContent(rendered)
		return m, nil
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m Model) renderContent(width int) string {
	if strings.TrimSpace(m.RawContent) == "" {
		return "\n  *(empty document)*\n"
	}

	rendered, err := render.Render(m.RawContent, width, m.Style)
	if err != nil {
		return fmt.Sprintf("\n  Error rendering markdown: %v\n", err)
	}

	return rendered
}

// View implements tea.Model.
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	if !m.Ready {
		return "Initializing madv..."
	}

	if m.Height < 3 || m.Width < 10 {
		return "Terminal too small"
	}

	vpView := m.Viewport.View()
	statusBar := RenderStatusBar(m.Name, m.Viewport, m.Width, m.Style)

	return fmt.Sprintf("%s\n%s", vpView, statusBar)
}
