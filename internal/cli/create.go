package cli

import (
	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/native"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"github.com/AkaraChen/gnpm/internal/runner"
)

var createYes bool

var createCmd = &cobra.Command{
	Use:     "create [template] [args...]",
	Aliases: []string{"c", "init"},
	Short:   "Create a new project or initialize package.json",
	Long: `Create a new project from a template or initialize package.json.

Without arguments, creates a new package.json (like npm init).
With a template name, creates a new project from that template.

Examples:
  gnpm create               # Create package.json (init)
  gnpm init                 # Same as above
  gnpm create -y            # Create package.json with defaults
  gnpm init -y              # Same as above
  gnpm create react-app     # Create project with create-react-app
  gnpm init react-app       # Same as above
  gnpm c react-app          # Same as above`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		// If no template specified, initialize package.json
		if len(args) == 0 {
			return native.Init(native.InitOptions{
				Dir:    workDir,
				Yes:    createYes,
				DryRun: dryRun,
			})
		}

		// Otherwise, create from template
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

func init() {
	createCmd.Flags().BoolVarP(&createYes, "yes", "y", false, "Skip prompts (for init)")
}
