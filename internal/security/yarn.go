package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	yarnRCFile               = ".yarnrc.yml"
	yarnNpmMinimalAgeGateMin = 1440

	// yarnMinimumSafeVersion is the version floor for the Yarn Berry security
	// settings gnpm writes below. npmPublishProvenance landed in Yarn 4.9.0,
	// while npmMinimalAgeGate and npmPreapprovedPackages landed in Yarn 4.10.0.
	// enableScripts and defaultSemverRangePrefix predate those settings, so
	// gnpm treats 4.10.0 as the minimum safe Yarn Berry version.
	yarnMinimumSafeVersion = "4.10.0"
)

type yarnVersionResult struct {
	Version string
	Source  string
}

var detectYarnVersion = detectYarnVersionFromSystem

func checkYarnMinimumSafeVersion(rootDir string) (string, string, error) {
	detected, err := detectYarnVersion(rootDir)
	if err != nil {
		return "", "", err
	}

	if compareSemver(detected.Version, yarnMinimumSafeVersion) >= 0 {
		return detected.Version, "", nil
	}

	return detected.Version, fmt.Sprintf(
		"yarn %s from %s is below gnpm's minimum safe Yarn version %s; upgrade Yarn to %s or newer to enable all recommended security settings",
		detected.Version,
		detected.Source,
		yarnMinimumSafeVersion,
		yarnMinimumSafeVersion,
	), nil
}

func detectYarnVersionFromSystem(rootDir string) (yarnVersionResult, error) {
	var errors []string

	for _, probe := range []struct {
		source string
		name   string
		args   []string
	}{
		{source: "corepack", name: "corepack", args: []string{"yarn", "-v"}},
		{source: "yarn", name: "yarn", args: []string{"-v"}},
	} {
		output, err := runPackageManagerVersionProbe(rootDir, probe.name, probe.args...)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		version, ok := parseSemver(output)
		if !ok {
			errors = append(errors, fmt.Sprintf("%s returned an unparseable version: %q", probe.source, strings.TrimSpace(output)))
			continue
		}

		return yarnVersionResult{Version: version.String(), Source: probe.source}, nil
	}

	return yarnVersionResult{}, fmt.Errorf("detect yarn version: %s", strings.Join(errors, "; "))
}

// EnsureYarnBestPractices enforces Yarn Berry supply-chain security settings.
func EnsureYarnBestPractices(rootDir string, version string, opts Options) (Result, error) {
	var result Result
	if rootDir == "" {
		return result, fmt.Errorf("empty project root")
	}

	path := filepath.Join(rootDir, yarnRCFile)
	doc, err := readOrCreateYAMLDocument(path)
	if err != nil {
		return result, err
	}

	root := ensureMappingDocument(doc)
	if !supportsPMSetting(version, "2.0.0") {
		result.Unsupported = append(result.Unsupported, "defaultSemverRangePrefix")
	} else if ensureString(root, "defaultSemverRangePrefix", "") {
		result.Settings = append(result.Settings, "defaultSemverRangePrefix=")
	}
	if !supportsPMSetting(version, "2.0.0") {
		result.Unsupported = append(result.Unsupported, "enableScripts")
	} else if ensureBool(root, "enableScripts", false) {
		result.Settings = append(result.Settings, "enableScripts=false")
	}
	if !supportsPMSetting(version, "4.10.0") {
		result.Unsupported = append(result.Unsupported, "npmMinimalAgeGate")
	} else if ensureMinYarnDuration(root, "npmMinimalAgeGate", yarnNpmMinimalAgeGateMin) {
		result.Settings = append(result.Settings, "npmMinimalAgeGate=1440")
	}
	if !supportsPMSetting(version, "4.9.0") {
		result.Unsupported = append(result.Unsupported, "npmPublishProvenance")
	} else if ensureBool(root, "npmPublishProvenance", true) {
		result.Settings = append(result.Settings, "npmPublishProvenance=true")
	}

	if _, err := os.Stat(filepath.Join(rootDir, "yarn.lock")); err != nil {
		if os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, "yarn lockfile is missing; run yarn install and commit yarn.lock")
		} else {
			return result, err
		}
	}

	result.Changed = len(result.Settings) > 0
	if !result.Changed || opts.DryRun {
		return result, nil
	}

	return result, writeYAMLDocument(path, doc)
}

func ensureMinYarnDuration(root *yaml.Node, key string, minMinutes int) bool {
	value := findMapValue(root, key)
	if value != nil && value.Kind == yaml.ScalarNode {
		current, ok := yarnDurationMinutes(value.Value)
		if ok && current >= minMinutes {
			if value.Tag == "" {
				value.Tag = "!!int"
			}
			return false
		}
	}

	setMapValue(root, key, intNode(minMinutes))
	return true
}

func yarnDurationMinutes(value string) (int, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false
	}

	unit := trimmed[len(trimmed)-1:]
	multiplier := 1
	number := trimmed
	switch unit {
	case "d":
		multiplier = 1440
		number = strings.TrimSpace(trimmed[:len(trimmed)-1])
	case "h":
		multiplier = 60
		number = strings.TrimSpace(trimmed[:len(trimmed)-1])
	case "m":
		number = strings.TrimSpace(trimmed[:len(trimmed)-1])
	}

	minutes, err := strconv.Atoi(number)
	if err != nil {
		return 0, false
	}
	return minutes * multiplier, true
}
