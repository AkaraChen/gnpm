package security

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	projectcontext "github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"gopkg.in/yaml.v3"
)

func TestEnsurePNPMBestPracticesCreatesMissingSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")

	result, err := EnsurePNPMBestPractices(rootDir, Options{})
	if err != nil {
		t.Fatalf("EnsurePNPMBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected missing settings to be written")
	}

	config := readWorkspaceConfig(t, rootDir)
	assertEqual(t, config["dangerouslyAllowAllBuilds"], false)
	assertEqual(t, config["strictDepBuilds"], true)
	assertEqual(t, config["blockExoticSubdeps"], true)
	assertEqual(t, config["minimumReleaseAge"], 1440)
	assertEqual(t, config["minimumReleaseAgeStrict"], true)
	assertEqual(t, config["trustPolicy"], "no-downgrade")
	if _, ok := config["allowBuilds"].(map[string]interface{}); !ok {
		t.Fatalf("expected allowBuilds to be an empty map, got %#v", config["allowBuilds"])
	}
}

func TestEnsurePNPMBestPracticesFixesUnsafeSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
	writeFile(t, rootDir, pnpmWorkspaceFile, `packages:
  - packages/*
dangerouslyAllowAllBuilds: true
strictDepBuilds: false
allowBuilds:
  esbuild: true
blockExoticSubdeps: false
minimumReleaseAge: 60
minimumReleaseAgeStrict: false
trustPolicy: off
`)

	result, err := EnsurePNPMBestPractices(rootDir, Options{})
	if err != nil {
		t.Fatalf("EnsurePNPMBestPractices failed: %v", err)
	}

	expectedSettings := []string{
		"dangerouslyAllowAllBuilds=false",
		"strictDepBuilds=true",
		"blockExoticSubdeps=true",
		"minimumReleaseAge=1440",
		"minimumReleaseAgeStrict=true",
		"trustPolicy=no-downgrade",
	}
	for _, setting := range expectedSettings {
		if !slices.Contains(result.Settings, setting) {
			t.Fatalf("expected %q in changed settings, got %#v", setting, result.Settings)
		}
	}

	config := readWorkspaceConfig(t, rootDir)
	assertEqual(t, config["dangerouslyAllowAllBuilds"], false)
	assertEqual(t, config["strictDepBuilds"], true)
	assertEqual(t, config["blockExoticSubdeps"], true)
	assertEqual(t, config["minimumReleaseAge"], 1440)
	assertEqual(t, config["minimumReleaseAgeStrict"], true)
	assertEqual(t, config["trustPolicy"], "no-downgrade")
	assertEqual(t, config["packages"], []interface{}{"packages/*"})

	allowBuilds, ok := config["allowBuilds"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected allowBuilds map, got %#v", config["allowBuilds"])
	}
	assertEqual(t, allowBuilds["esbuild"], true)
}

func TestEnsurePNPMBestPracticesKeepsStrongerMinimumReleaseAge(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
	writeFile(t, rootDir, pnpmWorkspaceFile, `dangerouslyAllowAllBuilds: false
strictDepBuilds: true
allowBuilds: {}
blockExoticSubdeps: true
minimumReleaseAge: 10080
minimumReleaseAgeStrict: true
trustPolicy: no-downgrade
`)

	result, err := EnsurePNPMBestPractices(rootDir, Options{})
	if err != nil {
		t.Fatalf("EnsurePNPMBestPractices failed: %v", err)
	}
	if result.Changed {
		t.Fatalf("expected no changes, got %#v", result.Settings)
	}

	config := readWorkspaceConfig(t, rootDir)
	assertEqual(t, config["minimumReleaseAge"], 10080)
}

func TestEnsurePNPMBestPracticesDryRunDoesNotWrite(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
	writeFile(t, rootDir, pnpmWorkspaceFile, "dangerouslyAllowAllBuilds: true\n")

	result, err := EnsurePNPMBestPractices(rootDir, Options{DryRun: true})
	if err != nil {
		t.Fatalf("EnsurePNPMBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected dry-run to report pending changes")
	}

	content, err := os.ReadFile(filepath.Join(rootDir, pnpmWorkspaceFile))
	if err != nil {
		t.Fatalf("read workspace config: %v", err)
	}
	if string(content) != "dangerouslyAllowAllBuilds: true\n" {
		t.Fatalf("dry-run changed file content:\n%s", content)
	}
}

func TestPNPMMinimumSafeVersion(t *testing.T) {
	if got, want := pnpmMinimumSafeVersion, "11.0.0"; got != want {
		t.Fatalf("expected minimum safe pnpm version %q, got %q", want, got)
	}
}

func TestCheckPNPMMinimumSafeVersionWarnsBelowMinimum(t *testing.T) {
	restore := stubPNPMVersionDetector(pnpmVersionResult{
		Version: "10.26.0",
		Source:  "corepack",
	}, nil)
	defer restore()

	warning, err := checkPNPMMinimumSafeVersion(t.TempDir())
	if err != nil {
		t.Fatalf("checkPNPMMinimumSafeVersion failed: %v", err)
	}
	if !strings.Contains(warning, "pnpm 10.26.0 from corepack") {
		t.Fatalf("expected detected version in warning, got %q", warning)
	}
	if !strings.Contains(warning, "11.0.0") {
		t.Fatalf("expected minimum version in warning, got %q", warning)
	}
}

func TestCheckPNPMMinimumSafeVersionAllowsMinimumOrNewer(t *testing.T) {
	for _, version := range []string{"11.0.0", "11.1.0"} {
		t.Run(version, func(t *testing.T) {
			restore := stubPNPMVersionDetector(pnpmVersionResult{
				Version: version,
				Source:  "pnpm",
			}, nil)
			defer restore()

			warning, err := checkPNPMMinimumSafeVersion(t.TempDir())
			if err != nil {
				t.Fatalf("checkPNPMMinimumSafeVersion failed: %v", err)
			}
			if warning != "" {
				t.Fatalf("expected no warning, got %q", warning)
			}
		})
	}
}

func TestCompareSemverTreatsPrereleaseAsBelowRelease(t *testing.T) {
	if compareSemver("11.0.0-rc.1", "11.0.0") >= 0 {
		t.Fatal("expected 11.0.0-rc.1 to compare below 11.0.0")
	}
	if compareSemver("11.0.0", "11.0.0-rc.1") <= 0 {
		t.Fatal("expected 11.0.0 to compare above 11.0.0-rc.1")
	}
}

func TestStartPackageManagerBestPracticeCheckRunsInBackground(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")

	StartPackageManagerBestPracticeCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.PNPM,
	}, Options{})
	WaitForPackageManagerBestPracticeChecks()

	config := readWorkspaceConfig(t, rootDir)
	assertEqual(t, config["trustPolicy"], "no-downgrade")
	assertEqual(t, config["blockExoticSubdeps"], true)
}

func stubPNPMVersionDetector(result pnpmVersionResult, err error) func() {
	previous := detectPNPMVersion
	detectPNPMVersion = func(rootDir string) (pnpmVersionResult, error) {
		return result, err
	}
	return func() {
		detectPNPMVersion = previous
	}
}

func writeFile(t *testing.T, rootDir, name, content string) {
	t.Helper()

	path := filepath.Join(rootDir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func readWorkspaceConfig(t *testing.T, rootDir string) map[string]interface{} {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(rootDir, pnpmWorkspaceFile))
	if err != nil {
		t.Fatalf("read workspace config: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("unmarshal workspace config: %v", err)
	}
	return config
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if !equalYAMLValue(got, want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func equalYAMLValue(got, want interface{}) bool {
	gotBytes, gotErr := yaml.Marshal(got)
	wantBytes, wantErr := yaml.Marshal(want)
	if gotErr != nil || wantErr != nil {
		return false
	}
	return string(gotBytes) == string(wantBytes)
}
