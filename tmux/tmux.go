package tmux

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

func (c *Client) Attach(sessionName string) error {
	return exec.Command("tmux", "attach", "-t", sessionName).Run()
}

func (c *Client) Switch(sessionName string) error {
	return exec.Command("tmux", "switch", "-t", sessionName).Run()
}

func (c *Client) Create(sessionName, path string) string {
	return execAndPrint("tmux", "new-session", "-d", "-P", "-F", "#{session_id}", "-s", sessionName, "-c", path)
}

func (c *Client) NewWindow(sessionName, windowName, path string, tool ...string) error {
	if len(tool) == 0 {
		// Plain terminal
		return exec.Command("tmux", "new-window", "-t", sessionName, "-c", path, "-n", windowName).Run()
	}
	toolCmd := fmt.Sprintf("%s; zsh", tool[0])
	return exec.Command("tmux", "new-window", "-t", sessionName, "-c", path, "-n", windowName, toolCmd).Run()
}

func (c *Client) RenameWindow(sessionName, oldName, newName string) error {
	return exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("%s:%s", sessionName, oldName), newName).Run()
}

func (c *Client) Kill(sessionNameOrId string) error {
	return exec.Command("tmux", "kill-session", "-t", sessionNameOrId).Run()
}

func (c *Client) HasSession(sessionName string) bool {
	_, err := exec.Command("tmux", "has-session", "-t", sessionName).Output()
	if err != nil {
		return false
	}
	return true
}

func (c *Client) GetSessionId(sessionName string) string {
	return execAndPrint("tmux", "display-message", "-t", sessionName, "-p", "#{session_id}")
}

