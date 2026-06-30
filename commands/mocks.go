package commands

type MockGit struct {
	ReadSessionIdCalled bool
	DestroyedWorktree   string
	SessionIdToReturn   string
	SaveSessionIdCalled bool
	SavedSessionId string
	WorktreePath string
	CreateWorktreeCalled bool
}

func (m *MockGit) ReadSessionId(path string) string {
	m.ReadSessionIdCalled = true
	return m.SessionIdToReturn
}

func (m *MockGit) DestroyWorktree(path string) error {
	m.DestroyedWorktree = path
	return nil
}

func (m *MockGit) CreateWorktree(branchName, path string) error {
	m.CreateWorktreeCalled = true
	return nil
}

func (m *MockGit) SaveSessionId(path, tmuxSessionId string) error {
	m.SaveSessionIdCalled = true
	m.SavedSessionId = tmuxSessionId
	m.WorktreePath = path
	return nil
}


type MockTmux struct {
	SessionIdToReturn string
	KilledSessionId string
	SessionExists bool
	AttachCalled bool
	SwitchCalled bool
	CreatedSessionName string
	CreatedWorktreePath string
	RenameWindowCalled bool
	NewWindowCalledTimes int
}

func (m *MockTmux) Kill(sessionId string) error {
	m.KilledSessionId = sessionId
	return nil
}

func (m *MockTmux) Attach(sessionName string) error {
	m.AttachCalled = true
	return nil
}

func (m *MockTmux) Switch(sessionName string) error {
	m.SwitchCalled = true
	return nil
}

func (m *MockTmux) Create(sessionName, path string) string {
	m.CreatedSessionName = sessionName
	m.CreatedWorktreePath = path
	return m.SessionIdToReturn
}

func (m *MockTmux) NewWindow(sessionName, windowName, path string, tool ...string) error {
	m.NewWindowCalledTimes += 1
	return nil
}

func (m *MockTmux) RenameWindow(sessionName, oldName, newName string) error {
	m.RenameWindowCalled = true
	return nil
}

func (m *MockTmux) HasSession(sessionName string) bool { return m.SessionExists }

func (m *MockTmux) GetSessionId(sessionName string) string { return m.SessionIdToReturn }
