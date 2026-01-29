package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var updateInteractive bool
var updateLatest bool

var updateCmd = &cobra.Command{
	Use:     "update [packages...]",
	Aliases: []string{"up", "upgrade"},
	Short:   "Update packages",
	Long:    `Update one or more packages to their latest versions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		updateCommand := pmcombo.NewUpdateCommand(pmcombo.UpdateOptions{
			Packages:    args,
			Interactive: updateInteractive,
			Latest:      updateLatest,
		})

		cmdArgs, err := updateCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateInteractive, "interactive", "i", false, "Interactive mode")
	updateCmd.Flags().BoolVarP(&updateLatest, "latest", "L", false, "Update to latest version (ignore semver)")
}
