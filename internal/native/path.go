package native

import (
	"os"
	"path/filepath"
	"strings"
)

// BuildNodeBinPath builds PATH with all node_modules/.bin directories up the tree
func BuildNodeBinPath(dir string) string {
	var binDirs []string
	current := dir

	for {
		binDir := filepath.Join(current, "node_modules", ".bin")
		if _, err := os.Stat(binDir); err == nil {
			binDirs = append(binDirs, binDir)
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	// Append original PATH
	pathEnv := os.Getenv("PATH")
	if len(binDirs) > 0 {
		return strings.Join(binDirs, string(os.PathListSeparator)) + string(os.PathListSeparator) + pathEnv
	}
	return pathEnv
}
