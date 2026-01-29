package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <pm>[@version]",
	Short: "Switch package manager version using corepack",
	Long: `Switch package manager version using corepack.

Examples:
  fnpm use pnpm@8
  fnpm use yarn@4
  fnpm use npm@10`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pmSpec := args[0]

		// Use corepack to enable the package manager
		corepackCmd := exec.Command("corepack", "use", pmSpec)
		corepackCmd.Stdout = os.Stdout
		corepackCmd.Stderr = os.Stderr
		corepackCmd.Stdin = os.Stdin

		if err := corepackCmd.Run(); err != nil {
			return fmt.Errorf("corepack use failed: %w", err)
		}

		return nil
	},
}
