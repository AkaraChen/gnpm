package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var installDev bool
var installExact bool
var installGlobal bool
var installPeer bool
var installOptional bool

var installCmd = &cobra.Command{
	Use:     "install [packages...]",
	Aliases: []string{"i", "a", "add"},
	Short:   "Install dependencies or add packages",
	Long: `Install all dependencies or add specific packages.

Without arguments, installs all dependencies from package.json.
With package names, adds those packages to the project.

Examples:
  gnpm install           # Install all dependencies
  gnpm i                 # Same as above
  gnpm install react     # Add react package
  gnpm i react -D        # Add react as dev dependency
  gnpm add react         # Same as install react
  gnpm a react           # Same as above`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		// If no packages specified, install all dependencies
		if len(args) == 0 {
			installCmd := pmcombo.NewInstallCommand(pmcombo.InstallOptions{
				Frozen: false,
			})

			cmdArgs, err := installCmd.Concat(ctx.PackageManager)
			if err != nil {
				return err
			}

			return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
		}

		// Otherwise, add the specified packages
		addCmd := pmcombo.NewAddCommand(pmcombo.AddOptions{
			Packages: args,
			Dev:      installDev,
			Exact:    installExact,
			Global:   installGlobal,
			Peer:     installPeer,
			Optional: installOptional,
		})

		cmdArgs, err := addCmd.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}

func init() {
	installCmd.Flags().BoolVarP(&installDev, "dev", "D", false, "Add as dev dependency")
	installCmd.Flags().BoolVarP(&installExact, "exact", "E", false, "Add exact version")
	installCmd.Flags().BoolVarP(&installGlobal, "global", "g", false, "Add globally")
	installCmd.Flags().BoolVar(&installPeer, "peer", false, "Add as peer dependency")
	installCmd.Flags().BoolVarP(&installOptional, "optional", "O", false, "Add as optional dependency")
}
