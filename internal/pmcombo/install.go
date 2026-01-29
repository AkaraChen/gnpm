package pmcombo

// InstallCommand generates install command for all dependencies
type InstallCommand struct {
	Options InstallOptions
}

// NewInstallCommand creates a new install command
func NewInstallCommand(opts InstallOptions) *InstallCommand {
	return &InstallCommand{Options: opts}
}

// Concat generates the command arguments for the given package manager
func (c *InstallCommand) Concat(pm PackageManager) ([]string, error) {
	switch pm {
	case NPM:
		if c.Options.Frozen {
			return []string{"ci"}, nil
		}
		return []string{"install"}, nil

	case Yarn:
		if c.Options.Frozen {
			return []string{"install", "--immutable"}, nil
		}
		return []string{"install"}, nil

	case YarnClassic:
		if c.Options.Frozen {
			return []string{"install", "--frozen-lockfile"}, nil
		}
		return []string{"install"}, nil

	case PNPM:
		if c.Options.Frozen {
			return []string{"install", "--frozen-lockfile"}, nil
		}
		return []string{"install"}, nil

	case Deno:
		return []string{"install"}, nil

	case Bun:
		if c.Options.Frozen {
			return []string{"install", "--frozen-lockfile"}, nil
		}
		return []string{"install"}, nil

	default:
		return []string{"install"}, nil
	}
}
