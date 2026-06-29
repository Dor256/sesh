package commands

import "fmt"

func (commander *Commander) Destroy(args []string) error {
	worktree, err := commander.picker()
	if err != nil {
		fmt.Println("Aborting")
		return err
	}
	fmt.Println("Killing worktree and session...")
	sessionId := commander.gitClient.ReadSessionId(worktree.Path)

	commander.tmuxClient.Kill(sessionId)
	commander.gitClient.DestroyWorktree(worktree.Path)
	return nil
}
