package main

import (
	"fmt"
	"os"
)

func list(args []string) error {
	fs := newFlagSet("list")

	showGtk := fs.Bool("gtk", false, "Show installed gtk themes")
	showIcon := fs.Bool("icon", false, "Show installed icon themes")
	showCursor := fs.Bool("cursor", false, "Show installed cursor themes")

	if err := parseFlag(fs, args, printListHelp); err != nil {
		return err
	}

	if fs.NArg() != 0 {
		return fmt.Errorf("'list' does not accept arguments; use flags instead")
	}

	if fs.NFlag() == 0 {
		*showGtk = true
	}

	tw := newTabWriter(os.Stdout)

	if *showGtk {
		fmt.Fprintf(tw, "GTK Themes:\n")

		gtkThemes := getInstalledThemes()

		for i, theme := range gtkThemes {
			fmt.Fprintf(tw, "\t[%d]\t%s\n", i+1, theme)
		}
	}

	if *showIcon {
		fmt.Fprintf(tw, "Icon Themes:\n")

		iconThemes := getInstalledIconThemes()

		for i, theme := range iconThemes {
			fmt.Fprintf(tw, "\t[%d]\t%s\n", i+1, theme)
		}
	}

	if *showCursor {
		fmt.Fprintf(tw, "Cursor Themes:\n")

		cursorThemes := getInstalledCursorThemes()

		for i, theme := range cursorThemes {
			fmt.Fprintf(tw, "\t[%d]\t%s\n", i+1, theme)
		}
	}

	tw.Flush()

	return nil
}

func current(args []string) error {
	fs := newFlagSet("current")

	if err := parseFlag(fs, args, printCurrentHelp); err != nil {
		return err
	}

	if fs.NFlag() > 0 || fs.NArg() > 0 {
		return fmt.Errorf("'current' accepts no flags or arguments")
	}

	currentTheme, err := getCurrentTheme()
	if err != nil {
		return err
	}

	tw := newTabWriter(os.Stdout)

	fmt.Fprintf(tw, "GTK Theme\t: %s\n", currentTheme.gtkTheme)
	fmt.Fprintf(tw, "Icon Theme\t: %s\n", currentTheme.iconTheme)
	fmt.Fprintf(tw, "Cursor Theme\t: %s\n", currentTheme.cursorTheme)

	colorScheme := "light"
	if currentTheme.preferDark {
		colorScheme = "dark"
	}

	fmt.Fprintf(tw, "Color Scheme\t: %s\n", colorScheme)

	tw.Flush()

	return nil
}

func set(args []string) error {
	fs := newFlagSet("set")

	gtkTheme := fs.String("gtk", "", "Set gtk theme")
	iconTheme := fs.String("icon", "", "Set icon theme")
	cursorTheme := fs.String("cursor", "", "Set cursor theme")
	colorScheme := fs.String("color-scheme", "", "Manually set color scheme")

	if err := parseFlag(fs, args, printSetHelp); err != nil {
		return err
	}

	if fs.NArg() != 0 {
		return fmt.Errorf("'set' does not accept arguments; use flags instead")
	}

	if fs.NFlag() == 0 {
		return fmt.Errorf("please specify one or more flags")
	}

	currentCfg, err := getCurrentTheme()
	if err != nil {
		return err
	}

	if *gtkTheme != "" {
		if err := setTheme(&currentCfg, *gtkTheme); err != nil {
			return err
		}
	}

	if *iconTheme != "" {
		if err := setIconTheme(&currentCfg, *iconTheme); err != nil {
			return err
		}
	}

	if *cursorTheme != "" {
		if err := setCursorTheme(&currentCfg, *cursorTheme); err != nil {
			return err
		}
	}

	if *colorScheme != "" {
		if err := setColorScheme(&currentCfg, *colorScheme); err != nil {
			return err
		}
	}

	if err := saveCurrentTheme(currentCfg); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "changes saved successfully!\n")

	return nil
}
