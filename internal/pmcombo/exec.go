package pmcombo

import "fmt"

// ExecCommand generates exec command for running binaries
type ExecCommand struct {
	Options ExecOptions
}

// NewExecCommand creates a new exec command
func NewExecCommand(opts ExecOptions) *ExecCommand {
	return &ExecCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *ExecCommand) Concat(pm PackageManager) ([]string, error) {
	if c.Options.Command == "" {
		return nil, fmt.Errorf("no command specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "exec", "--", c.Options.Command)
		args = append(args, c.Options.Args...)

	case Yarn:
		args = append(args, "exec", c.Options.Command)
		args = append(args, c.Options.Args...)

	case YarnClassic:
		// Yarn Classic doesn't have exec, use npx-like behavior
		args = append(args, "run", c.Options.Command)
		args = append(args, c.Options.Args...)

	case PNPM:
		args = append(args, "exec", c.Options.Command)
		args = append(args, c.Options.Args...)

	case Deno:
		// Deno doesn't have direct exec equivalent
		args = append(args, "run", "-A", "npm:"+c.Options.Command)
		args = append(args, c.Options.Args...)

	case Bun:
		args = append(args, "x", c.Options.Command)
		args = append(args, c.Options.Args...)
	}

	return args, nil
}
