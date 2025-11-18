package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := newFlagSet("lookctl")

	if err := parseFlag(fs, args, printMainHelp); err != nil {
		return err
	}

	tw := newTabWriter(os.Stderr)

	if fs.NArg() == 0 {
		printMainHelp(tw)
		return fmt.Errorf("please specify a command")
	}

	cmd := fs.Arg(0)
	cmdArgs := fs.Args()[1:]

	var err error
	switch cmd {
	case "list":
		err = list(cmdArgs)
	case "current":
		err = current(cmdArgs)
	case "set":
		err = set(cmdArgs)
	default:
		return fmt.Errorf("unknown command: '%s'. see 'lookctl -h' for more information", cmd)
	}

	if err != nil {
		return err
	}

	return nil
}
