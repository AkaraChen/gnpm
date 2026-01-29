package cli

import (
	"github.com/spf13/cobra"

	"github.com/user/fnpm/internal/logger"
	"github.com/user/fnpm/internal/native"
)

var registryGlobal bool

var registryCmd = &cobra.Command{
	Use:   "registry [url|preset]",
	Short: "Get or set the npm registry URL",
	Long: `Get or set the npm registry URL.

Without arguments, shows the current registry.
With a URL argument, sets the registry to that URL.
You can also use preset names: npm, yarn, taobao, cnpm, tencent

Examples:
  fnpm registry                    # Show current registry
  fnpm registry https://registry.npmmirror.com  # Set registry
  fnpm registry taobao             # Use taobao preset
  fnpm registry npm                # Reset to default npm registry
  fnpm registry -g taobao          # Set global registry`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir, err := getWorkingDir()
		if err != nil {
			return err
		}

		opts := native.RegistryOptions{
			Dir:    workDir,
			Global: registryGlobal,
			DryRun: dryRun,
		}

		if len(args) == 0 {
			registry, err := native.GetRegistry(opts)
			if err != nil {
				return err
			}
			logger.Plainln(registry)
			return nil
		}

		url := args[0]
		if preset, ok := native.GetPreset(url); ok {
			url = preset
		}

		return native.SetRegistry(opts, url)
	},
}

func init() {
	registryCmd.Flags().BoolVarP(&registryGlobal, "global", "g", false, "Use global .npmrc")
}
