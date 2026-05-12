package security

import (
	stdcontext "context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	projectcontext "github.com/AkaraChen/gnpm/internal/context"
	"github.com/AkaraChen/gnpm/internal/logger"
	"github.com/AkaraChen/gnpm/internal/pmcombo"
	"gopkg.in/yaml.v3"
)

const (
	pnpmWorkspaceFile        = "pnpm-workspace.yaml"
	pnpmMinimumReleaseAgeMin = 1440
	pnpmVersionProbeTimeout  = 3 * time.Second

	// pnpmMinimumSafeVersion is the version floor for the pnpm security
	// settings gnpm writes below. pnpm introduced the relevant settings across
	// v10: strictDepBuilds in 10.3, dangerouslyAllowAllBuilds in 10.9,
	// minimumReleaseAge in 10.16, trustPolicy in 10.21, and allowBuilds plus
	// blockExoticSubdeps in 10.26. minimumReleaseAgeStrict lands in pnpm 11,
	// which is also where pnpm's defaults move toward these safer settings, so
	// gnpm treats 11.0.0 as the minimum safe version.
	pnpmMinimumSafeVersion = "11.0.0"
)

var (
	checksMu sync.Mutex
	checks   []<-chan struct{}
)

// Options controls package-manager security checks.
type Options struct {
	DryRun  bool
	Verbose bool
}

// Result describes changes made by a security check.
type Result struct {
	Changed  bool
	Settings []string
	Warnings []string
}

type pnpmVersionResult struct {
	Version string
	Source  string
}

var detectPNPMVersion = detectPNPMVersionFromSystem

// WarnPackageManagerVersionIfUnsafe warns when the active package manager is too old.
func WarnPackageManagerVersionIfUnsafe(ctx *projectcontext.ProjectContext, opts Options) {
	if ctx == nil || ctx.PackageManager != pmcombo.PNPM {
		return
	}

	warning, err := checkPNPMMinimumSafeVersion(ctx.RootDir)
	if err != nil {
		if opts.Verbose {
			logger.Warn("pnpm version check failed: %v", err)
		}
		return
	}
	if warning != "" {
		logger.Warn("%s", warning)
	}
}

func checkPNPMMinimumSafeVersion(rootDir string) (string, error) {
	detected, err := detectPNPMVersion(rootDir)
	if err != nil {
		return "", err
	}

	if compareSemver(detected.Version, pnpmMinimumSafeVersion) >= 0 {
		return "", nil
	}

	return fmt.Sprintf(
		"pnpm %s from %s is below gnpm's minimum safe pnpm version %s; upgrade pnpm to %s or newer to enable all recommended security settings",
		detected.Version,
		detected.Source,
		pnpmMinimumSafeVersion,
		pnpmMinimumSafeVersion,
	), nil
}

func detectPNPMVersionFromSystem(rootDir string) (pnpmVersionResult, error) {
	var errors []string

	for _, probe := range []struct {
		source string
		name   string
		args   []string
	}{
		{source: "corepack", name: "corepack", args: []string{"pnpm", "-v"}},
		{source: "pnpm", name: "pnpm", args: []string{"-v"}},
	} {
		output, err := runPNPMVersionProbe(rootDir, probe.name, probe.args...)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		version, ok := parseSemver(output)
		if !ok {
			errors = append(errors, fmt.Sprintf("%s returned an unparseable version: %q", probe.source, strings.TrimSpace(output)))
			continue
		}

		return pnpmVersionResult{Version: version.String(), Source: probe.source}, nil
	}

	return pnpmVersionResult{}, fmt.Errorf("detect pnpm version: %s", strings.Join(errors, "; "))
}

func runPNPMVersionProbe(rootDir string, name string, args ...string) (string, error) {
	if _, err := exec.LookPath(name); err != nil {
		return "", err
	}

	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), pnpmVersionProbeTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = rootDir
	cmd.Env = append(os.Environ(), "COREPACK_ENABLE_DOWNLOAD_PROMPT=0")

	output, err := cmd.CombinedOutput()
	if ctx.Err() == stdcontext.DeadlineExceeded {
		return "", fmt.Errorf("%s %s timed out", name, strings.Join(args, " "))
	}
	if err != nil {
		details := strings.TrimSpace(string(output))
		if details != "" {
			return "", fmt.Errorf("%s %s failed: %w: %s", name, strings.Join(args, " "), err, details)
		}
		return "", fmt.Errorf("%s %s failed: %w", name, strings.Join(args, " "), err)
	}

	return string(output), nil
}

// StartPackageManagerBestPracticeCheck runs a best-practice check in the background.
func StartPackageManagerBestPracticeCheck(ctx *projectcontext.ProjectContext, opts Options) {
	done := make(chan struct{})

	checksMu.Lock()
	checks = append(checks, done)
	checksMu.Unlock()

	go func() {
		defer close(done)

		if ctx == nil || ctx.PackageManager != pmcombo.PNPM {
			return
		}

		result, err := EnsurePNPMBestPractices(ctx.RootDir, opts)
		if err != nil {
			logger.Warn("pnpm security config check failed: %v", err)
			return
		}

		for _, warning := range result.Warnings {
			logger.Warn("%s", warning)
		}

		if result.Changed && opts.Verbose {
			logger.Success("pnpm security config updated: %s", strings.Join(result.Settings, ", "))
		}
	}()
}

// WaitForPackageManagerBestPracticeChecks waits for all started checks to finish.
func WaitForPackageManagerBestPracticeChecks() {
	checksMu.Lock()
	pending := checks
	checks = nil
	checksMu.Unlock()

	for _, done := range pending {
		<-done
	}
}

// EnsurePNPMBestPractices enforces pnpm supply-chain security settings.
func EnsurePNPMBestPractices(rootDir string, opts Options) (Result, error) {
	var result Result
	if rootDir == "" {
		return result, fmt.Errorf("empty project root")
	}

	path := filepath.Join(rootDir, pnpmWorkspaceFile)
	doc, err := readOrCreateYAMLDocument(path)
	if err != nil {
		return result, err
	}

	root := ensureMappingDocument(doc)
	if ensureBool(root, "dangerouslyAllowAllBuilds", false) {
		result.Settings = append(result.Settings, "dangerouslyAllowAllBuilds=false")
	}
	if ensureBool(root, "strictDepBuilds", true) {
		result.Settings = append(result.Settings, "strictDepBuilds=true")
	}
	if ensureMap(root, "allowBuilds") {
		result.Settings = append(result.Settings, "allowBuilds={}")
	}
	if ensureBool(root, "blockExoticSubdeps", true) {
		result.Settings = append(result.Settings, "blockExoticSubdeps=true")
	}
	if ensureMinInt(root, "minimumReleaseAge", pnpmMinimumReleaseAgeMin) {
		result.Settings = append(result.Settings, "minimumReleaseAge=1440")
	}
	if ensureBool(root, "minimumReleaseAgeStrict", true) {
		result.Settings = append(result.Settings, "minimumReleaseAgeStrict=true")
	}
	if ensureString(root, "trustPolicy", "no-downgrade") {
		result.Settings = append(result.Settings, "trustPolicy=no-downgrade")
	}

	if _, err := os.Stat(filepath.Join(rootDir, "pnpm-lock.yaml")); err != nil {
		if os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, "pnpm lockfile is missing; run pnpm install and commit pnpm-lock.yaml")
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

func readOrCreateYAMLDocument(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newMappingDocument(), nil
		}
		return nil, err
	}

	var doc yaml.Node
	if len(strings.TrimSpace(string(data))) == 0 {
		return newMappingDocument(), nil
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func newMappingDocument() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		}},
	}
}

func ensureMappingDocument(doc *yaml.Node) *yaml.Node {
	if doc.Kind != yaml.DocumentNode {
		doc.Kind = yaml.DocumentNode
	}
	if len(doc.Content) == 0 || doc.Content[0] == nil {
		doc.Content = []*yaml.Node{{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
		}}
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		root.Kind = yaml.MappingNode
		root.Tag = "!!map"
		root.Value = ""
		root.Content = nil
	}
	if root.Tag == "" {
		root.Tag = "!!map"
	}
	return root
}

func writeYAMLDocument(path string, doc *yaml.Node) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ensureBool(root *yaml.Node, key string, desired bool) bool {
	value := findMapValue(root, key)
	if value != nil && value.Kind == yaml.ScalarNode && value.Tag == "!!bool" {
		current, err := strconv.ParseBool(value.Value)
		if err == nil && current == desired {
			return false
		}
	}

	setMapValue(root, key, boolNode(desired))
	return true
}

func ensureString(root *yaml.Node, key string, desired string) bool {
	value := findMapValue(root, key)
	if value != nil && value.Kind == yaml.ScalarNode && value.Value == desired {
		if value.Tag == "" {
			value.Tag = "!!str"
		}
		return false
	}

	setMapValue(root, key, stringNode(desired))
	return true
}

func ensureMinInt(root *yaml.Node, key string, min int) bool {
	value := findMapValue(root, key)
	if value != nil && value.Kind == yaml.ScalarNode && (value.Tag == "!!int" || value.Tag == "") {
		current, err := strconv.Atoi(value.Value)
		if err == nil && current >= min {
			if value.Tag == "" {
				value.Tag = "!!int"
			}
			return false
		}
	}

	setMapValue(root, key, intNode(min))
	return true
}

func ensureMap(root *yaml.Node, key string) bool {
	value := findMapValue(root, key)
	if value != nil && value.Kind == yaml.MappingNode {
		if value.Tag == "" {
			value.Tag = "!!map"
		}
		return false
	}

	setMapValue(root, key, mapNode())
	return true
}

func findMapValue(root *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Kind == yaml.ScalarNode && root.Content[i].Value == key {
			return root.Content[i+1]
		}
	}
	return nil
}

func setMapValue(root *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Kind == yaml.ScalarNode && root.Content[i].Value == key {
			root.Content[i+1] = value
			return
		}
	}

	root.Content = append(root.Content, stringNode(key), value)
}

func boolNode(value bool) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: strconv.FormatBool(value),
	}
}

func intNode(value int) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!int",
		Value: strconv.Itoa(value),
	}
}

func stringNode(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
}

func mapNode() *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.MappingNode,
		Tag:   "!!map",
		Style: yaml.FlowStyle,
	}
}

type semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

var semverPattern = regexp.MustCompile(`v?([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z.-]+))?`)

func parseSemver(value string) (semver, bool) {
	match := semverPattern.FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		return semver{}, false
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return semver{}, false
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return semver{}, false
	}
	patch, err := strconv.Atoi(match[3])
	if err != nil {
		return semver{}, false
	}

	return semver{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: match[4],
	}, true
}

func compareSemver(a string, b string) int {
	aVersion, aOK := parseSemver(a)
	bVersion, bOK := parseSemver(b)
	if !aOK || !bOK {
		return strings.Compare(a, b)
	}

	if aVersion.Major != bVersion.Major {
		return compareInt(aVersion.Major, bVersion.Major)
	}
	if aVersion.Minor != bVersion.Minor {
		return compareInt(aVersion.Minor, bVersion.Minor)
	}
	if aVersion.Patch != bVersion.Patch {
		return compareInt(aVersion.Patch, bVersion.Patch)
	}
	if aVersion.Prerelease == bVersion.Prerelease {
		return 0
	}
	if aVersion.Prerelease == "" {
		return 1
	}
	if bVersion.Prerelease == "" {
		return -1
	}
	return strings.Compare(aVersion.Prerelease, bVersion.Prerelease)
}

func compareInt(a int, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func (v semver) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		version += "-" + v.Prerelease
	}
	return version
}
