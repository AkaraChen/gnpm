package native

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/logger"
)

const defaultRegistry = "https://registry.npmjs.org/"

// RegistryOptions for registry operations
type RegistryOptions struct {
	Dir    string
	Global bool
	DryRun bool
}

// GetRegistry returns the current registry URL
func GetRegistry(opts RegistryOptions) (string, error) {
	// Check .npmrc files walking up the directory tree, then global
	paths := getNpmrcPaths(opts.Dir, opts.Global)

	for _, path := range paths {
		registry, err := readRegistryFromFile(path)
		if err == nil && registry != "" {
			return registry, nil
		}
	}

	return defaultRegistry, nil
}

// SetRegistry sets the registry URL
func SetRegistry(opts RegistryOptions, url string) error {
	var npmrcPath string

	if opts.Global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		npmrcPath = filepath.Join(home, ".npmrc")
	} else {
		// Find project root (where package.json is)
		projectRoot := opts.Dir
		if pkgPath, err := context.FindPackageJSON(opts.Dir); err == nil {
			projectRoot = filepath.Dir(pkgPath)
		}
		npmrcPath = filepath.Join(projectRoot, ".npmrc")
	}

	// Read existing content
	content, _ := os.ReadFile(npmrcPath)
	lines := strings.Split(string(content), "\n")

	// Update or add registry line
	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "registry=") {
			lines[i] = "registry=" + url
			found = true
			break
		}
	}

	if !found {
		// Add at the beginning
		lines = append([]string{"registry=" + url}, lines...)
	}

	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	output := strings.Join(lines, "\n") + "\n"

	if opts.DryRun {
		logger.DryRun("registry="+url, npmrcPath)
		return nil
	}

	if err := os.WriteFile(npmrcPath, []byte(output), 0644); err != nil {
		return err
	}

	logger.Success("registry set to %s", url)
	logger.Dim(npmrcPath)
	return nil
}

// getNpmrcPaths returns paths to check for .npmrc files (walks up directory tree)
func getNpmrcPaths(dir string, globalOnly bool) []string {
	var paths []string

	if !globalOnly {
		// Walk up directory tree looking for .npmrc files
		current := dir
		for {
			npmrcPath := filepath.Join(current, ".npmrc")
			paths = append(paths, npmrcPath)

			parent := filepath.Dir(current)
			if parent == current {
				break
			}
			current = parent
		}
	}

	// Global .npmrc (user home)
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, ".npmrc")
		// Avoid duplicate if we already walked to home
		if len(paths) == 0 || paths[len(paths)-1] != homePath {
			paths = append(paths, homePath)
		}
	}

	return paths
}

// readRegistryFromFile reads the registry value from an .npmrc file
func readRegistryFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "registry=") {
			return strings.TrimPrefix(line, "registry="), nil
		}
	}

	return "", nil
}

// Common registry presets
var RegistryPresets = map[string]string{
	"npm":       "https://registry.npmjs.org/",
	"yarn":      "https://registry.yarnpkg.com/",
	"taobao":    "https://registry.npmmirror.com/",
	"cnpm":      "https://r.cnpmjs.org/",
	"tencent":   "https://mirrors.cloud.tencent.com/npm/",
	"npmMirror": "https://skimdb.npmjs.com/registry/",
}

// GetPreset returns a preset registry URL by name
func GetPreset(name string) (string, bool) {
	url, ok := RegistryPresets[name]
	return url, ok
}
