package native

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/user/fnpm/internal/context"
	"github.com/user/fnpm/internal/logger"
)

// RunOptions for running scripts
type RunOptions struct {
	Dir     string
	Script  string
	Args    []string
	Verbose bool
	DryRun  bool
}

// Run executes a script from package.json
func Run(opts RunOptions) error {
	// Find package.json by walking up
	pkgPath, err := context.FindPackageJSON(opts.Dir)
	if err != nil {
		return err
	}

	pkgDir := filepath.Dir(pkgPath)
	pkg, err := context.ReadPackageJSON(pkgPath)
	if err != nil {
		return err
	}

	// Find script
	scriptCmd, ok := pkg.Scripts[opts.Script]
	if !ok {
		if len(pkg.Scripts) > 0 {
			logger.Header("Available scripts:")
			for name := range pkg.Scripts {
				logger.List(name)
			}
		}
		logger.Error("script %q not found", opts.Script)
		return nil
	}

	// Append extra args
	if len(opts.Args) > 0 {
		scriptCmd = scriptCmd + " " + strings.Join(opts.Args, " ")
	}

	if opts.DryRun {
		logger.DryRun(scriptCmd, pkgDir)
		return nil
	}

	if opts.Verbose {
		logger.Command(scriptCmd)
	}

	return executeScript(scriptCmd, pkgDir)
}

// executeScript runs a shell command with node_modules/.bin in PATH
func executeScript(script string, dir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", script)
	} else {
		cmd = exec.Command("sh", "-c", script)
	}

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), "PATH="+BuildNodeBinPath(dir))

	return cmd.Run()
}

// ListScripts returns all available scripts
func ListScripts(dir string) (map[string]string, error) {
	pkgPath, err := context.FindPackageJSON(dir)
	if err != nil {
		return nil, err
	}
	pkg, err := context.ReadPackageJSON(pkgPath)
	if err != nil {
		return nil, err
	}
	return pkg.Scripts, nil
}
