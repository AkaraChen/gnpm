package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var addDev bool
var addExact bool
var addGlobal bool
var addPeer bool
var addOptional bool

var addCmd = &cobra.Command{
	Use:   "add [packages...]",
	Short: "Add packages to the project",
	Long:  `Add one or more packages to the project dependencies.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		addCmd := pmcombo.NewAddCommand(pmcombo.AddOptions{
			Packages: args,
			Dev:      addDev,
			Exact:    addExact,
			Global:   addGlobal,
			Peer:     addPeer,
			Optional: addOptional,
		})

		cmdArgs, err := addCmd.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addDev, "dev", "D", false, "Add as dev dependency")
	addCmd.Flags().BoolVarP(&addExact, "exact", "E", false, "Add exact version")
	addCmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "Add globally")
	addCmd.Flags().BoolVar(&addPeer, "peer", false, "Add as peer dependency")
	addCmd.Flags().BoolVarP(&addOptional, "optional", "O", false, "Add as optional dependency")
}
