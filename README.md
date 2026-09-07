# madv

`madv` is a lightweight, terminal-based markdown viewer written in Go with Vim-style keybindings and responsive rendering (similar to `less`).

Built and managed with **Bazel** and **Bzlmod**.

## Features

- **Rich Markdown Formatting**: Headings, bold/italic, lists, blockquotes, syntax-highlighted code blocks, and formatted tables powered by [Glamour](https://github.com/charmbracelet/glamour).
- **Vim Navigation**: Familiar `less` / Vim keybindings for navigation.
- **Dynamic Terminal Wrapping**: Resizes gracefully on window changes without clipping or breaking text layouts.
- **Piped Stdin Support**: Pipe markdown output directly (`curl`, `cat`, or tool outputs) into `madv` and retain interactive navigation.
- **Safety**: Rejects binary files to protect the terminal and handles panics gracefully by restoring terminal state.

## Installation & Uninstallation

### Install to `/usr/bin`

To build the binary using Bazel and install it directly to `/usr/bin/madv`:

```bash
./install.sh
```

*(Note: Automatically prompts for `sudo` if `/usr/bin` is not directly writable by your user).*

### Uninstall from `/usr/bin`

To remove the binary from `/usr/bin`:

```bash
./uninstall.sh
```

## Building with Bazel

### Build the Binary

```bash
bazel build //:madv
```

The compiled binary will be located at `bazel-bin/cmd/madv/madv_/madv`.

### Run Directly via Bazel

```bash
# View a markdown file
bazel run //:madv -- README.md

# View from stdin pipe
cat testdata/sample.md | bazel run //:madv -- -

# View help
bazel run //:madv -- --help
```

### Run Tests

```bash
bazel test //...
```

### Regenerate BUILD Files with Gazelle

```bash
bazel run //:gazelle
```

## Keybindings

| Key | Description |
| :--- | :--- |
| `j`, `↓`, `Enter` | Scroll down 1 line |
| `k`, `↑` | Scroll up 1 line |
| `d`, `Ctrl+D` | Scroll down half page |
| `u`, `Ctrl+U` | Scroll up half page |
| `f`, `Space`, `PgDown` | Scroll down full page |
| `b`, `PgUp` | Scroll up full page |
| `g`, `Home` | Jump to beginning of document |
| `G`, `End` | Jump to end of document |
| `q`, `Esc`, `Ctrl+C` | Quit viewer cleanly |
