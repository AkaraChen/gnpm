package security

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	bunfigFile                  = "bunfig.toml"
	bunMinimumReleaseAgeSeconds = 259200

	// bunMinimumSafeVersion is the version floor for the Bun security settings
	// gnpm writes below. Bun 1.3 introduced the Security Scanner API and
	// minimumReleaseAge, while exact package saves and lockfile persistence
	// predate those settings. gnpm treats 1.3.0 as the minimum safe Bun version.
	bunMinimumSafeVersion = "1.3.0"
)

type bunVersionResult struct {
	Version string
	Source  string
}

var detectBunVersion = detectBunVersionFromSystem

func checkBunMinimumSafeVersion(rootDir string) (string, string, error) {
	detected, err := detectBunVersion(rootDir)
	if err != nil {
		return "", "", err
	}

	if compareSemver(detected.Version, bunMinimumSafeVersion) >= 0 {
		return detected.Version, "", nil
	}

	return detected.Version, fmt.Sprintf(
		"bun %s from %s is below gnpm's minimum safe Bun version %s; upgrade Bun to %s or newer to enable all recommended security settings",
		detected.Version,
		detected.Source,
		bunMinimumSafeVersion,
		bunMinimumSafeVersion,
	), nil
}

func detectBunVersionFromSystem(rootDir string) (bunVersionResult, error) {
	output, err := runPackageManagerVersionProbe(rootDir, "bun", "--version")
	if err != nil {
		return bunVersionResult{}, fmt.Errorf("detect bun version: %w", err)
	}

	version, ok := parseSemver(output)
	if !ok {
		return bunVersionResult{}, fmt.Errorf("bun returned an unparseable version: %q", strings.TrimSpace(output))
	}

	return bunVersionResult{Version: version.String(), Source: "bun"}, nil
}

// EnsureBunBestPractices enforces Bun supply-chain security settings.
func EnsureBunBestPractices(rootDir string, version string, opts Options) (Result, error) {
	var result Result
	if rootDir == "" {
		return result, fmt.Errorf("empty project root")
	}

	path := filepath.Join(rootDir, bunfigFile)
	config, err := readOrCreateTOMLMap(path)
	if err != nil {
		return result, err
	}

	install := ensureTOMLTable(config, "install")
	if !supportsPMSetting(version, "1.0.0") {
		result.Unsupported = append(result.Unsupported, "install.exact")
	} else if ensureTOMLBool(install, "exact", true) {
		result.Settings = append(result.Settings, "install.exact=true")
	}
	if !supportsPMSetting(version, "1.3.0") {
		result.Unsupported = append(result.Unsupported, "install.minimumReleaseAge")
	} else if ensureTOMLMinInt(install, "minimumReleaseAge", bunMinimumReleaseAgeSeconds) {
		result.Settings = append(result.Settings, "install.minimumReleaseAge=259200")
	}

	lockfile := ensureNestedTOMLTable(config, "install", "lockfile")
	if !supportsPMSetting(version, "1.0.0") {
		result.Unsupported = append(result.Unsupported, "install.lockfile.save")
	} else if ensureTOMLBool(lockfile, "save", true) {
		result.Settings = append(result.Settings, "install.lockfile.save=true")
	}

	if !hasAnyFile(rootDir, "bun.lock", "bun.lockb") {
		result.Warnings = append(result.Warnings, "bun lockfile is missing; run bun install and commit bun.lock")
	}

	result.Changed = len(result.Settings) > 0
	if !result.Changed || opts.DryRun {
		return result, nil
	}

	return result, writeTOMLMap(path, config)
}

func hasAnyFile(rootDir string, names ...string) bool {
	for _, name := range names {
		if _, err := os.Stat(filepath.Join(rootDir, name)); err == nil {
			return true
		}
	}
	return false
}

func readOrCreateTOMLMap(path string) (map[string]interface{}, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var config map[string]interface{}
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	if config == nil {
		config = make(map[string]interface{})
	}
	return config, nil
}

func writeTOMLMap(path string, config map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	var data bytes.Buffer
	if err := toml.NewEncoder(&data).Encode(config); err != nil {
		return err
	}
	return os.WriteFile(path, data.Bytes(), 0644)
}

func ensureTOMLTable(config map[string]interface{}, key string) map[string]interface{} {
	if table, ok := config[key].(map[string]interface{}); ok {
		return table
	}

	table := make(map[string]interface{})
	config[key] = table
	return table
}

func ensureNestedTOMLTable(config map[string]interface{}, parent string, key string) map[string]interface{} {
	parentTable := ensureTOMLTable(config, parent)
	return ensureTOMLTable(parentTable, key)
}

func ensureTOMLBool(table map[string]interface{}, key string, desired bool) bool {
	if current, ok := table[key].(bool); ok && current == desired {
		return false
	}

	table[key] = desired
	return true
}

func ensureTOMLMinInt(table map[string]interface{}, key string, min int64) bool {
	if current, ok := tomlInt(table[key]); ok && current >= min {
		return false
	}

	table[key] = min
	return true
}

func tomlInt(value interface{}) (int64, bool) {
	switch current := value.(type) {
	case int:
		return int64(current), true
	case int8:
		return int64(current), true
	case int16:
		return int64(current), true
	case int32:
		return int64(current), true
	case int64:
		return current, true
	case uint:
		return int64(current), true
	case uint8:
		return int64(current), true
	case uint16:
		return int64(current), true
	case uint32:
		return int64(current), true
	case uint64:
		if current > uint64(^uint64(0)>>1) {
			return 0, false
		}
		return int64(current), true
	default:
		return 0, false
	}
}
