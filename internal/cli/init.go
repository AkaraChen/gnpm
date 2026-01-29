package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/native"
)

var initYes bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new package.json",
	Long:  `Initialize a new package.json in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		return native.Init(native.InitOptions{
			Dir:    workDir,
			Yes:    initYes,
			DryRun: dryRun,
		})
	},
}

func init() {
	initCmd.Flags().BoolVarP(&initYes, "yes", "y", false, "Skip prompts")
}
