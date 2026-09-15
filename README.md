<h1 align="center">Lista</h1>

<p align="center">A minimal, good-looking todo list for your terminal — CLI and TUI, pick your mode.</p>

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

- **Two modes in one** — use the CLI for quick ops or the TUI for interactive management
- **Priorities & notes** — Low / Medium / High badges and longer notes on any todo
- **Smart lists** — filter by status, priority, or search; toggle and uncomplete with one command
- **Piped input** — `echo "Buy milk" | lista add`
- **Built for speed in the TUI** — add/edit forms, undo (`u`), purge (`c`), and a `?` help overlay
- **Themeable & portable** — Gruvbox by default, any color via config; plain JSON data with auto-backup of corrupt files

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

| Keys                 | Action                    |
| -------------------- | ------------------------- |
| `↑` / `k`, `↓` / `j` | Move up / down            |
| `space`              | Toggle complete           |
| `a`                  | Add a todo                |
| `e`                  | Edit the selected todo    |
| `d` / `x`            | Delete (confirm with `y`) |
| `c`                  | Purge completed todos     |
| `u`                  | Undo last action          |
| `?`                  | Show help overlay         |
| `q` / `ctrl+c`       | Quit                      |

Press `?` inside the TUI for the full keybinding reference, including the add/edit forms (`tab` to switch fields, `←/→/h/l` to change priority, `enter`/`ctrl+s` to save, `esc` to cancel).

### CLI mode

```bash
# Add a todo with priority and notes
lista add "Add new middleware" --priority high --notes "simple middleware to track visits"

# Add a todo from piped input
echo "Buy milk" | lista add

# Add a simple todo
lista add "Update docs"

# List todos
lista list

# Filter by status, priority, or search
lista list --pending
lista list --done --priority high
lista list --search docs

# Complete a todo
lista complete 1

# Toggle or revert completion
lista toggle 1
lista uncomplete 1

# View a todo with notes
lista view 1
```

| Command      | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| `add`        | Add a new todo (accepts piped input)                                   |
| `list`       | List todos (+ `--pending`, `--done`, `--priority`, `--search` filters) |
| `complete`   | Mark a todo done                                                       |
| `uncomplete` | Mark a todo pending                                                    |
| `toggle`     | Flip a todo's completion status                                        |
| `delete`     | Remove a todo                                                          |
| `edit`       | Change the title                                                       |
| `view`       | Show full details                                                      |
| `notes`      | Add notes to a todo                                                    |
| `export`     | Export todos as Markdown                                               |

## Tmux integration

Add this to your `~/.tmux.conf` to pop Lista in a floating window:

```bash
bind-key l display-popup -w 80% -h 80% -E "lista"
```

With `CTRL+a` as your prefix, hitting `CTRL+a` then `l` opens Lista in a centered floating pane.

https://github.com/user-attachments/assets/307b21a6-8c09-4eea-9948-2f6168fa772c



## Configuration

Lista ships with Gruvbox (because I like Gruvbox and you should too 🫡). To customize, edit `~/.config/lista/lista.config.json`:

```json
{
  "theme": {
    "background": "#3c3836",
    "background_alt": "#282828",
    "text_primary": "#ebdbb2",
    "text_secondary": "#c5c7bc",
    "text_muted": "#a89984",
    "priority_high": "#fb4934",
    "priority_medium": "#fe8019",
    "priority_low": "#b8bb26",
    "accent": "#fabd2f",
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
