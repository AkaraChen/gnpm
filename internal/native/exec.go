package native

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AkaraChen/gnpm/internal/logger"
)

// ExecOptions for executing binaries
type ExecOptions struct {
	Dir     string
	Command string
	Args    []string
	Verbose bool
	DryRun  bool
}

// Exec executes a binary from node_modules/.bin
func Exec(opts ExecOptions) error {
	binPath, err := FindBinary(opts.Dir, opts.Command)
	if err != nil {
		return err
	}

	cmdStr := binPath
	if len(opts.Args) > 0 {
		cmdStr += " " + strings.Join(opts.Args, " ")
	}

	if opts.DryRun {
		logger.DryRun(cmdStr, opts.Dir)
		return nil
	}

	if opts.Verbose {
		logger.Command(cmdStr)
	}

	cmd := exec.Command(binPath, opts.Args...)
	cmd.Dir = opts.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), "PATH="+BuildNodeBinPath(opts.Dir))

	return cmd.Run()
}

// FindBinary searches for a binary in node_modules/.bin, walking up the directory tree
func FindBinary(dir string, name string) (string, error) {
	var candidates []string
	if runtime.GOOS == "windows" {
		candidates = []string{name + ".cmd", name + ".ps1", name + ".exe", name}
	} else {
		candidates = []string{name}
	}

	current := dir
	for {
		binDir := filepath.Join(current, "node_modules", ".bin")

		if _, err := os.Stat(binDir); err == nil {
			for _, candidate := range candidates {
				binPath := filepath.Join(binDir, candidate)
				if info, err := os.Stat(binPath); err == nil {
					if runtime.GOOS != "windows" && info.Mode()&0111 == 0 {
						continue
					}
					return binPath, nil
				}
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	listAvailableBinaries(dir)
	return "", fmt.Errorf("binary %q not found", name)
}

func listAvailableBinaries(dir string) {
	current := dir
	for {
		binDir := filepath.Join(current, "node_modules", ".bin")
		entries, err := os.ReadDir(binDir)
		if err == nil && len(entries) > 0 {
			logger.Header("Available binaries:")
			seen := make(map[string]bool)
			for _, e := range entries {
				baseName := trimExtensions(e.Name())
				if !seen[baseName] {
					seen[baseName] = true
					logger.List(baseName)
				}
			}
			return
		}

		parent := filepath.Dir(current)
		if parent == current {
			return
		}
		current = parent
	}
}

func trimExtensions(name string) string {
	for _, ext := range []string{".cmd", ".ps1", ".exe"} {
		if strings.HasSuffix(name, ext) {
			return strings.TrimSuffix(name, ext)
		}
	}
	return name
}
