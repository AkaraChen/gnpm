package pmcombo

import "fmt"

// PackageManager represents supported package managers
type PackageManager string

const (
	NPM         PackageManager = "npm"
	Yarn        PackageManager = "yarn"
	YarnClassic PackageManager = "yarn-classic"
	PNPM        PackageManager = "pnpm"
	Deno        PackageManager = "deno"
	Bun         PackageManager = "bun"
)

// AllPackageManagers returns all supported package managers
func AllPackageManagers() []PackageManager {
	return []PackageManager{NPM, Yarn, YarnClassic, PNPM, Deno, Bun}
}

// String returns the string representation of the package manager
func (pm PackageManager) String() string {
	return string(pm)
}

// Executable returns the command-line executable name
func (pm PackageManager) Executable() string {
	switch pm {
	case YarnClassic:
		return "yarn"
	default:
		return string(pm)
	}
}

// ParsePackageManager parses a string into a PackageManager
func ParsePackageManager(s string) (PackageManager, error) {
	switch s {
	case "npm":
		return NPM, nil
	case "yarn":
		return Yarn, nil
	case "yarn-classic":
		return YarnClassic, nil
	case "pnpm":
		return PNPM, nil
	case "deno":
		return Deno, nil
	case "bun":
		return Bun, nil
	default:
		return "", fmt.Errorf("unknown package manager: %s", s)
	}
}

// Command interface for all pm-combo commands
type Command interface {
	// Concat generates the command arguments for the given package manager
	Concat(pm PackageManager) ([]string, error)
}

// AddOptions for the add command
type AddOptions struct {
	Packages []string
	Dev      bool
	Exact    bool
	Global   bool
	Peer     bool
	Optional bool
}

// RemoveOptions for the remove command
type RemoveOptions struct {
	Packages []string
	Global   bool
}

// InstallOptions for the install command
type InstallOptions struct {
	Frozen bool // For CI, use frozen lockfile
}

// InitOptions for the init command
type InitOptions struct {
	Yes bool // Skip prompts
}

// RunOptions for the run command
type RunOptions struct {
	Script string
	Args   []string
}

// ExecOptions for the exec command
type ExecOptions struct {
	Command string
	Args    []string
}

// TestOptions for the test command
type TestOptions struct {
	Args []string
}

// CreateOptions for the create command
type CreateOptions struct {
	Template string
	Args     []string
}

// DlxOptions for the dlx command
type DlxOptions struct {
	Package string
	Args    []string
}

// UpdateOptions for the update command
type UpdateOptions struct {
	Packages    []string
	Interactive bool
	Latest      bool
}

// ConfigOptions for the config command
type ConfigOptions struct {
	Action string // get, set, delete, list
	Key    string
	Value  string
}

// WhyOptions for the why command
type WhyOptions struct {
	Package string
}
