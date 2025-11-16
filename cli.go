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
		return fmt.Errorf("'help' only accepts one argument")
	}

	switch args[0] {
	case "list":
		printListHelp(os.Stdout)
	case "set":
		printSetHelp(os.Stdout)
	default:
		return fmt.Errorf("subcommand '%s' not found", args[0])
	}

	return nil
}

func list(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("'list' only accepts one argument")
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
		return fmt.Errorf("see 'lookctl help list' for more informations")
	}

	for i, asset := range assets {
		fmt.Fprintf(os.Stdout, "[%d] %s\n", i+1, asset)
	}

	return nil
}

func current(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("'current' takes no argument")
	}

	currentTheme, err := getCurrentTheme()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "GTK Theme	: %s\n", currentTheme.gtkTheme)
	fmt.Fprintf(os.Stdout, "Icon Theme	: %s\n", currentTheme.iconTheme)
	fmt.Fprintf(os.Stdout, "Cursor Theme	: %s\n", currentTheme.cursorTheme)

	return nil
}

func set(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("arguments cannot be empty. see 'lookctl help set' for more informations")
	}

	if len(args) != 2 && args[0] != "help" {
		return fmt.Errorf("invalid arguments. see 'lookctl help set' for more informations")
	}

	arg := args[0]

	var err error
	switch arg {
	case "theme":
		err = setTheme(args[1])
	case "cursor":
		err = setCursorTheme(args[1])
	case "icon":
		err = setIconTheme(args[1])
	case "help":
		printSetHelp(os.Stdout)
		return nil
	default:
		return fmt.Errorf("invalid arguments. see 'lookctl help set' for more informations")
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "theme changed successfully!")

	return nil
}
