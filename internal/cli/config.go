package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/native"
)

var configGlobal bool

var configCmd = &cobra.Command{
	Use:   "config <action> [key] [value]",
	Short: "Manage .npmrc configuration",
	Long: `Manage .npmrc configuration natively.

Actions:
  get <key>         Get a config value
  set <key> <value> Set a config value
  delete <key>      Delete a config value
  list              List all config values

By default, operates on project-level .npmrc.
Use -g flag for user-level ~/.npmrc.

Examples:
  fnpm config get registry
  fnpm config set registry https://registry.npmmirror.com
  fnpm config set @myorg:registry https://npm.myorg.com
  fnpm config delete registry
  fnpm config list
  fnpm config set registry https://registry.npmjs.org -g`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		action := args[0]
		key := ""
		value := ""

		if len(args) > 1 {
			key = args[1]
		}
		if len(args) > 2 {
			value = args[2]
		}

		return native.Config(native.ConfigOptions{
			Dir:    workDir,
			Action: action,
			Key:    key,
			Value:  value,
			Global: configGlobal,
			DryRun: dryRun,
		})
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configGlobal, "global", "g", false, "Operate on user-level ~/.npmrc")
}
