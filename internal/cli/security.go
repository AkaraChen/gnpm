package cli

import "github.com/AkaraChen/gnpm/internal/security"

func runPackageManagerSecurityCheck() {
	security.RunPackageManagerSecurityCheck(ctx, security.Options{
		DryRun:  dryRun,
		Verbose: verbose,
	})
}
