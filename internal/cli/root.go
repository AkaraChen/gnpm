package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/logger"
	"github.com/AkaraChen/gnpm/internal/native"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
	"github.com/AkaraChen/gnpm/internal/workspace"
)

var (
	// Global flags
	workspaceRoot bool
	fuzzySelect   bool
	verbose       bool
	dryRun        bool
	usePM         string

	// Detected context
	ctx *context.ProjectContext
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gnpm",
	Short: "A fast, unified package manager CLI",
	Long: `gnpm is a fast, unified CLI that wraps npm, yarn, pnpm, deno, and bun.
It detects your package manager from lock files and translates commands automatically.

Unknown commands are automatically resolved:
  1. As a script in package.json
  2. As a binary in node_modules/.bin
  3. As a system command`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip context detection for commands that don't need it
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// If -w flag is set, find workspace root
		if workspaceRoot {
			wsRoot, err := context.FindWorkspaceRoot(cwd)
			if err != nil {
				return fmt.Errorf("not in a workspace: %w", err)
			}
			cwd = wsRoot
		}

		// Detect project context
		ctx, err = context.Detect(cwd)
		if err != nil {
			// Not in a project, some commands might still work
			ctx = &context.ProjectContext{
				RootDir:        cwd,
				PackageManager: pmcombo.NPM,
			}
		}

		// Override PM if specified
		if usePM != "" {
			pm, err := pmcombo.ParsePackageManager(usePM)
			if err != nil {
				return err
			}
			ctx.PackageManager = pm
		}

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	// Check if the first argument is an unknown command (not a flag)
	if len(os.Args) > 1 && !isFlag(os.Args[1]) {
		cmd, _, err := rootCmd.Find(os.Args[1:])
		if err != nil || cmd == rootCmd {
			// Unknown command, try fallback
			return runFallback(os.Args[1], os.Args[2:])
		}
	}
	return rootCmd.Execute()
}

// isFlag checks if arg starts with - or --
func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

// runFallback handles unknown commands using the fallback mechanism
func runFallback(command string, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	result, err := native.Fallback(native.FallbackOptions{
		Dir:     cwd,
		Command: command,
		Args:    args,
		Verbose: verbose,
		DryRun:  dryRun,
	})

	if result == native.FallbackNotFound {
		logger.Error("unknown command %q", command)
		logger.Plainln("\nRun 'gnpm --help' for usage.")
		return nil
	}

	return err
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&workspaceRoot, "workspace", "w", false, "Run command in workspace root")
	rootCmd.PersistentFlags().BoolVarP(&fuzzySelect, "select", "s", false, "Fuzzy select a workspace package")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print command without executing")
	rootCmd.PersistentFlags().StringVar(&usePM, "pm", "", "Override detected package manager (npm, yarn, pnpm, bun, deno)")

	// Add all subcommands
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(whyCmd)
	rootCmd.AddCommand(ciCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(registryCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(scaffoldCmd)
}

// getWorkingDir returns the working directory for the command
// If -s flag is set, shows fuzzy finder to select a package
func getWorkingDir() (string, error) {
	if fuzzySelect && ctx.IsWorkspace {
		pkg, err := workspace.FuzzySelectPackage(ctx.RootDir)
		if err != nil {
			return "", err
		}
		return pkg.Dir, nil
	}
	return ctx.RootDir, nil
}

// runnerOpts returns the runner options from global flags
func runnerOpts() runner.Options {
	return runner.Options{
		Verbose: verbose,
		DryRun:  dryRun,
	}
}
