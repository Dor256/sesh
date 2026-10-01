package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Pane struct {
	Name string
	Path string
}

func (commander *Commander) Load(args []string) error {
	if os.Getenv("TMUX") != "" {
		fmt.Println("Can't load tmux state from within tmux")
		return nil
	}
	userConfigDir, err := userConfigDir()
	if err != nil {
		return err
	}
	seshStatePath := filepath.Join(userConfigDir, "sesh", "state")
	rawTmuxState, err := os.ReadFile(seshStatePath)
	if err != nil {
		fmt.Println("Couldn't read Sesh state")
		return err
	}
	tmuxState := parseTmuxState(string(rawTmuxState))
	if tmuxState == nil {
		// Tmux state is clean
		return nil
	}
	var worktreePath string
	var sessionId string
	for session, panes := range tmuxState {
		for idx, pane := range panes {
			// Worktree path is the same for each session
			worktreePath = pane.Path
			if idx == 0 {
				sessionId = commander.tmuxClient.Create(session, pane.Path)
				commander.tmuxClient.RenameWindow(session, "1", pane.Name)
			} else {
				commander.tmuxClient.NewWindow(session, pane.Name, pane.Path)
			}
		}
	}

	commander.gitClient.SaveSessionId(worktreePath, sessionId)
	commander.tmuxClient.AttachDefault()
	return nil
}

func parseTmuxState(rawState string) map[string][]Pane {
	rawSessions := strings.Split(rawState, "\n")
	state := make(map[string][]Pane)
	for _, rawSession := range rawSessions {
		sessionTuple := strings.Split(rawSession, " ")
		if len(sessionTuple) == 0 {
			return nil
		}
		sessionName := sessionTuple[0]
		pane := Pane{Name: sessionTuple[1], Path: sessionTuple[2]}
		if state[sessionName] == nil {
			state[sessionName] = []Pane{pane}
		} else {
			state[sessionName] = append(state[sessionName], pane)
		}
	}
	return state
}
