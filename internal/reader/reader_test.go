package reader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"empty", []byte(""), false},
		{"plain text", []byte("Hello, world!\nThis is markdown."), false},
		{"with null byte", []byte("Hello \x00 World"), true},
		{"elf header snippet", []byte("\x7fELF\x02\x01\x01\x00"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBinary(tt.data); got != tt.expected {
				t.Errorf("IsBinary() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReadInput(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Valid markdown file
	goodFile := filepath.Join(tempDir, "test.md")
	if err := os.WriteFile(goodFile, []byte("# Title\nContent"), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Binary file
	binFile := filepath.Join(tempDir, "blob.bin")
	if err := os.WriteFile(binFile, []byte("start\x00end"), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("valid file", func(t *testing.T) {
		in, err := ReadInput([]string{goodFile}, strings.NewReader(""), true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Name != goodFile {
			t.Errorf("expected name %s, got %s", goodFile, in.Name)
		}
		if in.Content != "# Title\nContent" {
			t.Errorf("unexpected content: %s", in.Content)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := ReadInput([]string{filepath.Join(tempDir, "missing.md")}, strings.NewReader(""), true)
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("expected does not exist error, got: %v", err)
		}
	})

	t.Run("directory target", func(t *testing.T) {
		_, err := ReadInput([]string{subDir}, strings.NewReader(""), true)
		if err == nil || !strings.Contains(err.Error(), "is a directory") {
			t.Errorf("expected directory error, got: %v", err)
		}
	})

	t.Run("binary file target", func(t *testing.T) {
		_, err := ReadInput([]string{binFile}, strings.NewReader(""), true)
		if err == nil || !strings.Contains(err.Error(), "appears to be a binary file") {
			t.Errorf("expected binary file error, got: %v", err)
		}
	})

	t.Run("no args interactive terminal", func(t *testing.T) {
		_, err := ReadInput([]string{}, strings.NewReader(""), true)
		if err == nil || !strings.Contains(err.Error(), "no file specified") {
			t.Errorf("expected no file specified error, got: %v", err)
		}
	})

	t.Run("no args piped stdin", func(t *testing.T) {
		stdinContent := "# Piped Input\nSome text"
		in, err := ReadInput([]string{}, strings.NewReader(stdinContent), false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Name != "[stdin]" {
			t.Errorf("expected name [stdin], got %s", in.Name)
		}
		if in.Content != stdinContent {
			t.Errorf("expected %q, got %q", stdinContent, in.Content)
		}
	})

	t.Run("explicit stdin dash", func(t *testing.T) {
		stdinContent := "Dashboard text"
		in, err := ReadInput([]string{"-"}, strings.NewReader(stdinContent), true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Name != "[stdin]" {
			t.Errorf("expected name [stdin], got %s", in.Name)
		}
		if in.Content != stdinContent {
			t.Errorf("expected %q, got %q", stdinContent, in.Content)
		}
	})

	t.Run("too many arguments", func(t *testing.T) {
		_, err := ReadInput([]string{"file1.md", "file2.md"}, strings.NewReader(""), true)
		if err == nil || !strings.Contains(err.Error(), "too many arguments") {
			t.Errorf("expected too many arguments error, got: %v", err)
		}
	})
}
