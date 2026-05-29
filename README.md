# Lista

A minimal, good-looking todo list for your terminal — CLI and TUI, pick your mode.

<p align="center">
  <a href="https://github.com/kwame-Owusu/lista/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kwame-Owusu/lista?style=flat" alt="License"></a>
  <a href="https://github.com/kwame-Owusu/lista/releases"><img src="https://img.shields.io/github/v/release/kwame-Owusu/lista?style=flat" alt="Release"></a>
</p>

---

## What is Lista?

Lista is a CLI-based todo list with a TUI mode, built for people who live in the terminal. Use it in **CLI mode** for quick add/list/complete actions, or drop into **TUI mode** for an interactive [Bubble Tea](https://github.com/charmbracelet/bubbletea) experience.

No sync, no cloud, no calendar. Just a JSON file, a terminal, and your todos.

## Why Lista?

I built this for myself. When I'm deep in a Neovim + Tmux, I don't want to reach for a GUI app or a browser tab to jot down a thought. I wanted something that:

- Fits in a tmux popup without breaking focus
- Works as fast as I can type
- Looks good enough that I don't mind looking at it

And it was a fun way to practice Go and build something with [Bubble Tea](https://github.com/charmbracelet/bubbletea). If it's useful to you too — great.

## Features

- **Two modes in one** — CLI for quick ops, TUI for interactive management
- **Priorities** — Low, Medium, High with color-coded badges
- **Notes** — Attach longer notes to any todo
- **Timestamps** — See when each todo was added ("added 5s ago", "added 10mins ago", etc.)
- **Themeable** — Ships with Gruvbox, customize any color in the config
- **Portable JSON** — Your data is just a file, sync it however you like

## Installation

### Homebrew

```bash
brew tap kwame-owusu/taps https://github.com/kwame-Owusu/homebrew-taps
brew install lista
```

### From source

```bash
git clone https://github.com/kwame-Owusu/lista.git
cd lista
go build -o lista
./lista --help
```

## Usage

### TUI mode

```bash
lista
```

Navigate with `↑/↓`, toggle with `space`, add with `a`, delete with `d`, edit with `e`, quit with `q`.

### CLI mode

```bash
# Add a todo with priority and notes
lista add "Buy groceries" --priority high --notes "get organic, check expiry"

# Add a simple todo
lista add "Walk the dog"

# List todos
lista list

# Complete a todo
lista complete 1

# View a todo with notes
lista view 1
```

| Command    | Description            |
|------------|-----------------------|
| `add`      | Add a new todo        |
| `list`     | List all todos        |
| `complete` | Mark a todo done      |
| `delete`   | Remove a todo         |
| `edit`     | Change the title      |
| `view`     | Show full details     |
| `notes`    | Add notes to a todo   |
| `export`   | Export todos as JSON  |

## Tmux integration

Add this to your `~/.tmux.conf` to pop Lista in a floating window:

```bash
bind-key l display-popup -w 80% -h 80% -E "lista"
```

With `CTRL+a` as your prefix, hitting `CTRL+a` then `l` opens Lista in a centered floating pane.


https://github.com/user-attachments/assets/705a15b5-110a-4e48-8944-27e65270bf2f

## Configuration

Lista ships with Gruvbox (because I like Gruvbox and you should too 🫡). To customize, edit `~/.config/lista/lista.config.json`:

```json
{
  "theme": {
    "background": "#282828",
    "foreground": "#ebdbb2",
    "accent": "#d79921",
    ...
  }
}
```

Just swap the hex values to make your own theme.

## Development

```bash
git clone https://github.com/kwame-Owusu/lista.git
cd lista
go build -o lista
```

Run tests:

```bash
make test
```
