package commands

import (
	"os"
	"path/filepath"
)

func (commander *Commander) Save(args []string) error {
	userConfigDir, err := userConfigDir()
	if err != nil {
		return err
	}
	seshStatePath := filepath.Join(userConfigDir, "sesh", "state")
	tmuxState := commander.tmuxClient.ListPanes()
	
	os.WriteFile(seshStatePath, []byte(tmuxState), 0644)
	
	return nil
}

