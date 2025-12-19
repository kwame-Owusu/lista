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


# Configuration
There is a little that can be configured about lista, and that is the UI theme, when installed it comes with a default theme which is using
gruvbox colors(because I like gruvbox and you should too 🫡), but if you do not like it, its possible to change that inside the config file at:
 ```
 ./config/lista/lista.config.json
 ```
The config is nothing complicated, just replace the colors in the theme object to create a new theme using hex color codes.


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
My bind-key is Ctrl+a, so in this case to invoke lista in our tmux session we use hit CTRL+a followed by l: 
```bash
bind-key l display-popup -w 80% -h 80% -E "lista"
```


https://github.com/user-attachments/assets/ba3af0fe-f913-4ca6-886e-eb86ec1b0329


