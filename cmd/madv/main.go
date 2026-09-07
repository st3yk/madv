package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"madv/internal/reader"
	"madv/internal/tui"
)

const version = "0.1.0"

func printUsage() {
	fmt.Println("madv - A simple terminal markdown viewer with Vim keybindings")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  madv <file>       View markdown file")
	fmt.Println("  madv -            Read and view markdown from stdin")
	fmt.Println("  cat doc.md | madv View piped markdown")
	fmt.Println()
	fmt.Println("Keybindings:")
	fmt.Println("  j, ↓, Enter       Scroll down 1 line")
	fmt.Println("  k, ↑              Scroll up 1 line")
	fmt.Println("  d, Ctrl+D         Scroll down half page")
	fmt.Println("  u, Ctrl+U         Scroll up half page")
	fmt.Println("  f, Space, PgDown  Scroll down full page")
	fmt.Println("  b, PgUp           Scroll up full page")
	fmt.Println("  g, Home           Jump to top")
	fmt.Println("  G, End            Jump to bottom")
	fmt.Println("  q, Esc, Ctrl+C    Quit viewer")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help        Show this help screen")
	fmt.Println("  -v, --version     Show madv version")
}

func main() {
	// Defensive panic recovery to ensure terminal is restored from raw/alt-screen mode
	defer func() {
		if r := recover(); r != nil {
			// Restore cursor and exit alternate screen
			fmt.Print("\033[?1049l\033[?25h")
			fmt.Fprintf(os.Stderr, "panic recovered in madv: %v\n", r)
			os.Exit(2)
		}
	}()

	args := os.Args[1:]

	// Flags check
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("madv v%s\n", version)
			os.Exit(0)
		}
	}

	isTerminal := term.IsTerminal(int(os.Stdin.Fd()))
	input, err := reader.ReadInput(args, os.Stdin, isTerminal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	programOpts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}

	// If stdin was piped, open /dev/tty for user keyboard interaction
	if !isTerminal || (len(args) == 1 && args[0] == "-") {
		if tty, err := os.Open("/dev/tty"); err == nil {
			defer tty.Close()
			programOpts = append(programOpts, tea.WithInput(tty))
		}
	}

	// Initialize Bubble Tea TUI
	model := tui.NewModel(input.Name, input.Content)
	p := tea.NewProgram(model, programOpts...)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running viewer: %v\n", err)
		os.Exit(1)
	}
}
