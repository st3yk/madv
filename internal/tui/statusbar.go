package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236"))

	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	percentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("240")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)
)

// RenderStatusBar produces a full-width footer status bar.
func RenderStatusBar(name string, vp viewport.Model, width int) string {
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

	left := nameStyle.Render(name)
	middle := percentStyle.Render(percentStr)
	right := helpStyle.Render("q: quit • j/k: scroll • g/G: top/bot")

	leftAndMid := lipgloss.JoinHorizontal(lipgloss.Top, left, middle)
	usedWidth := lipgloss.Width(leftAndMid) + lipgloss.Width(right)

	if width < usedWidth {
		// On narrower terminals, show compact version
		compact := fmt.Sprintf(" %s │ %s ", name, percentStr)
		return statusBarStyle.Width(width).Render(compact)
	}

	fillWidth := width - usedWidth
	fill := statusBarStyle.Render(strings.Repeat(" ", fillWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftAndMid, fill, right)
}
