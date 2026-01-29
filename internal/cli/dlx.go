package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/pmcombo"
	"github.com/user/fnpm/internal/runner"
)

var dlxCmd = &cobra.Command{
	Use:   "dlx <package> [args...]",
	Short: "Download and execute a package",
	Long:  `Download a package temporarily and execute it (like npx).`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		pkg := args[0]
		pkgArgs := []string{}
		if len(args) > 1 {
			pkgArgs = args[1:]
		}

		dlxCommand := pmcombo.NewDlxCommand(pmcombo.DlxOptions{
			Package: pkg,
			Args:    pkgArgs,
		})

		cmdArgs, err := dlxCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
