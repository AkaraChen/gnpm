package pmcombo

import (
	"reflect"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	tests := []struct {
		name     string
		pm       PackageManager
		opts     InstallOptions
		expected []string
	}{
		{
			name:     "npm install",
			pm:       NPM,
			opts:     InstallOptions{Frozen: false},
			expected: []string{"install"},
		},
		{
			name:     "npm ci",
			pm:       NPM,
			opts:     InstallOptions{Frozen: true},
			expected: []string{"ci"},
		},
		{
			name:     "pnpm install",
			pm:       PNPM,
			opts:     InstallOptions{Frozen: false},
			expected: []string{"install"},
		},
		{
			name:     "pnpm frozen",
			pm:       PNPM,
			opts:     InstallOptions{Frozen: true},
			expected: []string{"install", "--frozen-lockfile"},
		},
		{
			name:     "yarn install",
			pm:       Yarn,
			opts:     InstallOptions{Frozen: false},
			expected: []string{"install"},
		},
		{
			name:     "yarn frozen",
			pm:       Yarn,
			opts:     InstallOptions{Frozen: true},
			expected: []string{"install", "--immutable"},
		},
		{
			name:     "yarn classic frozen",
			pm:       YarnClassic,
			opts:     InstallOptions{Frozen: true},
			expected: []string{"install", "--frozen-lockfile"},
		},
		{
			name:     "bun install",
			pm:       Bun,
			opts:     InstallOptions{Frozen: false},
			expected: []string{"install"},
		},
		{
			name:     "bun frozen",
			pm:       Bun,
			opts:     InstallOptions{Frozen: true},
			expected: []string{"install", "--frozen-lockfile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewInstallCommand(tt.opts)
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
