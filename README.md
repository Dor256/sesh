# Sesh

Sesh is a coding session management tool.

## Prerequisites
- `tmux`
- `Neovim`
- `OpenCode`
- `fzf`

The motivation for this tool was to streamline my local development workflow with one command for a new task.
When starting a new task on linear:

- run `sesh <linear ticket number>`
- Sesh will create a detached branch from the current directory and attach the linear number to it. (it is possible to provide a branch name as a second argument to sesh)
- Sesh will then spin up a `tmux` session with 3 panes: `OpenCode`, `Neovim`, and `Terminal` all in the new worktree's working directory

