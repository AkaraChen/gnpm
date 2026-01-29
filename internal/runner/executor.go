package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/user/fnpm/internal/logger"
	"github.com/user/fnpm/internal/pmcombo"
)

// Options for running commands
type Options struct {
	Verbose bool
	DryRun  bool
}

// Run executes a package manager command
func Run(pm pmcombo.PackageManager, args []string, workDir string, opts Options) error {
	executable := pm.Executable()
	cmdStr := executable
	if len(args) > 0 {
		cmdStr += " " + strings.Join(args, " ")
	}

	if opts.DryRun {
		logger.DryRun(cmdStr, workDir)
		return nil
	}

	if opts.Verbose {
		logger.Command(cmdStr)
	}

	cmd := exec.Command(executable, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// RunOutput executes a package manager command and returns the output
func RunOutput(pm pmcombo.PackageManager, args []string, workDir string) (string, error) {
	executable := pm.Executable()

	cmd := exec.Command(executable, args...)
	cmd.Dir = workDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
