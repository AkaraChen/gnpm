package native

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/fnpm/internal/logger"
)

// InitOptions for creating package.json
type InitOptions struct {
	Dir    string
	Yes    bool
	Name   string
	DryRun bool
}

// PackageJSONTemplate represents the package.json structure
type PackageJSONTemplate struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description,omitempty"`
	Main        string            `json:"main,omitempty"`
	Scripts     map[string]string `json:"scripts,omitempty"`
	Keywords    []string          `json:"keywords,omitempty"`
	Author      string            `json:"author,omitempty"`
	License     string            `json:"license,omitempty"`
}

// Init creates a new package.json file
func Init(opts InitOptions) error {
	pkgPath := filepath.Join(opts.Dir, "package.json")

	if _, err := os.Stat(pkgPath); err == nil {
		return fmt.Errorf("package.json already exists")
	}

	var pkg PackageJSONTemplate

	if opts.Yes {
		pkg = getDefaults(opts.Dir, opts.Name)
	} else {
		var err error
		pkg, err = promptForValues(opts.Dir, opts.Name)
		if err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	if opts.DryRun {
		logger.DryRun(fmt.Sprintf("create %s", pkgPath), opts.Dir)
		logger.Plainln(string(data))
		return nil
	}

	if err := os.WriteFile(pkgPath, append(data, '\n'), 0644); err != nil {
		return err
	}

	logger.Success("created %s", pkgPath)
	return nil
}

func getDefaults(dir string, name string) PackageJSONTemplate {
	if name == "" {
		name = filepath.Base(dir)
	}
	return PackageJSONTemplate{
		Name:    name,
		Version: "1.0.0",
		Main:    "index.js",
		Scripts: map[string]string{
			"test": "echo \"Error: no test specified\" && exit 1",
		},
		License: "ISC",
	}
}

func promptForValues(dir string, defaultName string) (PackageJSONTemplate, error) {
	reader := bufio.NewReader(os.Stdin)

	if defaultName == "" {
		defaultName = filepath.Base(dir)
	}

	pkg := PackageJSONTemplate{
		Scripts: map[string]string{},
	}

	pkg.Name = prompt(reader, "package name", defaultName)
	pkg.Version = prompt(reader, "version", "1.0.0")
	pkg.Description = prompt(reader, "description", "")
	pkg.Main = prompt(reader, "entry point", "index.js")

	testCmd := prompt(reader, "test command", "")
	if testCmd != "" {
		pkg.Scripts["test"] = testCmd
	} else {
		pkg.Scripts["test"] = "echo \"Error: no test specified\" && exit 1"
	}

	keywords := prompt(reader, "keywords", "")
	if keywords != "" {
		pkg.Keywords = strings.Split(keywords, ",")
		for i := range pkg.Keywords {
			pkg.Keywords[i] = strings.TrimSpace(pkg.Keywords[i])
		}
	}

	pkg.Author = prompt(reader, "author", "")
	pkg.License = prompt(reader, "license", "ISC")

	return pkg, nil
}

func prompt(reader *bufio.Reader, field string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s: (%s) ", field, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", field)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}
