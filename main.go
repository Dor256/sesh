package main

import (
	"fmt"
	"os"
	"regexp"
	"sesh/commands"
)


func main() {
	userArgs := os.Args[1:]
	command := userArgs[0]
	isNumber := regexp.MustCompile(`^[0-9]+$`)

	switch {
	case command == "destroy":
		commands.Destroy(userArgs[1:])
	case command == "start":
		commands.Start(userArgs[1:])
	case isNumber.MatchString(command):
		fmt.Println("ticket number was provided, spawn fzf for branch name")
	default:
		fmt.Fprintf(os.Stderr, "No such argument %s\n", command)
	}
}
