package pmcombo

import "fmt"

// RunCommand generates run script command
type RunCommand struct {
	Options RunOptions
}

// NewRunCommand creates a new run command
func NewRunCommand(opts RunOptions) *RunCommand {
	return &RunCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *RunCommand) Concat(pm PackageManager) ([]string, error) {
	if c.Options.Script == "" {
		return nil, fmt.Errorf("no script specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "run", c.Options.Script)
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case Yarn:
		args = append(args, "run", c.Options.Script)
		args = append(args, c.Options.Args...)

	case YarnClassic:
		args = append(args, "run", c.Options.Script)
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case PNPM:
		args = append(args, "run", c.Options.Script)
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case Deno:
		args = append(args, "task", c.Options.Script)
		args = append(args, c.Options.Args...)

	case Bun:
		args = append(args, "run", c.Options.Script)
		args = append(args, c.Options.Args...)
	}

	return args, nil
}
