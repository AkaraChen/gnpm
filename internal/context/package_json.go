package context

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PackageJSON represents the relevant fields from package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	PackageManager  string            `json:"packageManager"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Workspaces      Workspaces        `json:"workspaces"`
}

// Workspaces can be either an array of strings or an object with packages field
type Workspaces struct {
	Patterns []string
}

// UnmarshalJSON handles both array and object formats for workspaces
func (w *Workspaces) UnmarshalJSON(data []byte) error {
	// Try array format first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		w.Patterns = arr
		return nil
	}

	// Try object format
	var obj struct {
		Packages []string `json:"packages"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		w.Patterns = obj.Packages
		return nil
	}

	return nil
}

// ReadPackageJSON reads and parses a package.json file
func ReadPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

// FindPackageJSON searches for package.json starting from the given directory
// and walking up the directory tree
func FindPackageJSON(startDir string) (string, error) {
	dir := startDir
	for {
		pkgPath := filepath.Join(dir, "package.json")
		if _, err := os.Stat(pkgPath); err == nil {
			return pkgPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// HasWorkspaces returns true if the package.json defines workspaces
func (p *PackageJSON) HasWorkspaces() bool {
	return len(p.Workspaces.Patterns) > 0
}
