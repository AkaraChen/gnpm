package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/native"
)

var runCmd = &cobra.Command{
	Use:     "run <script> [args...]",
	Aliases: []string{"r"},
	Short:   "Run a script from package.json",
	Long:  `Run a script defined in package.json.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		script := args[0]
		scriptArgs := []string{}
		if len(args) > 1 {
			scriptArgs = args[1:]
		}

		return native.Run(native.RunOptions{
			Dir:     workDir,
			Script:  script,
			Args:    scriptArgs,
			Verbose: verbose,
			DryRun:  dryRun,
		})
	},
}
