package render

import (
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

// Render renders raw markdown into ANSI-styled terminal text using the specified width.
func Render(markdown string, width int) (string, error) {
	if strings.TrimSpace(markdown) == "" {
		return "", nil
	}

	if width <= 0 {
		width = 80
	}

	style := "auto"
	if os.Getenv("NO_COLOR") != "" {
		style = "notty"
	}

	return RenderWithStyle(markdown, width, style)
}

// RenderWithStyle renders markdown using a specified Glamour style name or path (e.g. "dark", "light", "notty", "auto").
func RenderWithStyle(markdown string, width int, style string) (string, error) {
	if strings.TrimSpace(markdown) == "" {
		return "", nil
	}

	if width <= 0 {
		width = 80
	}

	var styleOpt glamour.TermRendererOption
	switch style {
	case "notty":
		styleOpt = glamour.WithStandardStyle("notty")
	case "dark":
		styleOpt = glamour.WithStandardStyle("dark")
	case "light":
		styleOpt = glamour.WithStandardStyle("light")
	case "auto":
		fallthrough
	default:
		styleOpt = glamour.WithAutoStyle()
	}

	renderer, err := glamour.NewTermRenderer(
		styleOpt,
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		return "", err
	}

	// Glamour often adds trailing newlines; trim excessive trailing whitespace while preserving layout
	return strings.TrimRight(out, "\n"), nil
}
