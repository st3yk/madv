package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"madv/internal/render"
)

// statusBarStyles holds the styles for each status bar segment.
type statusBarStyles struct {
	bar, name, percent, help lipgloss.Style
}

func newStatusBarStyles(barFg, barBg, nameFg, nameBg, percentBg, helpFg string) statusBarStyles {
	barStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(barFg)).
		Background(lipgloss.Color(barBg))
	return statusBarStyles{
		bar: barStyle,
		name: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(nameFg)).
			Background(lipgloss.Color(nameBg)).
			Padding(0, 1),
		percent: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(nameFg)).
			Background(lipgloss.Color(percentBg)).
			Padding(0, 1),
		help: barStyle.
			Foreground(lipgloss.Color(helpFg)).
			Padding(0, 1),
	}
}

var (
	darkStatusBar  = newStatusBarStyles("252", "236", "230", "62", "240", "245")
	lightStatusBar = newStatusBarStyles("236", "254", "230", "62", "244", "240")

	// Reverse video follows the terminal's own colors, so it is readable on
	// any background when we do not know which one we are on.
	neutralStatusBar = statusBarStyles{
		bar:     lipgloss.NewStyle().Reverse(true),
		name:    lipgloss.NewStyle().Reverse(true).Bold(true).Padding(0, 1),
		percent: lipgloss.NewStyle().Reverse(true).Bold(true).Padding(0, 1),
		help:    lipgloss.NewStyle().Reverse(true).Padding(0, 1),
	}
)

func stylesFor(style string) statusBarStyles {
	switch {
	case render.IsDarkStyle(style):
		return darkStatusBar
	case render.IsLightStyle(style):
		return lightStatusBar
	}
	return neutralStatusBar
}

// RenderStatusBar produces a full-width footer status bar.
func RenderStatusBar(name string, vp viewport.Model, width int, style string) string {
	if width <= 0 {
		return ""
	}

	// Calculate scroll percentage string
	var percentStr string
	switch {
	case vp.TotalLineCount() <= vp.Height:
		percentStr = "All"
	case vp.AtTop() || vp.YOffset == 0:
		percentStr = "Top"
	case vp.PastBottom() || vp.AtBottom():
		percentStr = "Bot"
	default:
		percentStr = fmt.Sprintf("%2.f%%", vp.ScrollPercent()*100)
	}

	st := stylesFor(style)

	left := st.name.Render(name)
	middle := st.percent.Render(percentStr)
	right := st.help.Render("q: quit • j/k: scroll • g/G: top/bot")

	leftAndMid := lipgloss.JoinHorizontal(lipgloss.Top, left, middle)
	usedWidth := lipgloss.Width(leftAndMid) + lipgloss.Width(right)

	if width < usedWidth {
		// On narrower terminals, show compact version
		compact := fmt.Sprintf(" %s │ %s ", name, percentStr)
		return st.bar.Width(width).Render(compact)
	}

	fillWidth := width - usedWidth
	fill := st.bar.Render(strings.Repeat(" ", fillWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftAndMid, fill, right)
}
