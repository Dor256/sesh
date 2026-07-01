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

type Client struct {}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateWorktree(branchName, path string) error {
	return exec.Command("git", "worktree", "add", "-b", branchName, path, "origin/main").Run()
}

func (c *Client) DestroyWorktree(path string) error {
	return exec.Command("git", "-C", path, "worktree", "remove", "--force", ".").Run()
}

func (c *Client) ensureWorktreeConfig(path string) error {
	return exec.Command("git", "-C", path, "config", "extensions.worktreeConfig", "true").Run()
}

func (c *Client) SaveSessionId(path, tmuxSessionId string) error {
	if err := c.ensureWorktreeConfig(path); err != nil {
		return err
	}
	return exec.Command("git", "-C", path, "config", "--worktree", "custom.worktree.tmuxid", tmuxSessionId).Run()
}

func (c *Client) ReadSessionId(path string) string {
	return execAndPrint("git", "-C", path, "config", "--worktree", "custom.worktree.tmuxid")
}

