package workspace

import (
	"fmt"
	"path/filepath"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// FuzzySelectPackage shows an interactive fuzzy finder for selecting a workspace package
func FuzzySelectPackage(rootDir string) (*Package, error) {
	packages, err := FindPackages(rootDir)
	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("no packages found in workspace")
	}

	idx, err := fuzzyfinder.Find(
		packages,
		func(i int) string {
			return packages[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			pkg := packages[i]
			relPath, _ := filepath.Rel(rootDir, pkg.Dir)
			return fmt.Sprintf("Name: %s\nPath: %s", pkg.Name, relPath)
		}),
	)

	if err != nil {
		return nil, err
	}

	return &packages[idx], nil
}

// FuzzySelectPackages shows an interactive fuzzy finder for selecting multiple packages
func FuzzySelectPackages(rootDir string) ([]Package, error) {
	packages, err := FindPackages(rootDir)
	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, fmt.Errorf("no packages found in workspace")
	}

	indices, err := fuzzyfinder.FindMulti(
		packages,
		func(i int) string {
			return packages[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			pkg := packages[i]
			relPath, _ := filepath.Rel(rootDir, pkg.Dir)
			return fmt.Sprintf("Name: %s\nPath: %s", pkg.Name, relPath)
		}),
	)

	if err != nil {
		return nil, err
	}

	result := make([]Package, len(indices))
	for i, idx := range indices {
		result[i] = packages[idx]
	}

	return result, nil
}
