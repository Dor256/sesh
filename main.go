package main

import (
	"fmt"
	"os"
	"sesh/commands"
)

func main() {
	userArgs := os.Args[1:]
	command := userArgs[0]

	switch command {
	case "destroy":
		commands.Destroy(userArgs[1:])
	case "start":
		commands.Start(userArgs[1:])
	default:
		fmt.Fprintf(os.Stderr, "No such argument %s\n", command)
		os.Exit(1)
	}
}
