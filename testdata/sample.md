# Welcome to `madv`

`madv` is a terminal markdown viewer written in Go with Vim keybindings.

## Features

- **Rich Markdown Formatting**: Headings, lists, blockquotes, code blocks.
- **Vim Navigation**:
  - `j` / `k`: Line scroll
  - `d` / `u`: Half page scroll
  - `f` / `b`: Full page scroll
  - `g` / `G`: Top / Bottom
  - `q`: Quit viewer

### Code Block Example

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello from madv!")
}
```

### Table Example

| Command | Keybinding | Action |
| :--- | :--- | :--- |
| Down | `j` or `↓` | Scroll down 1 line |
| Up | `k` or `↑` | Scroll up 1 line |
| Half Down | `d` | Half page down |
| Half Up | `u` | Half page up |
| Page Down | `f` or `Space` | Page down |
| Page Up | `b` | Page up |
| Top | `g` | Jump to beginning |
| Bottom | `G` | Jump to end |
| Quit | `q` | Exit cleanly |

> "Simplicity is prerequisite for reliability."
> — Edsger W. Dijkstra
