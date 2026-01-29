package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/native"
)

var testCmd = &cobra.Command{
	Use:     "test [args...]",
	Aliases: []string{"t"},
	Short:   "Run tests",
	Long:  `Run the test script from package.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		return native.Run(native.RunOptions{
			Dir:     workDir,
			Script:  "test",
			Args:    args,
			Verbose: verbose,
			DryRun:  dryRun,
		})
	},
}
