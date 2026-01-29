package pmcombo

import "fmt"

// CreateCommand generates create from template command
type CreateCommand struct {
	Options CreateOptions
}

// NewCreateCommand creates a new create command
func NewCreateCommand(opts CreateOptions) *CreateCommand {
	return &CreateCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *CreateCommand) Concat(pm PackageManager) ([]string, error) {
	if c.Options.Template == "" {
		return nil, fmt.Errorf("no template specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "create", c.Options.Template)
		args = append(args, c.Options.Args...)

	case Yarn:
		args = append(args, "create", c.Options.Template)
		args = append(args, c.Options.Args...)

	case YarnClassic:
		args = append(args, "create", c.Options.Template)
		args = append(args, c.Options.Args...)

	case PNPM:
		args = append(args, "create", c.Options.Template)
		args = append(args, c.Options.Args...)

	case Deno:
		// Deno doesn't have create, use init or npm equivalent
		args = append(args, "run", "-A", "npm:create-"+c.Options.Template)
		args = append(args, c.Options.Args...)

	case Bun:
		args = append(args, "create", c.Options.Template)
		args = append(args, c.Options.Args...)
	}

	return args, nil
}
