package pmcombo

import "fmt"

// RemoveCommand generates remove package command
type RemoveCommand struct {
	Options RemoveOptions
}

// NewRemoveCommand creates a new remove command
func NewRemoveCommand(opts RemoveOptions) *RemoveCommand {
	return &RemoveCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *RemoveCommand) Concat(pm PackageManager) ([]string, error) {
	if len(c.Options.Packages) == 0 {
		return nil, fmt.Errorf("no packages specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "uninstall")
		if c.Options.Global {
			args = append(args, "-g")
		}
		args = append(args, c.Options.Packages...)

	case Yarn:
		args = append(args, "remove")
		args = append(args, c.Options.Packages...)

	case YarnClassic:
		args = append(args, "remove")
		args = append(args, c.Options.Packages...)

	case PNPM:
		args = append(args, "remove")
		if c.Options.Global {
			args = append(args, "-g")
		}
		args = append(args, c.Options.Packages...)

	case Deno:
		args = append(args, "remove")
		args = append(args, c.Options.Packages...)

	case Bun:
		args = append(args, "remove")
		if c.Options.Global {
			args = append(args, "-g")
		}
		args = append(args, c.Options.Packages...)
	}

	return args, nil
}
