package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Clean install dependencies (frozen lockfile)",
	Long:  `Install dependencies using frozen lockfile, suitable for CI environments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		installCmd := pmcombo.NewInstallCommand(pmcombo.InstallOptions{
			Frozen: true,
		})

		cmdArgs, err := installCmd.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
