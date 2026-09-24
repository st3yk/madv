package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"madv/internal/reader"
	"madv/internal/render"
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
	fmt.Println("  -s, --style NAME  Color style: neutral, dark, light, notty, dracula,")
	fmt.Println("                    tokyo-night, pink, or a path to a Glamour JSON style")
	fmt.Println("  -h, --help        Show this help screen")
	fmt.Println("  -v, --version     Show madv version")
	fmt.Println()
	fmt.Println("Environment:")
	fmt.Println("  MADV_STYLE        Default style (overridden by --style)")
	fmt.Println("  GLAMOUR_STYLE     Used when MADV_STYLE is unset")
	fmt.Println("  NO_COLOR          Disable colors")
	fmt.Println()
	fmt.Println("Without a style, madv detects the terminal background. Under tmux/screen")
	fmt.Println("that is not possible, so the 'neutral' style is used; set MADV_STYLE to")
	fmt.Println("'dark' or 'light' for a full theme.")
}

// parseStyleFlag extracts "-s NAME", "--style NAME" or "--style=NAME" from args
// and returns the remaining arguments.
func parseStyleFlag(args []string) (style string, rest []string, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-s" || a == "--style":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("%s requires a style name", a)
			}
			style = args[i+1]
			i++
		case strings.HasPrefix(a, "--style="):
			style = strings.TrimPrefix(a, "--style=")
		default:
			rest = append(rest, a)
		}
	}
	return style, rest, nil
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

	styleFlag, args, err := parseStyleFlag(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

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

	// Resolve the color style before Bubble Tea takes over the terminal: the
	// background query reads the terminal's reply and would otherwise race
	// with Bubble Tea's input reader.
	style := render.ResolveStyle(styleFlag, render.OSEnv())
	if _, err := render.Render("x", 80, style); err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid style: %v\n", err)
		os.Exit(1)
	}

	// Initialize Bubble Tea TUI
	model := tui.NewModel(input.Name, input.Content, style)
	p := tea.NewProgram(model, programOpts...)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running viewer: %v\n", err)
		os.Exit(1)
	}
}
