package render

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func TestRenderEmpty(t *testing.T) {
	out, err := Render("", 80, StyleNeutral)
	if err != nil {
		t.Fatalf("unexpected error rendering empty markdown: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}

	outWhitespace, err := Render("   \n\n  \t ", 80, StyleNeutral)
	if err != nil {
		t.Fatalf("unexpected error rendering whitespace: %v", err)
	}
	if outWhitespace != "" {
		t.Errorf("expected empty string for whitespace, got %q", outWhitespace)
	}
}

func TestRenderHeader(t *testing.T) {
	md := "# Hello World\nThis is a paragraph."
	out, err := Render(md, 80, "dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify ANSI styling sequences are present in styled output
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected styled output to contain ANSI escape sequences")
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "Hello World") {
		t.Errorf("expected stripped output to contain 'Hello World', got: %s", plain)
	}
	if !strings.Contains(plain, "This is a paragraph.") {
		t.Errorf("expected stripped output to contain paragraph, got: %s", plain)
	}
}

func TestRenderWordWrap(t *testing.T) {
	longText := "This is a very long line of text that is intended to test word wrapping inside the madv markdown renderer component."
	width := 40
	out, err := Render(longText, width, "notty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(out, "\n")
	if len(lines) <= 1 {
		t.Errorf("expected multiple lines due to wrapping at width %d, got %d lines", width, len(lines))
	}
}

func TestRenderCodeBlock(t *testing.T) {
	codeMd := "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"hi\")\n}\n```"
	out, err := Render(codeMd, 80, "dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify syntax highlighting injected ANSI codes
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected code block to contain syntax highlighting ANSI escapes")
	}

	plain := stripANSI(out)
	if !strings.Contains(plain, "package main") {
		t.Errorf("expected stripped code block to contain 'package main', got: %s", plain)
	}
	if !strings.Contains(plain, "fmt.Println") {
		t.Errorf("expected stripped code block to contain 'fmt.Println', got: %s", plain)
	}
}

func TestRenderNotty(t *testing.T) {
	out, err := Render("# Plain Title", 80, StyleNotty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected notty style to emit no ANSI sequences, got: %q", out)
	}
	if !strings.Contains(out, "Plain Title") {
		t.Errorf("expected output to contain 'Plain Title', got: %q", out)
	}
}

// Regression test: under tmux madv used the dark style, whose near-white body
// text and dark backgrounds made light terminals unreadable. The neutral style
// must set no background colors and leave body text in the terminal default.
func TestRenderNeutralHasNoBackgrounds(t *testing.T) {
	// Covers every element that uses a background or near-white text in the
	// dark style: H1, body text, inline code, code blocks, tables, quotes.
	md := "# Welcome to `madv`\n\nBody text with `inline code` and a [link](https://example.com).\n\n" +
		"## Code\n\n```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"hi\")\n}\n```\n\n" +
		"| Key | Action |\n| --- | --- |\n| `j` | Down |\n\n> Quote by Dijkstra\n\n---\n"
	out, err := Render(md, 80, StyleNeutral)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sawColor := false
	for _, seq := range ansiRegex.FindAllString(out, -1) {
		fg, bg := sgrColors(seq)
		sawColor = sawColor || len(fg) > 0
		if len(bg) > 0 {
			t.Errorf("neutral style emitted a background color %v in %q", bg, seq)
		}
		for _, c := range fg {
			if c == "38;5;252" {
				t.Errorf("neutral style emitted near-white text in %q", seq)
			}
		}
	}

	if !sawColor {
		t.Fatalf("expected neutral style to emit foreground colors; output had none")
	}

	plain := stripANSI(out)
	for _, want := range []string{"Welcome to", "fmt.Println", "inline code", "Dijkstra"} {
		if !strings.Contains(plain, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestRenderUnknownStyle(t *testing.T) {
	if _, err := Render("# x", 80, "no-such-style"); err == nil {
		t.Errorf("expected an error for an unknown style")
	}
}

// sgrColors returns the foreground and background color parameters of an SGR
// escape sequence, e.g. "\x1b[1;38;5;32;48;5;236m" -> ["38;5;32"], ["48;5;236"].
func sgrColors(seq string) (fg, bg []string) {
	if !strings.HasSuffix(seq, "m") {
		return nil, nil
	}
	params := strings.Split(strings.TrimSuffix(strings.TrimPrefix(seq, "\x1b["), "m"), ";")
	for i := 0; i < len(params); i++ {
		p := params[i]
		n, _ := strconv.Atoi(p)
		switch {
		case p == "38" || p == "48":
			end := i + 1
			if end < len(params) && params[end] == "5" {
				end += 2
			} else if end < len(params) && params[end] == "2" {
				end += 4
			}
			if end > len(params) {
				end = len(params)
			}
			c := strings.Join(params[i:end], ";")
			if p == "38" {
				fg = append(fg, c)
			} else {
				bg = append(bg, c)
			}
			i = end - 1
		case (n >= 30 && n <= 37) || (n >= 90 && n <= 97):
			fg = append(fg, p)
		case (n >= 40 && n <= 47) || (n >= 100 && n <= 107):
			bg = append(bg, p)
		}
	}
	return fg, bg
}
