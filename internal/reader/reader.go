package reader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const maxFileSize = 10 * 1024 * 1024 // 10 MB limit for safety

// Input represents markdown content and its source identifier.
type Input struct {
	Name    string
	Content string
}

// IsBinary detects whether the provided byte slice represents binary data.
func IsBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	// Check first 512 bytes (standard sniff buffer size)
	sniffLen := 512
	if len(data) < sniffLen {
		sniffLen = len(data)
	}
	chunk := data[:sniffLen]

	// Check for null bytes (standard heuristic for binary files)
	if bytes.IndexByte(chunk, 0) != -1 {
		return true
	}

	contentType := http.DetectContentType(chunk)
	// If it doesn't start with text/ and isn't application/octet-stream fallback on empty text
	if !strings.HasPrefix(contentType, "text/") &&
		contentType != "application/octet-stream" &&
		contentType != "application/json" &&
		contentType != "application/xml" {
		return true
	}

	return false
}

// ReadInput parses arguments and stdin to produce an Input struct.
func ReadInput(args []string, stdin io.Reader, isStdinTerminal bool) (*Input, error) {
	if len(args) > 1 {
		return nil, errors.New("too many arguments; usage: madv [file | -]")
	}

	// Case 1: No arguments provided
	if len(args) == 0 {
		if isStdinTerminal {
			return nil, errors.New("no file specified; usage: madv [file | -]")
		}
		// Stdin is piped
		return readFromReader(stdin, "[stdin]")
	}

	target := args[0]

	// Case 2: Explicit stdin "-"
	if target == "-" {
		return readFromReader(stdin, "[stdin]")
	}

	// Case 3: File path
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %q does not exist", target)
		}
		return nil, fmt.Errorf("cannot access %q: %w", target, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%q is a directory, not a file", target)
	}

	if info.Size() > maxFileSize {
		return nil, fmt.Errorf("file %q exceeds maximum supported size (10MB)", target)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", target, err)
	}

	if IsBinary(data) {
		return nil, fmt.Errorf("%q appears to be a binary file", target)
	}

	return &Input{
		Name:    target,
		Content: string(data),
	}, nil
}

func readFromReader(r io.Reader, name string) (*Input, error) {
	lr := io.LimitReader(r, maxFileSize+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s: %w", name, err)
	}

	if int64(len(data)) > maxFileSize {
		return nil, fmt.Errorf("input exceeds maximum supported size (10MB)")
	}

	if IsBinary(data) {
		return nil, fmt.Errorf("input from %s appears to be binary data", name)
	}

	return &Input{
		Name:    name,
		Content: string(data),
	}, nil
}
