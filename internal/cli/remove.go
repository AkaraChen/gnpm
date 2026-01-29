package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var removeGlobal bool

var removeCmd = &cobra.Command{
	Use:     "remove [packages...]",
	Aliases: []string{"rm", "uninstall", "un"},
	Short:   "Remove packages from the project",
	Long:    `Remove one or more packages from the project dependencies.`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		removeCmd := pmcombo.NewRemoveCommand(pmcombo.RemoveOptions{
			Packages: args,
			Global:   removeGlobal,
		})

		cmdArgs, err := removeCmd.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeGlobal, "global", "g", false, "Remove globally")
}
