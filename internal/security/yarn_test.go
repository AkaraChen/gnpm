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

func TestEnsureYarnBestPracticesCreatesMissingSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")

	result, err := EnsureYarnBestPractices(rootDir, yarnMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureYarnBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected missing settings to be written")
	}

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["defaultSemverRangePrefix"], "")
	assertEqual(t, config["enableScripts"], false)
	assertEqual(t, config["npmMinimalAgeGate"], 1440)
	assertEqual(t, config["npmPublishProvenance"], true)
}

func TestEnsureYarnBestPracticesFixesUnsafeSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")
	writeFile(t, rootDir, yarnRCFile, `nodeLinker: node-modules
defaultSemverRangePrefix: ^
enableScripts: true
npmMinimalAgeGate: 60
npmPublishProvenance: false
`)

	result, err := EnsureYarnBestPractices(rootDir, yarnMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureYarnBestPractices failed: %v", err)
	}

	expectedSettings := []string{
		"defaultSemverRangePrefix=",
		"enableScripts=false",
		"npmMinimalAgeGate=1440",
		"npmPublishProvenance=true",
	}
	for _, setting := range expectedSettings {
		if !slices.Contains(result.Settings, setting) {
			t.Fatalf("expected %q in changed settings, got %#v", setting, result.Settings)
		}
	}

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["nodeLinker"], "node-modules")
	assertEqual(t, config["defaultSemverRangePrefix"], "")
	assertEqual(t, config["enableScripts"], false)
	assertEqual(t, config["npmMinimalAgeGate"], 1440)
	assertEqual(t, config["npmPublishProvenance"], true)
}

func TestEnsureYarnBestPracticesKeepsStrongerMinimalAgeGate(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")
	writeFile(t, rootDir, yarnRCFile, `defaultSemverRangePrefix: ""
enableScripts: false
npmMinimalAgeGate: 3d
npmPublishProvenance: true
`)

	result, err := EnsureYarnBestPractices(rootDir, yarnMinimumSafeVersion, Options{})
	if err != nil {
		t.Fatalf("EnsureYarnBestPractices failed: %v", err)
	}
	if result.Changed {
		t.Fatalf("expected no changes, got %#v", result.Settings)
	}

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["npmMinimalAgeGate"], "3d")
}

func TestEnsureYarnBestPracticesDryRunDoesNotWrite(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")
	writeFile(t, rootDir, yarnRCFile, "enableScripts: true\n")

	result, err := EnsureYarnBestPractices(rootDir, yarnMinimumSafeVersion, Options{DryRun: true})
	if err != nil {
		t.Fatalf("EnsureYarnBestPractices failed: %v", err)
	}
	if !result.Changed {
		t.Fatal("expected dry-run to report pending changes")
	}

	content, err := os.ReadFile(filepath.Join(rootDir, yarnRCFile))
	if err != nil {
		t.Fatalf("read yarn config: %v", err)
	}
	if string(content) != "enableScripts: true\n" {
		t.Fatalf("dry-run changed file content:\n%s", content)
	}
}

func TestYarnMinimumSafeVersion(t *testing.T) {
	if got, want := yarnMinimumSafeVersion, "4.10.0"; got != want {
		t.Fatalf("expected minimum safe Yarn version %q, got %q", want, got)
	}
}

func TestCheckYarnMinimumSafeVersionWarnsBelowMinimum(t *testing.T) {
	restore := stubYarnVersionDetector(yarnVersionResult{
		Version: "4.9.0",
		Source:  "corepack",
	}, nil)
	defer restore()

	version, warning, err := checkYarnMinimumSafeVersion(t.TempDir())
	if err != nil {
		t.Fatalf("checkYarnMinimumSafeVersion failed: %v", err)
	}
	if version != "4.9.0" {
		t.Fatalf("expected detected version %q, got %q", "4.9.0", version)
	}
	if !strings.Contains(warning, "yarn 4.9.0 from corepack") {
		t.Fatalf("expected detected version in warning, got %q", warning)
	}
	if !strings.Contains(warning, "4.10.0") {
		t.Fatalf("expected minimum version in warning, got %q", warning)
	}
}

func TestCheckYarnMinimumSafeVersionAllowsMinimumOrNewer(t *testing.T) {
	for _, version := range []string{"4.10.0", "4.11.0"} {
		t.Run(version, func(t *testing.T) {
			restore := stubYarnVersionDetector(yarnVersionResult{
				Version: version,
				Source:  "yarn",
			}, nil)
			defer restore()

			detected, warning, err := checkYarnMinimumSafeVersion(t.TempDir())
			if err != nil {
				t.Fatalf("checkYarnMinimumSafeVersion failed: %v", err)
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

func TestRunPackageManagerSecurityCheckRunsForYarn(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")

	restore := stubYarnVersionDetector(yarnVersionResult{
		Version: "4.10.0",
		Source:  "corepack",
	}, nil)
	defer restore()

	RunPackageManagerSecurityCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.Yarn,
	}, Options{})

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["enableScripts"], false)
	assertEqual(t, config["npmPublishProvenance"], true)
}

func TestRunPackageManagerSecurityCheckWritesSupportedYarnSettingsBelowMinimum(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")

	restore := stubYarnVersionDetector(yarnVersionResult{
		Version: "4.9.0",
		Source:  "corepack",
	}, nil)
	defer restore()

	RunPackageManagerSecurityCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.Yarn,
	}, Options{})

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["defaultSemverRangePrefix"], "")
	assertEqual(t, config["enableScripts"], false)
	assertEqual(t, config["npmPublishProvenance"], true)
	if _, ok := config["npmMinimalAgeGate"]; ok {
		t.Fatal("did not expect unsupported npmMinimalAgeGate to be written")
	}
}

func TestRunPackageManagerSecurityCheckSkipsYarnClassicWrite(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")

	restore := stubYarnVersionDetector(yarnVersionResult{
		Version: "1.22.22",
		Source:  "yarn",
	}, nil)
	defer restore()

	RunPackageManagerSecurityCheck(&projectcontext.ProjectContext{
		RootDir:        rootDir,
		PackageManager: pmcombo.YarnClassic,
	}, Options{})

	if _, err := os.Stat(filepath.Join(rootDir, yarnRCFile)); !os.IsNotExist(err) {
		t.Fatalf("expected no yarn berry config to be written for classic yarn, stat err=%v", err)
	}
}

func TestEnsureYarnBestPracticesWritesOnlyVersionSupportedSettings(t *testing.T) {
	rootDir := t.TempDir()
	writeFile(t, rootDir, "yarn.lock", "# yarn lockfile\n")

	result, err := EnsureYarnBestPractices(rootDir, "4.9.0", Options{})
	if err != nil {
		t.Fatalf("EnsureYarnBestPractices failed: %v", err)
	}

	config := readYarnConfig(t, rootDir)
	assertEqual(t, config["defaultSemverRangePrefix"], "")
	assertEqual(t, config["enableScripts"], false)
	assertEqual(t, config["npmPublishProvenance"], true)
	if _, ok := config["npmMinimalAgeGate"]; ok {
		t.Fatal("did not expect unsupported npmMinimalAgeGate to be written")
	}
	if !slices.Contains(result.Unsupported, "npmMinimalAgeGate") {
		t.Fatalf("expected npmMinimalAgeGate to be unsupported, got %#v", result.Unsupported)
	}
}

func stubYarnVersionDetector(result yarnVersionResult, err error) func() {
	previous := detectYarnVersion
	detectYarnVersion = func(rootDir string) (yarnVersionResult, error) {
		return result, err
	}
	return func() {
		detectYarnVersion = previous
	}
}

func readYarnConfig(t *testing.T, rootDir string) map[string]interface{} {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(rootDir, yarnRCFile))
	if err != nil {
		t.Fatalf("read yarn config: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("unmarshal yarn config: %v", err)
	}
	return config
}
