package workspace

import (
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/user/fnpm/internal/context"
)

// Package represents a workspace package
type Package struct {
	Name string
	Path string
	Dir  string
}

// FindPackages finds all packages in a workspace
func FindPackages(rootDir string) ([]Package, error) {
	patterns, err := getWorkspacePatterns(rootDir)
	if err != nil {
		return nil, err
	}

	var packages []Package
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		// Handle glob patterns
		fullPattern := filepath.Join(rootDir, pattern)

		// If pattern doesn't contain wildcards, treat it as a direct path
		if !containsGlob(pattern) {
			fullPattern = filepath.Join(rootDir, pattern, "package.json")
		} else {
			// Append package.json to glob pattern
			fullPattern = filepath.Join(fullPattern, "package.json")
		}

		matches, err := filepath.Glob(fullPattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			dir := filepath.Dir(match)
			if seen[dir] {
				continue
			}
			seen[dir] = true

			pkg, err := context.ReadPackageJSON(match)
			if err != nil {
				continue
			}

			packages = append(packages, Package{
				Name: pkg.Name,
				Path: match,
				Dir:  dir,
			})
		}
	}

	// Sort by name for consistent ordering
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages, nil
}

// getWorkspacePatterns extracts workspace patterns from package.json or pnpm-workspace.yaml
func getWorkspacePatterns(rootDir string) ([]string, error) {
	// Check pnpm-workspace.yaml first
	pnpmWorkspacePath := filepath.Join(rootDir, "pnpm-workspace.yaml")
	if patterns, err := readPnpmWorkspace(pnpmWorkspacePath); err == nil && len(patterns) > 0 {
		return patterns, nil
	}

	// Fall back to package.json workspaces
	pkgPath := filepath.Join(rootDir, "package.json")
	pkg, err := context.ReadPackageJSON(pkgPath)
	if err != nil {
		return nil, err
	}

	return pkg.Workspaces.Patterns, nil
}

// PnpmWorkspace represents pnpm-workspace.yaml structure
type PnpmWorkspace struct {
	Packages []string `yaml:"packages"`
}

// readPnpmWorkspace reads workspace patterns from pnpm-workspace.yaml
func readPnpmWorkspace(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var workspace PnpmWorkspace
	if err := yaml.Unmarshal(data, &workspace); err != nil {
		return nil, err
	}

	return workspace.Packages, nil
}

// containsGlob checks if a pattern contains glob characters
func containsGlob(pattern string) bool {
	for _, c := range pattern {
		if c == '*' || c == '?' || c == '[' {
			return true
		}
	}
	return false
}

// FindPackageByName finds a package by name in the workspace
func FindPackageByName(rootDir string, name string) (*Package, error) {
	packages, err := FindPackages(rootDir)
	if err != nil {
		return nil, err
	}

	for _, pkg := range packages {
		if pkg.Name == name {
			return &pkg, nil
		}
	}

	return nil, os.ErrNotExist
}
