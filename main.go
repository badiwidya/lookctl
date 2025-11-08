package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrUsage = errors.New("")

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, ErrUsage) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printMainHelp(os.Stderr)
		return ErrUsage
	}

	cmd := os.Args[1]
	cmdArgs := os.Args[2:]

	var err error
	switch cmd {
	case "help":
		err = help(cmdArgs)
	case "list":
		err = list(cmdArgs)
	case "current":
		err = current(cmdArgs)
	}

	if err != nil {
		return err
	}

	return nil
}
