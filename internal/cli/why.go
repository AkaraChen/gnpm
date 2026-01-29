package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var whyCmd = &cobra.Command{
	Use:   "why <package>",
	Short: "Show why a package is installed",
	Long:  `Show why a package is installed and what depends on it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		whyCommand := pmcombo.NewWhyCommand(pmcombo.WhyOptions{
			Package: args[0],
		})

		cmdArgs, err := whyCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
