package main

import (
	"fmt"
	"os"
)

func help(args []string) error {
	if len(args) == 0 {
		printMainHelp(os.Stdout)
		return nil
	}

	if len(args) > 1 {
		return fmt.Errorf("error: 'help' only accepts one argument")
	}

	switch args[0] {
	case "list":
		printListHelp(os.Stdout)
	default:
		return fmt.Errorf("error: subcommand '%s' not found", args[0])
	}

	return nil
}

func list(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("error: 'list' only accepts one argument")
	}

	arg := ""

	if len(args) == 0 {
		arg = "theme"
	} else {
		arg = args[0]
	}

	assets := []string{}

	switch arg {
	case "theme":
		assets = getInstalledThemes()
	case "cursor":
		assets = getInstalledCursorThemes()
	case "icon":
		assets = getInstalledIconThemes()
	case "help":
		printListHelp(os.Stdout)
		return nil
	default:
		return fmt.Errorf("error: see 'lookctl help list' for more informations")
	}

	for i, asset := range assets {
		fmt.Fprintf(os.Stdout, "[%d] %s\n", i+1, asset)
	}

	return nil
}
