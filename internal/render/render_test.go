package render

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func TestRenderEmpty(t *testing.T) {
	out, err := Render("", 80)
	if err != nil {
		t.Fatalf("unexpected error rendering empty markdown: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}

	outWhitespace, err := Render("   \n\n  \t ", 80)
	if err != nil {
		t.Fatalf("unexpected error rendering whitespace: %v", err)
	}
	if outWhitespace != "" {
		t.Errorf("expected empty string for whitespace, got %q", outWhitespace)
	}
}

func TestRenderHeader(t *testing.T) {
	md := "# Hello World\nThis is a paragraph."
	out, err := RenderWithStyle(md, 80, "dark")
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
	out, err := RenderWithStyle(longText, width, "notty")
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
	out, err := RenderWithStyle(codeMd, 80, "dark")
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

func TestRenderNoColor(t *testing.T) {
	orig := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", orig)

	os.Setenv("NO_COLOR", "1")
	out, err := Render("# Plain Title", 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// In notty style with NO_COLOR, there should be no ANSI escape sequences
	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected NO_COLOR to prevent ANSI color sequences, got: %q", out)
	}
	if !strings.Contains(out, "Plain Title") {
		t.Errorf("expected output to contain 'Plain Title', got: %q", out)
	}
}
