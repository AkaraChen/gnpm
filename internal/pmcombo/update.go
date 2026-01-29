package pmcombo

// UpdateCommand generates update packages command
type UpdateCommand struct {
	Options UpdateOptions
}

// NewUpdateCommand creates a new update command
func NewUpdateCommand(opts UpdateOptions) *UpdateCommand {
	return &UpdateCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *UpdateCommand) Concat(pm PackageManager) ([]string, error) {
	var args []string

	switch pm {
	case NPM:
		args = append(args, "update")
		args = append(args, c.Options.Packages...)

	case Yarn:
		if c.Options.Interactive {
			args = append(args, "upgrade-interactive")
		} else {
			args = append(args, "up")
		}
		args = append(args, c.Options.Packages...)

	case YarnClassic:
		if c.Options.Interactive {
			args = append(args, "upgrade-interactive")
			if c.Options.Latest {
				args = append(args, "--latest")
			}
		} else {
			args = append(args, "upgrade")
		}
		args = append(args, c.Options.Packages...)

	case PNPM:
		if c.Options.Interactive {
			args = append(args, "update", "-i")
		} else {
			args = append(args, "update")
		}
		if c.Options.Latest {
			args = append(args, "--latest")
		}
		args = append(args, c.Options.Packages...)

	case Deno:
		args = append(args, "cache", "--reload")
		args = append(args, c.Options.Packages...)

	case Bun:
		args = append(args, "update")
		args = append(args, c.Options.Packages...)
	}

	return args, nil
}
