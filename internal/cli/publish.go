package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/runner"
)

var publishTag string
var publishAccess string
var publishDryRun bool

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the package to npm registry",
	Long:  `Publish the current package to the npm registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		// Publish is a direct passthrough to the package manager
		cmdArgs := []string{"publish"}

		if publishTag != "" {
			cmdArgs = append(cmdArgs, "--tag", publishTag)
		}
		if publishAccess != "" {
			cmdArgs = append(cmdArgs, "--access", publishAccess)
		}
		if publishDryRun {
			cmdArgs = append(cmdArgs, "--dry-run")
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}

func init() {
	publishCmd.Flags().StringVar(&publishTag, "tag", "", "Publish with a specific tag")
	publishCmd.Flags().StringVar(&publishAccess, "access", "", "Set access level (public/restricted)")
	publishCmd.Flags().BoolVar(&publishDryRun, "dry-run", false, "Run without actually publishing")
}
