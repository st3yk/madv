package render

import (
	"github.com/alecthomas/chroma/v2"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/glamour/ansi"
	glamourstyles "github.com/charmbracelet/glamour/styles"
)

// neutralChromaTheme is registered under its own name: Glamour registers its
// per-style chroma theme under a single global name ("charm") on first use,
// so reusing that mechanism would leak colors between styles.
const neutralChromaTheme = "madv-neutral"

func init() {
	// Mid-luminance colors only and no backgrounds, so code stays readable on
	// both light and dark terminals. Plain text keeps the terminal foreground.
	chromastyles.Register(chroma.MustNewStyle(neutralChromaTheme, chroma.StyleEntries{
		chroma.Error:               "bold #d70000",
		chroma.Comment:             "italic #808080",
		chroma.CommentPreproc:      "#d75f00",
		chroma.Keyword:             "#0087d7",
		chroma.KeywordReserved:     "#af00af",
		chroma.KeywordNamespace:    "#d7005f",
		chroma.KeywordType:         "#5f5fd7",
		chroma.Operator:            "#d75f5f",
		chroma.NameBuiltin:         "#af005f",
		chroma.NameTag:             "#8700d7",
		chroma.NameAttribute:       "#5f5faf",
		chroma.NameClass:           "bold underline",
		chroma.NameDecorator:       "#af8700",
		chroma.NameFunction:        "#00875f",
		chroma.LiteralNumber:       "#008787",
		chroma.LiteralString:       "#af5f00",
		chroma.LiteralStringEscape: "#00af87",
		chroma.GenericDeleted:      "#d70000",
		chroma.GenericInserted:     "#00af00",
		chroma.GenericEmph:         "italic",
		chroma.GenericStrong:       "bold",
		chroma.GenericSubheading:   "#808080",
	}))
}

// NeutralStyleConfig returns a Glamour style that sets no background colors
// and leaves body text in the terminal's own foreground color, so it is
// readable regardless of the terminal background.
func NeutralStyleConfig() ansi.StyleConfig {
	// Shallow copy: only reassign pointer fields, never write through them,
	// or the shared DarkStyleConfig global would be modified.
	cfg := glamourstyles.DarkStyleConfig

	cfg.Document.Color = nil

	cfg.Heading.Color = strPtr("32")
	cfg.H1.Color = nil
	cfg.H1.BackgroundColor = nil
	cfg.H1.Prefix = "# "
	cfg.H1.Suffix = ""

	cfg.Code.Color = strPtr("161")
	cfg.Code.BackgroundColor = nil

	cfg.CodeBlock.Chroma = nil
	cfg.CodeBlock.Theme = neutralChromaTheme

	cfg.Image.Color = strPtr("170")
	cfg.HorizontalRule.Color = strPtr("244")

	return cfg
}

func strPtr(s string) *string { return &s }
