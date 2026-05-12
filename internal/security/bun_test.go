package security

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	projectcontext "github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
)

func TestEnsureBunBestPracticesCreatesMissingSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")

	result, err := EnsureBunBestPractices(rootDir, bunMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureBunBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected missing settings to be written")
	}

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "[install]")
	assertContains(t, content, "exact = true")
	assertContains(t, content, "minimumReleaseAge = 259200")
	assertContains(t, content, "[install.lockfile]")
	assertContains(t, content, "save = true")
}

func TestEnsureBunBestPracticesFixesUnsafeSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lockb", "bun lockfile v0\x00")
	writeFile(t, rootDir, bunfigFile, `telemetry = false

[install]
registry = "https://registry.npmjs.org"
exact = false
minimumReleaseAge = 60

[install.lockfile]
save = false
`)

	result, err := EnsureBunBestPractices(rootDir, bunMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureBunBestPractices failed: %v", err)
	}

	expectedSettings := []string{
		"install.exact=true",
		"install.minimumReleaseAge=259200",
		"install.lockfile.save=true",
	}
	for _, setting := range expectedSettings {
		if !slices.Contains(result.Settings, setting) {
			t.Fatalf("expected %q in changed settings, got %#v", setting, result.Settings)
		}
	}

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "telemetry = false")
	assertContains(t, content, "registry = \"https://registry.npmjs.org\"")
	assertContains(t, content, "exact = true")
	assertContains(t, content, "minimumReleaseAge = 259200")
	assertContains(t, content, "save = true")
}

func TestEnsureBunBestPracticesDryRunDoesNotWrite(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")
	writeFile(t, rootDir, bunfigFile, "exact = false\n")

	result, err := EnsureBunBestPractices(rootDir, bunMinimumSafeVersion, Options{DryRun: true})
	if err != nil {
		t.Fatalf("EnsureBunBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected dry-run to report pending changes")
	}

	content, err := os.ReadFile(filepath.Join(rootDir, bunfigFile))
	if err != nil {
		t.Fatalf("read bun config: %v", err)
	}
	if string(content) != "exact = false\n" {
		t.Fatalf("dry-run changed file content:\n%s", content)
	}
}

func TestEnsureBunBestPracticesKeepsStrongerMinimumReleaseAge(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")
	writeFile(t, rootDir, bunfigFile, `[install]
exact = true
minimumReleaseAge = 604800

[install.lockfile]
save = true
`)

	result, err := EnsureBunBestPractices(rootDir, bunMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureBunBestPractices failed: %v", err)
	}
	if result.Changed {
		t.Fatalf("expected no changes, got %#v", result.Settings)
	}

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "minimumReleaseAge = 604800")
}

func TestBunMinimumSafeVersion(t *testing.T) {
	if got, want := bunMinimumSafeVersion, "1.3.0"; got != want {
		t.Fatalf("expected minimum safe Bun version %q, got %q", want, got)
	}
}

func TestCheckBunMinimumSafeVersionWarnsBelowMinimum(t *testing.T) {
	restore := stubBunVersionDetector(bunVersionResult{
		Version: "1.2.22",
		Source:  "bun",
	}, nil)
	defer restore()

	version, warning, err := checkBunMinimumSafeVersion(t.TempDir())
	if err != nil {
		t.Fatalf("checkBunMinimumSafeVersion failed: %v", err)
	}
	if version != "1.2.22" {
		t.Fatalf("expected detected version %q, got %q", "1.2.22", version)
	}
	if !strings.Contains(warning, "bun 1.2.22 from bun") {
		t.Fatalf("expected detected version in warning, got %q", warning)
	}
	if !strings.Contains(warning, "1.3.0") {
		t.Fatalf("expected minimum version in warning, got %q", warning)
	}
}

func TestCheckBunMinimumSafeVersionAllowsMinimumOrNewer(t *testing.T) {
	for _, version := range []string{"1.3.0", "1.3.1"} {
		t.Run(version, func(t *testing.T) {
			restore := stubBunVersionDetector(bunVersionResult{
				Version: version,
				Source:  "bun",
			}, nil)
			defer restore()

			detected, warning, err := checkBunMinimumSafeVersion(t.TempDir())
			if err != nil {
				t.Fatalf("checkBunMinimumSafeVersion failed: %v", err)
			}
			if detected != version {
				t.Fatalf("expected detected version %q, got %q", version, detected)
			}
			if warning != "" {
				t.Fatalf("expected no warning, got %q", warning)
			}
		})
	}
}

func TestRunPackageManagerSecurityCheckRunsForBun(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")

	restore := stubBunVersionDetector(bunVersionResult{
		Version: "1.3.0",
		Source:  "bun",
	}, nil)
	defer restore()

	RunPackageManagerSecurityCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.Bun,
	}, Options{})

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "exact = true")
	assertContains(t, content, "minimumReleaseAge = 259200")
}

func TestRunPackageManagerSecurityCheckWritesSupportedBunSettingsBelowMinimum(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")

	restore := stubBunVersionDetector(bunVersionResult{
		Version: "1.2.22",
		Source:  "bun",
	}, nil)
	defer restore()

	RunPackageManagerSecurityCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.Bun,
	}, Options{})

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "exact = true")
	assertContains(t, content, "save = true")
	if strings.Contains(content, "minimumReleaseAge") {
		t.Fatal("did not expect unsupported minimumReleaseAge to be written")
	}
}

func TestEnsureBunBestPracticesWritesOnlyVersionSupportedSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "bun.lock", "# bun lockfile\n")

	result, err := EnsureBunBestPractices(rootDir, "1.2.22", Options{})
	if err != nil {
		t.Fatalf("EnsureBunBestPractices failed: %v", err)
	}

	content := readBunConfig(t, rootDir)
	assertContains(t, content, "exact = true")
	assertContains(t, content, "save = true")
	if strings.Contains(content, "minimumReleaseAge") {
		t.Fatal("did not expect unsupported minimumReleaseAge to be written")
	}
	if !slices.Contains(result.Unsupported, "install.minimumReleaseAge") {
		t.Fatalf("expected install.minimumReleaseAge to be unsupported, got %#v", result.Unsupported)
	}
}

func stubBunVersionDetector(result bunVersionResult, err error) func() {
	previous := detectBunVersion
	detectBunVersion = func(rootDir string) (bunVersionResult, error) {
		return result, err
	}
	return func() {
		detectBunVersion = previous
	}
}

func readBunConfig(t *testing.T, rootDir string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(rootDir, bunfigFile))
	if err != nil {
		t.Fatalf("read bun config: %v", err)
	}
	return string(data)
}

func assertContains(t *testing.T, value string, want string) {
	t.Helper()

	if !strings.Contains(value, want) {
		t.Fatalf("expected %q to contain %q", value, want)
	}
}
