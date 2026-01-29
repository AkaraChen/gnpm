package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AkaraChen/gnpm/internal/pmcombo"
)

var defaultCmd = &cobra.Command{
	Use:   "default [pm]",
	Short: "Get or set the default package manager",
	Long: `Get or set the default package manager for new projects.

Without arguments, shows the current detected package manager.
With a PM argument, sets it as default for the current project.

Examples:
  gnpm default        # Show current PM
  gnpm default pnpm   # Set pnpm as default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show current package manager
			fmt.Printf("Current package manager: %s\n", ctx.PackageManager)
			return nil
		}

		// Validate the package manager
		pm, err := pmcombo.ParsePackageManager(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Default package manager set to: %s\n", pm)
		fmt.Println("Note: This only affects the current session. Use 'gnpm use' to change the packageManager field in package.json.")
		return nil
	},
}
