# 🤔 What is Lista?

Lista is a simple CLI-based todo list application with support for notes.  
As the name suggests, it’s intentionally minimal, but also pleasant to use, thanks to
[Bubble Tea](https://github.com/charmbracelet/bubbletea).

Lista can be used in two ways:

- **Normal mode** — where you interact with it using commands and flags
- **TUI mode** — an interactive terminal UI that lets you manage your todos in real time

Both modes are designed to be simple and fast to work with.

# 🤔 Why Lista?

I built Lista primarily for myself.

When I’m programming, I often need a lightweight way to jot down small notes — what to work on next, reminders, or ideas that pop up mid-flow. While there are plenty of existing tools, I’ve been going deep into the **Neovim + TMUX + Bash** workflow, and I wanted something that fit naturally into that environment. But also a way for me to practice my Go skills.

Lista is meant to be:

- Easy to invoke inside a TMUX session
- Fast enough to not break focus
- Simple enough to stay out of the way

It’s nothing fancy, just a small tool that helps me stay organized. If it’s useful to others who enjoy the terminal like I do then even better.

# 🔨 Installation

You can install Lista using Homebrew. For now, it’s available via my personal tap:

```bash
brew tap kwame-owusu/tap https://github.com/kwame-owusu/homebrew-tap
brew install lista
lista --version
```

# 🛠️ Development

To build the project locally, you must:

```bash
 git clone https://github.com/kwame-Owusu/lista.git
 cd lista
```

Build the binary and run it:

```bash
  go build -o lista
  ./lista
```

# Tmux + Nvim + Lista demo
For this to work, you need to have tmux installed and have a few different bindings, but all boil down to preference.
My bind-key is Ctrl+a: 
```bash
bind-key l display-popup -w 80% -h 80% -E "lista"
```
