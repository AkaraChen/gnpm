package pmcombo

import "fmt"

// DlxCommand generates download and execute command
type DlxCommand struct {
	Options DlxOptions
}

// NewDlxCommand creates a new dlx command
func NewDlxCommand(opts DlxOptions) *DlxCommand {
	return &DlxCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *DlxCommand) Concat(pm PackageManager) ([]string, error) {
	if c.Options.Package == "" {
		return nil, fmt.Errorf("no package specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "exec", "--", c.Options.Package)
		args = append(args, c.Options.Args...)

	case Yarn:
		args = append(args, "dlx", c.Options.Package)
		args = append(args, c.Options.Args...)

	case YarnClassic:
		// Yarn Classic doesn't have dlx, approximate with npx
		args = append(args, "exec", c.Options.Package)
		args = append(args, c.Options.Args...)

	case PNPM:
		args = append(args, "dlx", c.Options.Package)
		args = append(args, c.Options.Args...)

	case Deno:
		args = append(args, "run", "-A", "npm:"+c.Options.Package)
		args = append(args, c.Options.Args...)

	case Bun:
		args = append(args, "x", c.Options.Package)
		args = append(args, c.Options.Args...)
	}

	return args, nil
}
