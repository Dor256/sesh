package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (commander *Commander) createWorktree(worktree Worktree, worktreePath string) error {
	if worktree.Branch == "" {
		fmt.Println("Aborting")
		return fmt.Errorf("Branch was not provided")
	}
	err := commander.gitClient.CreateWorktree(worktree.Branch, worktreePath)
	if err != nil {
		fmt.Printf("Git failed: %v\n", err)
		return err
	}
	return nil
}

func parseArgs(args []string) string {
	scanner := bufio.NewScanner(os.Stdin)
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	sessionName := startCmd.String("s", "", "Name of session (required)")
	startCmd.Parse(args)
	
	if *sessionName == "" {
		fmt.Print("Enter Session Name: ")
		if scanner.Scan() {
			input := scanner.Text()
			*sessionName = strings.TrimSpace(input)
		}
	}
	return *sessionName
}

func (commander *Commander) Start(args []string) error {
	sessionName := parseArgs(args)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't extract CWD. Aborting")
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't extract user's HOME dir")
		return err
	}

	worktreeDir := fmt.Sprintf("%s/dev/.worktrees", home)
	worktree, err := commander.picker()
	cwdBase := filepath.Base(cwd)
	var worktreePath string = worktree.Path

	// New worktree
	if err != nil {
		worktreePath = fmt.Sprintf("%s/%s-%s", worktreeDir, cwdBase, strings.ToLower(sessionName))
		err := commander.createWorktree(*worktree, worktreePath)
		if err != nil {
			return err
		}
	}
	var sessionId string
	if commander.tmuxClient.HasSession(sessionName) {
		sessionId = commander.tmuxClient.GetSessionId(sessionName)
	} else {
		sessionId = commander.tmuxClient.Create(sessionName, worktreePath)

		// Rename the default window
		commander.tmuxClient.RenameWindow(sessionName, "1", "Terminal")
		commander.tmuxClient.NewWindow(sessionName, "OpenCode", worktreePath, "opencode")
		commander.tmuxClient.NewWindow(sessionName, "Neovim", worktreePath, "nvim .")
		commander.tmuxClient.NewWindow(sessionName, "Claude", worktreePath, "claude")
	}

	commander.gitClient.SaveSessionId(worktreePath, sessionId)
	// If inside tmux switch to session if outside attach to session
	if os.Getenv("TMUX") != "" {
		commander.tmuxClient.Switch(sessionName)
	} else {
		commander.tmuxClient.Attach(sessionName)
	}
	return nil
}
