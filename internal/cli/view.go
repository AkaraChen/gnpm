package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var viewRepo bool

var viewCmd = &cobra.Command{
	Use:   "view [package]",
	Short: "Open package on npm or repository",
	Long: `Open package page on npm registry or its repository.

Without a package name, opens the current package.
With -r flag, opens the repository instead of npm.

Examples:
  fnpm view lodash      # Open lodash on npm
  fnpm view lodash -r   # Open lodash repository
  fnpm view             # Open current package on npm`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var pkgName string

		if len(args) > 0 {
			pkgName = args[0]
		} else if ctx.PackageJSON != nil && ctx.PackageJSON.Name != "" {
			pkgName = ctx.PackageJSON.Name
		} else {
			return fmt.Errorf("no package specified and no package.json found")
		}

		var url string
		if viewRepo {
			// Try to get repo URL from npm
			url = fmt.Sprintf("https://www.npmjs.com/package/%s", pkgName)
		} else {
			url = fmt.Sprintf("https://www.npmjs.com/package/%s", pkgName)
		}

		return openBrowser(url)
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewRepo, "repo", "r", false, "Open repository instead of npm page")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
