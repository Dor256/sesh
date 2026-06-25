package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sesh/git"
	"sesh/tmux"
)

func createWorktree(worktree Worktree, worktreePath string) {
	if worktree.Branch == "" {
		fmt.Println("Aborting")
		os.Exit(1)
	}
	err := git.CreateWorktree(worktree.Branch, worktreePath)
	if err != nil {
		fmt.Printf("Git failed: %v\n", err)
		return
	}
}

func Start(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't extract CWD. Aborting")
		os.Exit(1)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't extract user's HOME dir")
	}

	if len(args) > 0 && args[0] != "-n" {
		fmt.Printf("Unknown argument %s\n", args[0])
		os.Exit(1)
	}
	worktreeDir := fmt.Sprintf("%s/dev/.worktrees", home)
	worktree, err := Picker()
	cwdBase := filepath.Base(cwd)
	var worktreePath string = worktree.Path

	// New worktree
	if err != nil {
		worktreePath = fmt.Sprintf("%s/%s-%s", worktreeDir, cwdBase, args[1])
		createWorktree(*worktree, worktreePath)
	}
	sessionName := fmt.Sprintf("ENG%s", args[1])
	var sessionId string
	if tmux.HasSession(sessionName) {
		sessionId = tmux.GetSessionId(sessionName)
	} else {
		sessionId = tmux.Create(sessionName, worktreePath)

		tmux.NewWindow(sessionName, "OpenCode", worktreePath, "opencode")
		tmux.NewWindow(sessionName, "Neovim", worktreePath, "nvim .")
		tmux.NewWindow(sessionName, "Terminal", worktreePath)
		tmux.NewWindow(sessionName, "Claude", worktreePath, "claude")
	}

	git.SaveSessionId(worktreePath, sessionId)
	// If inside tmux switch to session if outside attach to session
	if os.Getenv("TMUX") != "" {
		tmux.Switch(sessionName)
	} else {
		tmux.Attach(sessionName)
	}
}
