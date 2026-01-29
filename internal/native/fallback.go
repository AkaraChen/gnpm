package native

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/logger"
)

// FallbackOptions for running unknown commands
type FallbackOptions struct {
	Dir     string
	Command string
	Args    []string
	Verbose bool
	DryRun  bool
}

// FallbackResult indicates what type of command was executed
type FallbackResult int

const (
	FallbackScript FallbackResult = iota
	FallbackBinary
	FallbackSystem
	FallbackNotFound
)

// Fallback tries to run an unknown command by checking scripts, binaries, and system commands
func Fallback(opts FallbackOptions) (FallbackResult, error) {
	// 1. Check if it's a script in package.json
	if hasScript(opts.Dir, opts.Command) {
		err := Run(RunOptions{
			Dir:     opts.Dir,
			Script:  opts.Command,
			Args:    opts.Args,
			Verbose: opts.Verbose,
			DryRun:  opts.DryRun,
		})
		return FallbackScript, err
	}

	// 2. Check if it's a binary in node_modules/.bin
	if binPath, err := FindBinary(opts.Dir, opts.Command); err == nil {
		err := Exec(ExecOptions{
			Dir:     opts.Dir,
			Command: opts.Command,
			Args:    opts.Args,
			Verbose: opts.Verbose,
			DryRun:  opts.DryRun,
		})
		_ = binPath // unused but checked for existence
		return FallbackBinary, err
	}

	// 3. Try to run as a system command
	if hasSystemCommand(opts.Command) {
		err := runSystemCommand(opts)
		return FallbackSystem, err
	}

	return FallbackNotFound, nil
}

// hasScript checks if a script exists in package.json
func hasScript(dir string, name string) bool {
	pkgPath, err := context.FindPackageJSON(dir)
	if err != nil {
		return false
	}
	pkg, err := context.ReadPackageJSON(pkgPath)
	if err != nil {
		return false
	}
	_, ok := pkg.Scripts[name]
	return ok
}

// hasSystemCommand checks if a command exists in PATH
func hasSystemCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// runSystemCommand executes a system command
func runSystemCommand(opts FallbackOptions) error {
	cmdStr := opts.Command
	for _, arg := range opts.Args {
		cmdStr += " " + arg
	}

	if opts.DryRun {
		logger.DryRun(cmdStr, opts.Dir)
		return nil
	}

	if opts.Verbose {
		logger.Command(cmdStr)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		args := append([]string{"/C", opts.Command}, opts.Args...)
		cmd = exec.Command("cmd", args...)
	} else {
		cmd = exec.Command(opts.Command, opts.Args...)
	}

	cmd.Dir = opts.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), "PATH="+BuildNodeBinPath(opts.Dir))

	return cmd.Run()
}
