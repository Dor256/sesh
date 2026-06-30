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


func Test_StartCreatesTmuxSessionAndOpensWindows(t *testing.T) {
	mockGit := &MockGit{SessionIdToReturn: "sessionId"}
	mockTmux := &MockTmux{
		SessionIdToReturn: "sessionId",
		SessionExists: false,
		NewWindowCalledTimes: 0,
	}
	expectedSessionName := "ENG123"
	expectedWorktreePath := ".workrees/path/to/worktree"

	commander := Commander{
		gitClient: mockGit,
		tmuxClient: mockTmux,
		picker: func() (*Worktree, error) {
			return &Worktree{Path: expectedWorktreePath, Branch: "some-branch"}, nil
		},
	}

	commander.Start([]string{"-s", expectedSessionName})

	if mockTmux.CreatedSessionName != expectedSessionName {
		t.Errorf("Expected %s, Got: %s", expectedSessionName, mockTmux.CreatedSessionName)
	}	

	if mockTmux.CreatedWorktreePath != expectedWorktreePath {
		t.Errorf("Expected: %s, Got: %s", expectedWorktreePath, mockTmux.CreatedWorktreePath)
	}

	if !mockTmux.RenameWindowCalled {
		t.Error("Expectd RenameWindow to have been called")
	}

	if mockTmux.NewWindowCalledTimes != 3 {
		t.Errorf("Expected NewWindow to have been called %d times, Got: %d", 3, mockTmux.NewWindowCalledTimes)
	}
}

