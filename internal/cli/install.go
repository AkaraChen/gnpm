package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "Install all dependencies",
	Long:    `Install all dependencies from package.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		installCmd := pmcombo.NewInstallCommand(pmcombo.InstallOptions{
			Frozen: false,
		})

		cmdArgs, err := installCmd.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
