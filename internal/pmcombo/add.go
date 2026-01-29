package pmcombo

import "fmt"

// AddCommand generates add package command
type AddCommand struct {
	Options AddOptions
}

// NewAddCommand creates a new add command
func NewAddCommand(opts AddOptions) *AddCommand {
	return &AddCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *AddCommand) Concat(pm PackageManager) ([]string, error) {
	if len(c.Options.Packages) == 0 {
		return nil, fmt.Errorf("no packages specified")
	}

	var args []string

	switch pm {
	case NPM:
		args = append(args, "install")
		args = c.appendNpmFlags(args)
		args = append(args, c.Options.Packages...)

	case Yarn:
		args = append(args, "add")
		args = c.appendYarnFlags(args)
		args = append(args, c.Options.Packages...)

	case YarnClassic:
		args = append(args, "add")
		args = c.appendYarnClassicFlags(args)
		args = append(args, c.Options.Packages...)

	case PNPM:
		args = append(args, "add")
		args = c.appendPnpmFlags(args)
		args = append(args, c.Options.Packages...)

	case Deno:
		args = append(args, "add")
		// Deno uses npm: protocol for packages
		for _, pkg := range c.Options.Packages {
			args = append(args, "npm:"+pkg)
		}

	case Bun:
		args = append(args, "add")
		args = c.appendBunFlags(args)
		args = append(args, c.Options.Packages...)
	}

	return args, nil
}

func (c *AddCommand) appendNpmFlags(args []string) []string {
	if c.Options.Global {
		args = append(args, "-g")
	}
	if c.Options.Dev {
		args = append(args, "--save-dev")
	}
	if c.Options.Exact {
		args = append(args, "--save-exact")
	}
	if c.Options.Peer {
		args = append(args, "--save-peer")
	}
	if c.Options.Optional {
		args = append(args, "--save-optional")
	}
	return args
}

func (c *AddCommand) appendYarnFlags(args []string) []string {
	if c.Options.Dev {
		args = append(args, "-D")
	}
	if c.Options.Exact {
		args = append(args, "-E")
	}
	if c.Options.Peer {
		args = append(args, "--peer")
	}
	if c.Options.Optional {
		args = append(args, "-O")
	}
	return args
}

func (c *AddCommand) appendYarnClassicFlags(args []string) []string {
	if c.Options.Dev {
		args = append(args, "-D")
	}
	if c.Options.Exact {
		args = append(args, "-E")
	}
	if c.Options.Peer {
		args = append(args, "--peer")
	}
	if c.Options.Optional {
		args = append(args, "-O")
	}
	return args
}

func (c *AddCommand) appendPnpmFlags(args []string) []string {
	if c.Options.Global {
		args = append(args, "-g")
	}
	if c.Options.Dev {
		args = append(args, "-D")
	}
	if c.Options.Exact {
		args = append(args, "-E")
	}
	if c.Options.Peer {
		args = append(args, "--save-peer")
	}
	if c.Options.Optional {
		args = append(args, "-O")
	}
	return args
}

func (c *AddCommand) appendBunFlags(args []string) []string {
	if c.Options.Global {
		args = append(args, "-g")
	}
	if c.Options.Dev {
		args = append(args, "-d")
	}
	if c.Options.Exact {
		args = append(args, "--exact")
	}
	if c.Options.Peer {
		args = append(args, "--peer")
	}
	if c.Options.Optional {
		args = append(args, "--optional")
	}
	return args
}
