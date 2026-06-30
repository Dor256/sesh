package commands

type TmuxClient interface {
	Create(sessionName, worktreePath string) string
	HasSession(sessionName string) bool
	GetSessionId(sessionName string) string
	Switch(sessionName string) error
	Attach(sessionName string) error
	RenameWindow(sessionName, windowIndex, windowName string) error
	NewWindow(sessionName, windowName, worktreePath string, command ...string) error
	Kill(sessionId string) error
}

type GitClient interface {
	SaveSessionId(worktreePath, sessionId string) error
	CreateWorktree(branchName, worktreePath string) error
	ReadSessionId(worktreePath string) string
	DestroyWorktree(worktreePath string) error
}

type Commander struct {
	gitClient  GitClient
	tmuxClient TmuxClient
	picker func() (*Worktree, error)
}

func NewCommander(tmuxClient TmuxClient, gitClient GitClient) *Commander {
	return &Commander{gitClient: gitClient, tmuxClient: tmuxClient, picker: Picker}
}
