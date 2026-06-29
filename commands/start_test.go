package commands

import (
	"fmt"
	"os"
	"testing"
)

func Test_StartResumesExistingSession(t *testing.T) {
	mockGit := &MockGit{SessionIdToReturn: "sessionId"}
	mockTmux := &MockTmux{SessionIdToReturn: "sessionId", SessionExists: true}

	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Path: "/.worktrees/path/to/worktree"}, nil
		},
	}

	commander.Start([]string{"-s", "ENG123"})

	if mockGit.CreateWorktreeCalled {
		t.Error("Called CreateWorktree even though exists")
	}

	if !mockGit.SaveSessionIdCalled {
		t.Error("Did not save session id")
	}

	if mockGit.SavedSessionId != "sessionId" {
		t.Error("Incorrect sessionId saved")
	}
}

func Test_StartAttachesTmuxSessionWhenOutside(t *testing.T) {
	os.Setenv("TMUX", "")
	mockGit := &MockGit{SessionIdToReturn: "sessionId"}
	mockTmux := &MockTmux{SessionIdToReturn: "sessionId", SessionExists: true}

	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Path: "/.worktrees/path/to/worktree"}, nil
		},
	}

	commander.Start([]string{"-s", "ENG123"})

	if !mockTmux.AttachCalled {
		t.Error("Attach was not called")
	}
	if mockTmux.SwitchCalled {
		t.Error("Switch was called")
	}
}

func Test_StartSwitchesTmuxSessionWhenInside(t *testing.T) {
	os.Setenv("TMUX", "true")
	mockGit := &MockGit{SessionIdToReturn: "sessionId"}
	mockTmux := &MockTmux{SessionIdToReturn: "sessionId", SessionExists: true}

	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Path: "/.worktrees/path/to/worktree"}, nil
		},
	}

	commander.Start([]string{"-s", "ENG123"})

	if !mockTmux.SwitchCalled {
		t.Error("Switch was not called")
	}

	if mockTmux.AttachCalled {
		t.Error("Attach was called")
	}
}


func Test_StartCreatesWorktree(t *testing.T) {
	mockGit := &MockGit{SessionIdToReturn: "sessionId"}
	mockTmux := &MockTmux{SessionIdToReturn: "sessionId", SessionExists: true}

	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Branch: "some-branch"}, fmt.Errorf("Picker error")
		},
	}

	commander.Start([]string{"-s", "ENG123"})
	
	if !mockGit.CreateWorktreeCalled {
		t.Error("Create worktree was not called")
	}
}

