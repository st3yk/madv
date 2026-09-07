package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func createTestModel(lines int) Model {
	var sb strings.Builder
	for i := 1; i <= lines; i++ {
		sb.WriteString("Line item number in document\n\n")
	}
	m := NewModel("test.md", sb.String())
	// Send initial WindowSizeMsg (width: 80, height: 20)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	return updated.(Model)
}

func TestModelInitialization(t *testing.T) {
	m := NewModel("sample.md", "# Title")
	if m.Ready {
		t.Errorf("model should not be ready before WindowSizeMsg")
	}
	if m.Init() != nil {
		t.Errorf("Init() should return nil")
	}

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(Model)
	if !m.Ready {
		t.Errorf("model should be ready after WindowSizeMsg")
	}
	if m.Viewport.Width != 80 {
		t.Errorf("expected viewport width 80, got %d", m.Viewport.Width)
	}
	if m.Viewport.Height != 23 { // 24 - 1 footer
		t.Errorf("expected viewport height 23, got %d", m.Viewport.Height)
	}
}

func TestVimNavigation(t *testing.T) {
	m := createTestModel(100)

	// Test 1: j (down 1 line)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.Viewport.YOffset != 1 {
		t.Errorf("expected YOffset 1 after 'j', got %d", m.Viewport.YOffset)
	}

	// Test 2: k (up 1 line)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.Viewport.YOffset != 0 {
		t.Errorf("expected YOffset 0 after 'k', got %d", m.Viewport.YOffset)
	}

	// Test 3: Boundary clamping at top
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.Viewport.YOffset != 0 {
		t.Errorf("expected YOffset 0 when clamping top, got %d", m.Viewport.YOffset)
	}

	// Test 4: d (half page down)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(Model)
	halfHeight := m.Viewport.Height / 2
	if m.Viewport.YOffset != halfHeight {
		t.Errorf("expected YOffset %d after 'd', got %d", halfHeight, m.Viewport.YOffset)
	}

	// Test 5: u (half page up)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	m = updated.(Model)
	if m.Viewport.YOffset != 0 {
		t.Errorf("expected YOffset 0 after 'u', got %d", m.Viewport.YOffset)
	}

	// Test 6: G (jump to bottom)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = updated.(Model)
	if !m.Viewport.AtBottom() {
		t.Errorf("expected viewport to be at bottom after 'G'")
	}

	// Test 7: g (jump to top)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = updated.(Model)
	if m.Viewport.YOffset != 0 {
		t.Errorf("expected YOffset 0 after 'g', got %d", m.Viewport.YOffset)
	}

	// Test 8: q (quit)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Errorf("expected tea.Quit command after 'q'")
	}
}

func TestResizeReflow(t *testing.T) {
	m := createTestModel(50)

	// Resize to smaller width
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 30})
	m = updated.(Model)

	if m.Viewport.Width != 40 {
		t.Errorf("expected viewport width 40, got %d", m.Viewport.Width)
	}
	if m.Viewport.Height != 29 { // 30 - 1 footer
		t.Errorf("expected viewport height 29, got %d", m.Viewport.Height)
	}
}

func TestSmallTerminalDisplay(t *testing.T) {
	m := NewModel("test.md", "content")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 8, Height: 2})
	m = updated.(Model)

	view := m.View()
	if !strings.Contains(view, "Terminal too small") {
		t.Errorf("expected 'Terminal too small' for tiny dimensions, got: %s", view)
	}
}

func TestEmptyFileView(t *testing.T) {
	m := NewModel("empty.md", "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(Model)

	view := m.View()
	if !strings.Contains(view, "empty document") {
		t.Errorf("expected empty document placeholder in view, got: %s", view)
	}
}

func TestStatusBarRendering(t *testing.T) {
	// Short content (<= height)
	vpShort := viewport.New(80, 20)
	vpShort.SetContent("Line 1\nLine 2")
	barShort := RenderStatusBar("sample.md", vpShort, 80)
	if !strings.Contains(barShort, "sample.md") {
		t.Errorf("expected status bar to contain filename, got: %s", barShort)
	}
	if !strings.Contains(barShort, "All") {
		t.Errorf("expected status bar to contain 'All' for short document, got: %s", barShort)
	}

	// Long content (> height)
	vpLong := viewport.New(80, 5)
	vpLong.SetContent("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")
	barLong := RenderStatusBar("long.md", vpLong, 80)
	if !strings.Contains(barLong, "Top") {
		t.Errorf("expected status bar to contain 'Top' for long document at line 0, got: %s", barLong)
	}
}
