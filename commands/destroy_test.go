package commands

import (
	"fmt"
	"testing"
)

func Test_DestroysWorktreeAndSession(t *testing.T) {
	mockGit := &MockGit{SessionIdToReturn: "blamos"}
	mockTmux := &MockTmux{}
	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Path: "/.worktrees/path/to/worktree"}, nil
		},
	}

	commander.Destroy([]string{})

	if mockGit.ReadSessionIdCalled != true {
		t.Error("expected ReadSessionIdCalled to be true got false")
	}
	if mockGit.SessionIdToReturn != "blamos" {
		t.Error("Got wrong session id")
	}
	if mockTmux.KilledSessionId != "blamos" {
		t.Error("Destroyed wrong sessionId")
	}
}

func Test_AbortsWhenNoWorktreeSelected(t *testing.T) {
	mockGit := &MockGit{SessionIdToReturn: "blamos"}
	mockTmux := &MockTmux{}
	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return nil, fmt.Errorf("%s", "Aborted")
		},
	}

	commander.Destroy([]string{})

	if mockGit.ReadSessionIdCalled == true {
		t.Error("Destroy was not aborted")
	}
	if mockTmux.KilledSessionId != "" {
		t.Error("Tmux session was killed")
	}
}
