package fixture_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/fnpm/internal/context"
	"github.com/user/fnpm/internal/workspace"
)

// fixtureDir returns the absolute path to a fixture directory
func fixtureDir(t *testing.T, name string) string {
	t.Helper()
	absPath, err := filepath.Abs(name)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	return absPath
}

// buildFnpm compiles fnpm and returns the binary path
func buildFnpm(t *testing.T) string {
	t.Helper()

	binPath := filepath.Join(t.TempDir(), "fnpm")
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/fnpm")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build fnpm: %v\n%s", err, output)
	}
	return binPath
}

// =============================================================================
// E2E Tests: Package Manager Detection
// =============================================================================

func TestE2E_PMDetection(t *testing.T) {
	cases := []struct {
		dir      string
		expected string // expected executable name (yarn-classic → yarn)
	}{
		{"npm-single", "npm"},
		{"npm-mono-5", "npm"},
		{"npm-mono-100", "npm"},
		{"yarn-classic-single", "yarn"},
		{"yarn-classic-mono-5", "yarn"},
		{"yarn-classic-mono-100", "yarn"},
		{"yarn-berry-single", "yarn"},
		{"yarn-berry-mono-5", "yarn"},
		{"yarn-berry-mono-100", "yarn"},
		{"pnpm-single", "pnpm"},
		{"pnpm-mono-5", "pnpm"},
		{"pnpm-mono-100", "pnpm"},
		{"deno-single", "deno"},
		{"bun-single", "bun"},
	}

	fnpm := buildFnpm(t)

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			cmd := exec.Command(fnpm, "install", "--dry-run")
			cmd.Dir = fixtureDir(t, tc.dir)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("fnpm failed: %v\nstderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("failed to run fnpm: %v", err)
			}

			// dry-run outputs command to stdout
			got := strings.Fields(strings.TrimSpace(string(output)))[0]
			if got != tc.expected {
				t.Errorf("expected PM %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// E2E Tests: Install Command
// =============================================================================

func TestE2E_Install(t *testing.T) {
	cases := []struct {
		dir      string
		expected string
	}{
		{"npm-single", "npm install"},
		{"yarn-classic-single", "yarn install"},
		{"yarn-berry-single", "yarn install"},
		{"pnpm-single", "pnpm install"},
		{"deno-single", "deno install"},
		{"bun-single", "bun install"},
	}

	fnpm := buildFnpm(t)

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			cmd := exec.Command(fnpm, "install", "--dry-run")
			cmd.Dir = fixtureDir(t, tc.dir)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("fnpm failed: %v\nstderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("failed to run fnpm: %v", err)
			}

			got := strings.TrimSpace(string(output))
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// E2E Tests: Add Command
// =============================================================================

func TestE2E_Add(t *testing.T) {
	cases := []struct {
		dir      string
		args     []string
		expected string
	}{
		// Basic add
		{"npm-single", []string{"lodash"}, "npm install lodash"},
		{"yarn-classic-single", []string{"lodash"}, "yarn add lodash"},
		{"yarn-berry-single", []string{"lodash"}, "yarn add lodash"},
		{"pnpm-single", []string{"lodash"}, "pnpm add lodash"},
		{"bun-single", []string{"lodash"}, "bun add lodash"},
		{"deno-single", []string{"lodash"}, "deno add npm:lodash"},

		// Dev dependency
		{"npm-single", []string{"-D", "typescript"}, "npm install --save-dev typescript"},
		{"yarn-classic-single", []string{"-D", "typescript"}, "yarn add -D typescript"},
		{"yarn-berry-single", []string{"-D", "typescript"}, "yarn add -D typescript"},
		{"pnpm-single", []string{"-D", "typescript"}, "pnpm add -D typescript"},
		{"bun-single", []string{"-D", "typescript"}, "bun add -d typescript"},

		// Multiple packages
		{"npm-single", []string{"lodash", "express"}, "npm install lodash express"},
		{"pnpm-single", []string{"lodash", "express"}, "pnpm add lodash express"},
	}

	fnpm := buildFnpm(t)

	for _, tc := range cases {
		name := tc.dir + "_" + strings.Join(tc.args, "_")
		t.Run(name, func(t *testing.T) {
			args := append([]string{"add", "--dry-run"}, tc.args...)
			cmd := exec.Command(fnpm, args...)
			cmd.Dir = fixtureDir(t, tc.dir)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("fnpm failed: %v\nstderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("failed to run fnpm: %v", err)
			}

			got := strings.TrimSpace(string(output))
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// E2E Tests: Remove Command
// =============================================================================

func TestE2E_Remove(t *testing.T) {
	cases := []struct {
		dir      string
		args     []string
		expected string
	}{
		{"npm-single", []string{"lodash"}, "npm uninstall lodash"},
		{"yarn-classic-single", []string{"lodash"}, "yarn remove lodash"},
		{"yarn-berry-single", []string{"lodash"}, "yarn remove lodash"},
		{"pnpm-single", []string{"lodash"}, "pnpm remove lodash"},
		{"bun-single", []string{"lodash"}, "bun remove lodash"},
		{"deno-single", []string{"lodash"}, "deno remove lodash"},
	}

	fnpm := buildFnpm(t)

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			args := append([]string{"remove", "--dry-run"}, tc.args...)
			cmd := exec.Command(fnpm, args...)
			cmd.Dir = fixtureDir(t, tc.dir)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("fnpm failed: %v\nstderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("failed to run fnpm: %v", err)
			}

			got := strings.TrimSpace(string(output))
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// E2E Tests: Run Command
// =============================================================================

// Note: fnpm run uses native script execution (directly runs the script content
// from package.json) rather than delegating to the package manager. The dry-run
// output shows the script content, not "pm run <script>".
func TestE2E_Run(t *testing.T) {
	cases := []struct {
		dir      string
		script   string
		expected string // Expected: the script content from package.json
	}{
		{"npm-single", "build", "echo build"},
		{"yarn-classic-single", "build", "echo build"},
		{"yarn-berry-single", "build", "echo build"},
		{"pnpm-single", "build", "echo build"},
		{"bun-single", "build", "echo build"},
		{"deno-single", "build", "echo build"},
	}

	fnpm := buildFnpm(t)

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			cmd := exec.Command(fnpm, "run", "--dry-run", tc.script)
			cmd.Dir = fixtureDir(t, tc.dir)

			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("fnpm failed: %v\nstderr: %s", err, exitErr.Stderr)
				}
				t.Fatalf("failed to run fnpm: %v", err)
			}

			got := strings.TrimSpace(string(output))
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// E2E Tests: Workspace Detection
// =============================================================================

func TestE2E_WorkspaceDetection(t *testing.T) {
	cases := []struct {
		dir         string
		isWorkspace bool
		pkgCount    int
	}{
		// Single packages - not workspaces
		{"npm-single", false, 0},
		{"yarn-classic-single", false, 0},
		{"yarn-berry-single", false, 0},
		{"pnpm-single", false, 0},
		{"deno-single", false, 0},
		{"bun-single", false, 0},

		// Monorepos - are workspaces
		{"npm-mono-5", true, 5},
		{"npm-mono-100", true, 100},
		{"yarn-classic-mono-5", true, 5},
		{"yarn-classic-mono-100", true, 100},
		{"yarn-berry-mono-5", true, 5},
		{"yarn-berry-mono-100", true, 100},
		{"pnpm-mono-5", true, 5},
		{"pnpm-mono-100", true, 100},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			dir := fixtureDir(t, tc.dir)
			ctx, err := context.Detect(dir)
			if err != nil {
				t.Fatalf("failed to detect context: %v", err)
			}

			if ctx.IsWorkspace != tc.isWorkspace {
				t.Errorf("expected IsWorkspace=%v, got %v", tc.isWorkspace, ctx.IsWorkspace)
			}

			if tc.isWorkspace {
				packages, err := workspace.FindPackages(dir)
				if err != nil {
					t.Fatalf("failed to find packages: %v", err)
				}
				if len(packages) != tc.pkgCount {
					t.Errorf("expected %d packages, got %d", tc.pkgCount, len(packages))
				}
			}
		})
	}
}

// =============================================================================
// E2E Tests: Internal PM Type Detection (yarn classic vs berry)
// =============================================================================

func TestE2E_YarnVersionDetection(t *testing.T) {
	cases := []struct {
		dir      string
		expected string // internal PM type
	}{
		{"yarn-classic-single", "yarn-classic"},
		{"yarn-classic-mono-5", "yarn-classic"},
		{"yarn-berry-single", "yarn"},
		{"yarn-berry-mono-5", "yarn"},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			dir := fixtureDir(t, tc.dir)
			ctx, err := context.Detect(dir)
			if err != nil {
				t.Fatalf("failed to detect context: %v", err)
			}

			got := string(ctx.PackageManager)
			if got != tc.expected {
				t.Errorf("expected PM type %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// Benchmarks: Context Detection
// =============================================================================

func BenchmarkDetect_NPM_Single(b *testing.B) {
	dir, _ := filepath.Abs("npm-single")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_NPM_Mono5(b *testing.B) {
	dir, _ := filepath.Abs("npm-mono-5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_NPM_Mono100(b *testing.B) {
	dir, _ := filepath.Abs("npm-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_YarnBerry_Single(b *testing.B) {
	dir, _ := filepath.Abs("yarn-berry-single")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_YarnBerry_Mono100(b *testing.B) {
	dir, _ := filepath.Abs("yarn-berry-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_PNPM_Single(b *testing.B) {
	dir, _ := filepath.Abs("pnpm-single")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

func BenchmarkDetect_PNPM_Mono100(b *testing.B) {
	dir, _ := filepath.Abs("pnpm-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = context.Detect(dir)
	}
}

// =============================================================================
// Benchmarks: Workspace Package Finding
// =============================================================================

func BenchmarkFindPackages_NPM_5(b *testing.B) {
	dir, _ := filepath.Abs("npm-mono-5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}

func BenchmarkFindPackages_NPM_100(b *testing.B) {
	dir, _ := filepath.Abs("npm-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}

func BenchmarkFindPackages_YarnBerry_5(b *testing.B) {
	dir, _ := filepath.Abs("yarn-berry-mono-5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}

func BenchmarkFindPackages_YarnBerry_100(b *testing.B) {
	dir, _ := filepath.Abs("yarn-berry-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}

func BenchmarkFindPackages_PNPM_5(b *testing.B) {
	dir, _ := filepath.Abs("pnpm-mono-5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}

func BenchmarkFindPackages_PNPM_100(b *testing.B) {
	dir, _ := filepath.Abs("pnpm-mono-100")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = workspace.FindPackages(dir)
	}
}
