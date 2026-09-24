package render

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

// Render renders raw markdown into ANSI-styled terminal text using the given
// width and style. The style is either "neutral", one of Glamour's standard
// style names (e.g. "dark", "light", "notty") or a path to a JSON style file.
// Resolve the style once with ResolveStyle; Render never probes the terminal.
func Render(markdown string, width int, style string) (string, error) {
	if strings.TrimSpace(markdown) == "" {
		return "", nil
	}

	if width <= 0 {
		width = 80
	}

	var styleOpt glamour.TermRendererOption
	switch {
	case style == StyleNeutral || style == "" || style == styles.AutoStyle:
		styleOpt = glamour.WithStyles(NeutralStyleConfig())
	case styles.DefaultStyles[style] != nil:
		styleOpt = glamour.WithStandardStyle(style)
	default:
		styleOpt = glamour.WithStylePath(style)
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
