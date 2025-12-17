# 🤔 What is Lista?

Lista is a simple CLI-based todo list application with support for notes.  
As the name suggests, it’s intentionally minimal,  but also pleasant to use — thanks to
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

It’s nothing fancy — just a small tool that helps me stay organized. If it’s useful to others who enjoy the terminal like I do then even better.


 # 🔨 Installation
You can install Lista using Homebrew. For now, it’s available via my personal tap:
 ```bash
  brew tap kwame-Owusu/lista
  brew install lista
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
  ./
 ```

 # Tmux + Nvim + Lista demo
 add link to the bash script which we use in TMUX

# License

MIT License

Copyright (c) 2025 kwame owusu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
