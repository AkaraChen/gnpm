package pmcombo

import (
	"reflect"
	"testing"
)

func TestAddCommand(t *testing.T) {
	tests := []struct {
		name     string
		pm       PackageManager
		opts     AddOptions
		expected []string
	}{
		{
			name:     "npm add simple",
			pm:       NPM,
			opts:     AddOptions{Packages: []string{"lodash"}},
			expected: []string{"install", "lodash"},
		},
		{
			name:     "npm add dev",
			pm:       NPM,
			opts:     AddOptions{Packages: []string{"lodash"}, Dev: true},
			expected: []string{"install", "--save-dev", "lodash"},
		},
		{
			name:     "pnpm add simple",
			pm:       PNPM,
			opts:     AddOptions{Packages: []string{"lodash"}},
			expected: []string{"add", "lodash"},
		},
		{
			name:     "pnpm add dev",
			pm:       PNPM,
			opts:     AddOptions{Packages: []string{"lodash"}, Dev: true},
			expected: []string{"add", "-D", "lodash"},
		},
		{
			name:     "yarn add simple",
			pm:       Yarn,
			opts:     AddOptions{Packages: []string{"lodash"}},
			expected: []string{"add", "lodash"},
		},
		{
			name:     "yarn add dev exact",
			pm:       Yarn,
			opts:     AddOptions{Packages: []string{"lodash"}, Dev: true, Exact: true},
			expected: []string{"add", "-D", "-E", "lodash"},
		},
		{
			name:     "bun add dev",
			pm:       Bun,
			opts:     AddOptions{Packages: []string{"lodash"}, Dev: true},
			expected: []string{"add", "-d", "lodash"},
		},
		{
			name:     "deno add",
			pm:       Deno,
			opts:     AddOptions{Packages: []string{"lodash"}},
			expected: []string{"add", "npm:lodash"},
		},
		{
			name:     "npm add multiple",
			pm:       NPM,
			opts:     AddOptions{Packages: []string{"lodash", "express"}},
			expected: []string{"install", "lodash", "express"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewAddCommand(tt.opts)
			result, err := cmd.Concat(tt.pm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAddCommandError(t *testing.T) {
	cmd := NewAddCommand(AddOptions{Packages: []string{}})
	_, err := cmd.Concat(NPM)
	if err == nil {
		t.Error("expected error for empty packages")
	}
}
