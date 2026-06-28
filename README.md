# Sesh
<div align="center">
    <img src="./sesh-logo.png" width="300">
</div>

Sesh is a coding session management tool.

## Prerequisites
- `tmux`
- `Neovim`
- `OpenCode`
- `fzf`
- `Claude Code`

The motivation for this tool was to streamline my local development workflow with one command for a new task.
When starting a new task:

- run `sesh create` (optional flag `-s session-name` to skip the prompt)
- Sesh will create a detached branch from the current directory and attach the session name to it
- Sesh will then open an `fzf` picker where you can select the worktree to work on. If selected - Sesh will either open the existing session on `tmux` or create a new session if none exists
- Sesh will then spin up a `tmux` session with 4 panes: `OpenCode`, `Neovim`, `Terminal`, and `Claude` all in the new worktree's working directory

