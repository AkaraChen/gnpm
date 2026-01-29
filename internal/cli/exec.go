package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/native"
)

var execCmd = &cobra.Command{
	Use:   "exec <command> [args...]",
	Short: "Execute a binary from node_modules",
	Long:  `Execute a binary from node_modules/.bin.`,
	Args:  cobra.MinimumNArgs(1),
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

		return native.Exec(native.ExecOptions{
			Dir:     workDir,
			Command: command,
			Args:    commandArgs,
			Verbose: verbose,
			DryRun:  dryRun,
		})
	},
}
