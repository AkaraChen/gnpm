package pmcombo

// TestCommand generates test command
type TestCommand struct {
	Options TestOptions
}

// NewTestCommand creates a new test command
func NewTestCommand(opts TestOptions) *TestCommand {
	return &TestCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *TestCommand) Concat(pm PackageManager) ([]string, error) {
	var args []string

	switch pm {
	case NPM:
		args = append(args, "test")
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case Yarn:
		args = append(args, "test")
		args = append(args, c.Options.Args...)

	case YarnClassic:
		args = append(args, "test")
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case PNPM:
		args = append(args, "test")
		if len(c.Options.Args) > 0 {
			args = append(args, "--")
			args = append(args, c.Options.Args...)
		}

	case Deno:
		args = append(args, "test")
		args = append(args, c.Options.Args...)

	case Bun:
		args = append(args, "test")
		args = append(args, c.Options.Args...)
	}

	return args, nil
}
