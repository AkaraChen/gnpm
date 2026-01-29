package context

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/user/fnpm/internal/pmcombo"
)

// ProjectContext holds information about the current project
type ProjectContext struct {
	RootDir        string
	PackageJSON    *PackageJSON
	PackageManager pmcombo.PackageManager
	IsWorkspace    bool
}

// Detect detects the project context from the given directory
func Detect(startDir string) (*ProjectContext, error) {
	pkgPath, err := FindPackageJSON(startDir)
	if err != nil {
		return nil, err
	}

	rootDir := filepath.Dir(pkgPath)
	pkg, err := ReadPackageJSON(pkgPath)
	if err != nil {
		return nil, err
	}

	pm := DetectPackageManager(rootDir, pkg)

	return &ProjectContext{
		RootDir:        rootDir,
		PackageJSON:    pkg,
		PackageManager: pm,
		IsWorkspace:    pkg.HasWorkspaces(),
	}, nil
}

// DetectPackageManager detects the package manager from lock files and package.json
func DetectPackageManager(rootDir string, pkg *PackageJSON) pmcombo.PackageManager {
	// Check lock files in order of specificity
	lockFiles := []struct {
		name string
		pm   pmcombo.PackageManager
	}{
		{"bun.lockb", pmcombo.Bun},
		{"bun.lock", pmcombo.Bun},
		{"deno.lock", pmcombo.Deno},
		{"pnpm-lock.yaml", pmcombo.PNPM},
		{"yarn.lock", pmcombo.Yarn}, // Will check for classic vs berry below
		{"package-lock.json", pmcombo.NPM},
		{"npm-shrinkwrap.json", pmcombo.NPM},
	}

	for _, lf := range lockFiles {
		lockPath := filepath.Join(rootDir, lf.name)
		if _, err := os.Stat(lockPath); err == nil {
			// Special handling for yarn - check if it's classic or berry
			if lf.pm == pmcombo.Yarn {
				return detectYarnVersion(rootDir)
			}
			return lf.pm
		}
	}

	// Fallback to packageManager field in package.json
	if pkg != nil && pkg.PackageManager != "" {
		return parsePackageManagerField(pkg.PackageManager)
	}

	// Default to npm
	return pmcombo.NPM
}

// detectYarnVersion checks if the project uses Yarn Classic or Yarn Berry
func detectYarnVersion(rootDir string) pmcombo.PackageManager {
	// Check for .yarnrc.yml (Yarn Berry)
	yarnrcPath := filepath.Join(rootDir, ".yarnrc.yml")
	if _, err := os.Stat(yarnrcPath); err == nil {
		return pmcombo.Yarn // Berry
	}

	// Check for .yarnrc (Yarn Classic)
	yarnrcClassicPath := filepath.Join(rootDir, ".yarnrc")
	if _, err := os.Stat(yarnrcClassicPath); err == nil {
		return pmcombo.YarnClassic
	}

	// Check yarn.lock content for clues
	yarnLockPath := filepath.Join(rootDir, "yarn.lock")
	content, err := os.ReadFile(yarnLockPath)
	if err == nil {
		// Yarn Berry lock files start with specific header
		if strings.Contains(string(content), "__metadata:") {
			return pmcombo.Yarn // Berry
		}
	}

	// Default to Yarn Classic for older projects
	return pmcombo.YarnClassic
}

// parsePackageManagerField parses the packageManager field from package.json
// Format: "npm@10.2.0", "pnpm@8.6.0", etc.
func parsePackageManagerField(field string) pmcombo.PackageManager {
	// Extract package manager name before @
	name := field
	if idx := strings.Index(field, "@"); idx != -1 {
		name = field[:idx]
	}

	switch strings.ToLower(name) {
	case "npm":
		return pmcombo.NPM
	case "yarn":
		return pmcombo.Yarn
	case "pnpm":
		return pmcombo.PNPM
	case "bun":
		return pmcombo.Bun
	case "deno":
		return pmcombo.Deno
	default:
		return pmcombo.NPM
	}
}

// FindWorkspaceRoot finds the root of a workspace/monorepo
func FindWorkspaceRoot(startDir string) (string, error) {
	dir := startDir
	var lastWorkspaceRoot string

	for {
		pkgPath := filepath.Join(dir, "package.json")
		if _, err := os.Stat(pkgPath); err == nil {
			pkg, err := ReadPackageJSON(pkgPath)
			if err == nil && pkg.HasWorkspaces() {
				lastWorkspaceRoot = dir
			}
		}

		// Also check for pnpm-workspace.yaml
		pnpmWorkspacePath := filepath.Join(dir, "pnpm-workspace.yaml")
		if _, err := os.Stat(pnpmWorkspacePath); err == nil {
			lastWorkspaceRoot = dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if lastWorkspaceRoot != "" {
		return lastWorkspaceRoot, nil
	}

	return "", os.ErrNotExist
}
