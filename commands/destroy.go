package commands

import (
	"fmt"
	"os"
	"sesh/git"
	"sesh/tmux"
)

func Destroy(args []string) {
	worktree, err := Picker()
	if err != nil {
		fmt.Println("Aborting")
		os.Exit(1)
	}
	fmt.Println("Killing worktree and session...")
	sessionId := git.ReadSessionId(worktree.Path)

	tmux.Kill(sessionId)
	git.DestroyWorktree(worktree.Path)
}
