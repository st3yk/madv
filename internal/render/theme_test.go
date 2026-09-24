package render

import "testing"

func TestResolveStyle(t *testing.T) {
	tests := []struct {
		name    string
		flag    string
		env     map[string]string
		tty     bool
		dark    bool
		want    string
		queried bool
	}{
		{name: "tmux falls back to neutral", env: map[string]string{"TMUX": "/tmp/tmux-1000/default,1,0", "TERM": "screen-256color"}, tty: true, want: StyleNeutral},
		{name: "screen TERM without TMUX", env: map[string]string{"TERM": "screen-256color"}, tty: true, want: StyleNeutral},
		{name: "GNU screen", env: map[string]string{"STY": "1234.pts-0", "TERM": "xterm"}, tty: true, want: StyleNeutral},
		{name: "plain terminal, dark background", env: map[string]string{"TERM": "xterm-256color"}, tty: true, dark: true, want: StyleDark, queried: true},
		{name: "plain terminal, light background", env: map[string]string{"TERM": "xterm-256color"}, tty: true, want: StyleLight, queried: true},
		{name: "not a tty", env: map[string]string{"TERM": "xterm-256color"}, want: StyleNotty},
		{name: "flag wins", flag: "dracula", env: map[string]string{"MADV_STYLE": "light", "GLAMOUR_STYLE": "dark", "NO_COLOR": "1"}, tty: true, want: "dracula"},
		{name: "MADV_STYLE beats GLAMOUR_STYLE", env: map[string]string{"MADV_STYLE": "light", "GLAMOUR_STYLE": "dark"}, tty: true, want: StyleLight},
		{name: "GLAMOUR_STYLE", env: map[string]string{"GLAMOUR_STYLE": "dark", "TMUX": "x"}, tty: true, want: StyleDark},
		{name: "auto is treated as unset", flag: "auto", env: map[string]string{"GLAMOUR_STYLE": "auto", "TMUX": "x"}, tty: true, want: StyleNeutral},
		{name: "NO_COLOR", env: map[string]string{"NO_COLOR": "1", "TMUX": "x"}, tty: true, want: StyleNotty},
		{name: "COLORFGBG light", env: map[string]string{"COLORFGBG": "0;15", "TMUX": "x"}, tty: true, want: StyleLight},
		{name: "COLORFGBG dark", env: map[string]string{"COLORFGBG": "15;default;0", "TMUX": "x"}, tty: true, want: StyleDark},
		{name: "COLORFGBG garbage ignored", env: map[string]string{"COLORFGBG": "15;default", "TMUX": "x"}, tty: true, want: StyleNeutral},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queried := false
			env := Env{
				Getenv: func(k string) string { return tt.env[k] },
				IsTTY:  func() bool { return tt.tty },
				HasDarkBackground: func() bool {
					queried = true
					return tt.dark
				},
			}
			if got := ResolveStyle(tt.flag, env); got != tt.want {
				t.Errorf("ResolveStyle() = %q, want %q", got, tt.want)
			}
			if queried != tt.queried {
				t.Errorf("background queried = %v, want %v", queried, tt.queried)
			}
		})
	}
}
