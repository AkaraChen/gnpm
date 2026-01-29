package pmcombo

// InitCommand generates init command for creating package.json
type InitCommand struct {
	Options InitOptions
}

// NewInitCommand creates a new init command
func NewInitCommand(opts InitOptions) *InitCommand {
	return &InitCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *InitCommand) Concat(pm PackageManager) ([]string, error) {
	var args []string

	switch pm {
	case NPM:
		args = append(args, "init")
		if c.Options.Yes {
			args = append(args, "-y")
		}

	case Yarn:
		args = append(args, "init")
		if c.Options.Yes {
			args = append(args, "-y")
		}

	case YarnClassic:
		args = append(args, "init")
		if c.Options.Yes {
			args = append(args, "-y")
		}

	case PNPM:
		args = append(args, "init")

	case Deno:
		args = append(args, "init")

	case Bun:
		args = append(args, "init")
		if c.Options.Yes {
			args = append(args, "-y")
		}
	}

	return args, nil
}
