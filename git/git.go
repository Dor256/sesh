package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execAndPrint(cmd string, args ...string) string {
	stdout, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println("Error reading stdout", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(stdout))
}

func CreateWorktree(branchName, path string) error {
	return exec.Command("git", "worktree", "add", "-b", branchName, path, "origin/main").Run()
}

func DestroyWorktree(path string) error {
	return exec.Command("git", "-C", path, "worktree", "remove", "--force", ".").Run()
}

func SaveSessionId(path, tmuxSessionId string) error {
	return exec.Command("git", "-C", path, "config", "--local", "custom.worktree.tmuxid", tmuxSessionId).Run()
}

func ReadSessionId(path string) string {
	return execAndPrint("git", "-C", path, "config", "--local", "custom.worktree.tmuxid")
}

