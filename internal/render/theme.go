package render

import (
	"os"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// Style names understood by madv in addition to Glamour's standard styles.
const (
	StyleDark    = "dark"
	StyleLight   = "light"
	StyleNotty   = "notty"
	StyleNeutral = "neutral"
)

// Env holds the environment lookups and terminal probes used by ResolveStyle.
// It exists so the resolution order can be tested without a real terminal.
type Env struct {
	Getenv func(string) string
	// IsTTY reports whether stdout is a terminal.
	IsTTY func() bool
	// HasDarkBackground queries the terminal for its background color.
	// It must only be called before the TUI takes over the terminal.
	HasDarkBackground func() bool
}

// OSEnv returns an Env backed by the real process environment and stdout.
func OSEnv() Env {
	return Env{
		Getenv: os.Getenv,
		IsTTY:  func() bool { return term.IsTerminal(int(os.Stdout.Fd())) },
		HasDarkBackground: func() bool {
			return termenv.NewOutput(os.Stdout).HasDarkBackground()
		},
	}
}

// ResolveStyle picks the style to render with, in priority order:
//
//  1. the --style flag
//  2. MADV_STYLE
//  3. GLAMOUR_STYLE
//  4. NO_COLOR -> notty
//  5. COLORFGBG
//  6. terminal background query (skipped under tmux/screen, where it cannot work)
//  7. neutral, a style readable on both light and dark backgrounds
//
// It must be called before the TUI starts: the background query reads from the
// terminal and would race with the TUI's input reader.
func ResolveStyle(flag string, env Env) string {
	for _, s := range []string{flag, env.Getenv("MADV_STYLE"), env.Getenv("GLAMOUR_STYLE")} {
		if s != "" && s != "auto" {
			return s
		}
	}

	if env.Getenv("NO_COLOR") != "" {
		return StyleNotty
	}

	if s, ok := styleFromColorFGBG(env.Getenv("COLORFGBG")); ok {
		return s
	}

	if !env.IsTTY() {
		return StyleNotty
	}

	if !inMultiplexer(env) {
		if env.HasDarkBackground() {
			return StyleDark
		}
		return StyleLight
	}

	return StyleNeutral
}

// styleFromColorFGBG parses COLORFGBG ("fg;bg" or "fg;default;bg").
func styleFromColorFGBG(v string) (string, bool) {
	if !strings.Contains(v, ";") {
		return "", false
	}
	parts := strings.Split(v, ";")
	bg, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil || bg < 0 || bg > 15 {
		return "", false
	}
	if bg == 7 || bg >= 9 {
		return StyleLight, true
	}
	return StyleDark, true
}

// inMultiplexer reports whether we run inside tmux or screen. termenv refuses
// to query the background there and silently reports black instead.
func inMultiplexer(env Env) bool {
	if env.Getenv("TMUX") != "" || env.Getenv("STY") != "" {
		return true
	}
	t := env.Getenv("TERM")
	return strings.HasPrefix(t, "screen") || strings.HasPrefix(t, "tmux")
}

// IsLightStyle reports whether a style is designed for light backgrounds.
func IsLightStyle(style string) bool {
	return style == StyleLight
}

// IsDarkStyle reports whether a style is designed for dark backgrounds.
func IsDarkStyle(style string) bool {
	switch style {
	case StyleDark, "dracula", "tokyo-night", "pink":
		return true
	}
	return false
}
