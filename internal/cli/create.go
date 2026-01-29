package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var createCmd = &cobra.Command{
	Use:   "create <template> [args...]",
	Short: "Create a new project from a template",
	Long:  `Create a new project from a create-* template.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		template := args[0]
		templateArgs := []string{}
		if len(args) > 1 {
			templateArgs = args[1:]
		}

		createCommand := pmcombo.NewCreateCommand(pmcombo.CreateOptions{
			Template: template,
			Args:     templateArgs,
		})

		cmdArgs, err := createCommand.Concat(ctx.PackageManager)
		if err != nil {
			return err
		}

		return runner.Run(ctx.PackageManager, cmdArgs, workDir, runnerOpts())
	},
}
