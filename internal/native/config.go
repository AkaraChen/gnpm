package native

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/user/fnpm/internal/logger"
)

// ConfigOptions for config operations
type ConfigOptions struct {
	Dir     string
	Action  string // get, set, delete, list
	Key     string
	Value   string
	Global  bool // operate on ~/.npmrc instead of project .npmrc
	DryRun  bool
}

// Config manages .npmrc configuration
func Config(opts ConfigOptions) error {
	switch opts.Action {
	case "get":
		return configGet(opts)
	case "set":
		return configSet(opts)
	case "delete":
		return configDelete(opts)
	case "list":
		return configList(opts)
	default:
		return fmt.Errorf("unknown config action: %s", opts.Action)
	}
}

// configGet retrieves a config value
func configGet(opts ConfigOptions) error {
	if opts.Key == "" {
		return fmt.Errorf("no key specified")
	}

	// Read from both project and user level, project takes precedence
	configs := mergeConfigs(opts.Dir)

	if value, ok := configs[opts.Key]; ok {
		logger.Plainln("%s", value)
		return nil
	}

	// Key not found - silent (npm behavior)
	return nil
}

// configSet sets a config value
func configSet(opts ConfigOptions) error {
	if opts.Key == "" {
		return fmt.Errorf("no key specified")
	}

	rcPath := getNpmrcPath(opts.Dir, opts.Global)

	if opts.DryRun {
		logger.DryRun(fmt.Sprintf("set %s=%s in %s", opts.Key, opts.Value, rcPath), opts.Dir)
		return nil
	}

	// Read existing config
	configs, lines := readNpmrc(rcPath)

	// Update or add the key
	configs[opts.Key] = opts.Value

	// Write back
	return writeNpmrc(rcPath, configs, lines, opts.Key)
}

// configDelete removes a config key
func configDelete(opts ConfigOptions) error {
	if opts.Key == "" {
		return fmt.Errorf("no key specified")
	}

	rcPath := getNpmrcPath(opts.Dir, opts.Global)

	if opts.DryRun {
		logger.DryRun(fmt.Sprintf("delete %s from %s", opts.Key, rcPath), opts.Dir)
		return nil
	}

	// Read existing config
	configs, lines := readNpmrc(rcPath)

	if _, ok := configs[opts.Key]; !ok {
		// Key doesn't exist, nothing to do
		return nil
	}

	delete(configs, opts.Key)

	// Write back, removing the key
	return writeNpmrcDelete(rcPath, lines, opts.Key)
}

// configList lists all config values
func configList(opts ConfigOptions) error {
	configs := mergeConfigs(opts.Dir)

	if len(configs) == 0 {
		logger.Dim("(no configuration)")
		return nil
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		logger.Plainln("%s=%s", k, configs[k])
	}

	return nil
}

// getNpmrcPath returns the path to .npmrc file
func getNpmrcPath(dir string, global bool) string {
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return ".npmrc"
		}
		return filepath.Join(home, ".npmrc")
	}
	return filepath.Join(dir, ".npmrc")
}

// readNpmrc reads .npmrc file and returns configs map and original lines
func readNpmrc(path string) (map[string]string, []string) {
	configs := make(map[string]string)
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return configs, lines
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		// Skip comments and empty lines
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
			continue
		}

		// Parse key=value
		if key, value, ok := strings.Cut(line, "="); ok {
			configs[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}

	return configs, lines
}

// mergeConfigs reads both project and user .npmrc, merging them
func mergeConfigs(dir string) map[string]string {
	// Start with user-level config
	userPath := getNpmrcPath(dir, true)
	configs, _ := readNpmrc(userPath)

	// Override with project-level config
	projectPath := getNpmrcPath(dir, false)
	projectConfigs, _ := readNpmrc(projectPath)

	for k, v := range projectConfigs {
		configs[k] = v
	}

	return configs
}

// writeNpmrc writes config back to .npmrc, preserving comments and order
func writeNpmrc(path string, configs map[string]string, originalLines []string, updatedKey string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var newLines []string
	keyWritten := false

	// Process original lines, updating the key if found
	for _, line := range originalLines {
		trimmed := strings.TrimSpace(line)

		// Keep comments and empty lines as-is
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
			newLines = append(newLines, line)
			continue
		}

		// Check if this line is the key we're updating
		if key, _, ok := strings.Cut(line, "="); ok {
			if strings.TrimSpace(key) == updatedKey {
				newLines = append(newLines, fmt.Sprintf("%s=%s", updatedKey, configs[updatedKey]))
				keyWritten = true
				continue
			}
		}

		newLines = append(newLines, line)
	}

	// If key wasn't found in original, append it
	if !keyWritten {
		newLines = append(newLines, fmt.Sprintf("%s=%s", updatedKey, configs[updatedKey]))
	}

	// Write to file
	content := strings.Join(newLines, "\n")
	if len(newLines) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// writeNpmrcDelete writes config back to .npmrc, removing a key
func writeNpmrcDelete(path string, originalLines []string, deleteKey string) error {
	var newLines []string

	for _, line := range originalLines {
		trimmed := strings.TrimSpace(line)

		// Keep comments and empty lines as-is
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
			newLines = append(newLines, line)
			continue
		}

		// Skip the line if it's the key we're deleting
		if key, _, ok := strings.Cut(line, "="); ok {
			if strings.TrimSpace(key) == deleteKey {
				continue
			}
		}

		newLines = append(newLines, line)
	}

	// Write to file
	content := strings.Join(newLines, "\n")
	if len(newLines) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return os.WriteFile(path, []byte(content), 0644)
}
