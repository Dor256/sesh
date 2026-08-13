package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func (commander *Commander) Setup(args []string) error {
	configDir, err := userConfigDir()
	if err != nil {
		return err
	}
	fmt.Println("Setting up Sesh...")
	seshDir := filepath.Join(configDir, "sesh")
	seshStateFile := filepath.Join(seshDir, "state")
	if _, err := os.Stat(seshStateFile); err == nil {
		err = setupSeshStateFile(seshDir)
		if err != nil {
			return err
		}
	}

	tmuxConfigFile, err := getTmuxConfig()
	if err != nil {
		return err
	}
	err = setupTmuxHooks(tmuxConfigFile)
	if err != nil {
		return err
	}
	fmt.Println("Sesh is all set up")
	return nil
}

func setupSeshStateFile(seshDir string) error {
	os.Mkdir(seshDir, 0755)
	stateFilePath := filepath.Join(seshDir, "state")
	_, err := os.Create(stateFilePath)
	if err != nil {
		fmt.Println("Failed to create state file!")
		return err
	}
	return nil
}

func setupTmuxHooks(tmuxConfigFile string) error {
	clientDetachedHook := "set-hook -g client-detached 'run-shell \"sesh save\"'"
	sessionClosedHook := "set-hook -g session-closed 'run-shell \"sesh save\"'"
	sessionCreatedHook := "set-hook -g session-created 'run-shell \"sesh save\"'"
	grepErr := exec.Command("grep", "-qxF", clientDetachedHook + "\n" + sessionClosedHook + "\n" + sessionCreatedHook, tmuxConfigFile).Run()

	// Return if hooks are already set up
	if grepErr == nil {
		return nil
	}

	var exitErr *exec.ExitError
	file, err := os.OpenFile(tmuxConfigFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open tmux config file at %s. Err: %w", tmuxConfigFile, err)
	}
	defer file.Close()

	var stringBuilder strings.Builder
	if errors.As(grepErr, &exitErr) && exitErr.ExitCode() == 1 {
		stringBuilder.WriteString("\n# Sesh owned tmux hooks for session state recovery\n")
		stringBuilder.WriteString(clientDetachedHook);
		stringBuilder.WriteString("\n")
		stringBuilder.WriteString(sessionClosedHook)
		stringBuilder.WriteString("\n")
		stringBuilder.WriteString(sessionCreatedHook)
	}

	if _, err := file.WriteString(stringBuilder.String()); err != nil {
        return fmt.Errorf("failed to write to file %s: %w", tmuxConfigFile, err)
    }

	return nil
}

func getTmuxConfig() (string, error) {
	configDir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	tmuxConfigFile := filepath.Join(configDir, "tmux", "tmux.conf")
	_, err = os.Stat(tmuxConfigFile)
	if err == nil {
		return tmuxConfigFile, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	legacyConfigPath := filepath.Join(home, ".tmux.conf")
	_, err = os.Stat(legacyConfigPath)
	if err == nil {
		return legacyConfigPath, nil
	}

	err = os.MkdirAll(filepath.Join(configDir, "tmux"), 0755)
	if err != nil {
		return "", err
	}
	_, err = os.Create(tmuxConfigFile)
	if err != nil {
		return "", err
	}
	return tmuxConfigFile, nil
}

func userConfigDir() (string, error) {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return xdg, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't extract user's HOME dir")
		return "", err
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		return filepath.Join(home, ".config"), nil
	}
	return "", fmt.Errorf("No config dir found!")
}
