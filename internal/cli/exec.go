package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/native"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var execCmd = &cobra.Command{
	Use:     "exec <command> [args...]",
	Aliases: []string{"x", "npx", "dlx"},
	Short:   "Execute a binary from node_modules or download and run",
	Long: `Execute a binary from node_modules/.bin, or download and execute if not found.

First searches for the binary in node_modules/.bin (walking up the directory tree).
If not found locally, falls back to downloading and executing (like npx/dlx).

Examples:
  gnpm exec eslint .     # Run local eslint
  gnpm x eslint .        # Same as above
  gnpm x cowsay hello    # Download and run cowsay if not local
  gnpm npx cowsay hello  # Same as above
  gnpm dlx cowsay hello  # Same as above`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		command := args[0]
		commandArgs := []string{}
		if len(args) > 1 {
			commandArgs = args[1:]
		}

		// Try to find the binary locally first
		_, err = native.FindBinary(workDir, command)
		if err == nil {
			// Binary found locally, execute it
			return native.Exec(native.ExecOptions{
				Dir:     workDir,
				Command: command,
				Args:    commandArgs,
				Verbose: verbose,
				DryRun:  dryRun,
			})
		}

		// Binary not found locally, fall back to dlx
		dlxCommand := pmcombo.NewDlxCommand(pmcombo.DlxOptions{
			Package: command,
			Args:    commandArgs,
		})

		cmdArgs, err := dlxCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
