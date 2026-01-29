package pmcombo

import "fmt"

// WhyCommand generates why command for dependency explanation
type WhyCommand struct {
	Options WhyOptions
}

// NewWhyCommand creates a new why command
func NewWhyCommand(opts WhyOptions) *WhyCommand {
	return &WhyCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *WhyCommand) Concat(pm PackageManager) ([]string, error) {
	if c.Options.Package == "" {
		return nil, fmt.Errorf("no package specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "why", c.Options.Package)

	case Yarn:
		args = append(args, "why", c.Options.Package)

	case YarnClassic:
		args = append(args, "why", c.Options.Package)

	case PNPM:
		args = append(args, "why", c.Options.Package)

	case Deno:
		return nil, fmt.Errorf("deno does not support why command")

	case Bun:
		// Bun doesn't have why, but we can try
		args = append(args, "pm", "ls", "--all")
	}

	return args, nil
}
