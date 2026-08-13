package main

import (
	"fmt"
	"os"
	"sesh/commands"
	"sesh/git"
	"sesh/tmux"
)

func main() {
	userArgs := os.Args[1:]
	command := userArgs[0]
	tmuxClient := tmux.NewClient()
	gitClient := git.NewClient()
	commander := commands.NewCommander(tmuxClient, gitClient)

	switch command {
	case "destroy":
		err := commander.Destroy(userArgs[1:])
		if err != nil {
			os.Exit(1)
		}
	case "start":
		err := commander.Start(userArgs[1:])
		if err != nil {
			os.Exit(1)
		}
	case "setup":
		err := commander.Setup(userArgs[1:])
		if err != nil {
			os.Exit(1)
		}
	case "load":
		err := commander.Load(userArgs[1:])
		if err != nil {
			os.Exit(1)
		}
	case "save":
		err := commander.Save(userArgs[1:])
		if err != nil {
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "No such argument %s\n", command)
		os.Exit(1)
	}
}
