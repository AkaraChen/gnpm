package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold [args...]",
	Short: "Scaffold a project using create-akrc",
	Long:  `Scaffold a new project using the create-akrc template.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		// scaffold is just dlx create-akrc
		dlxCommand := pmcombo.NewDlxCommand(pmcombo.DlxOptions{
			Package: "create-akrc",
			Args:    args,
		})

		cmdArgs, err := dlxCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
