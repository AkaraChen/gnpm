# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build ./cmd/gnpm          # Build binary to current directory
make build                   # Build to bin/gnpm
make test                    # Run all tests
go test ./internal/pmcombo   # Run tests for specific package
go test -run TestAddCommand ./internal/pmcombo  # Run single test
go test ./fixture -v         # Run E2E fixture tests
go test ./fixture -bench=.   # Run benchmarks
make install                 # Install to $GOPATH/bin
```

## Architecture

gnpm is a unified CLI that wraps npm, yarn, pnpm, deno, and bun. It detects the package manager from lock files and translates commands.

### Three Execution Paths

Commands fall into three categories:

1. **PM Delegation** (`pmcombo` → `runner.Run`): Translates gnpm commands to PM-specific commands and executes via the detected package manager. Used by: install (with packages), remove, ci, update, create (with template), scaffold, publish, why.

2. **Native Go** (`native.*`): Implements functionality directly in Go without calling a PM. Used by: run, test, exec (local binaries), create (without args), config, registry.

3. **Fallback** (`native.Fallback`): For unknown commands, tries: package.json scripts → node_modules/.bin → system commands.

### Merged Commands

Several commands are merged with smart behavior:

- **install** (aliases: `i`, `a`, `add`): Without args installs all deps, with args adds packages
- **exec** (aliases: `x`, `npx`, `dlx`): Tries local binary first, falls back to dlx
- **create** (aliases: `c`, `init`): Without args creates package.json, with args scaffolds from template

### Core Packages

**internal/pmcombo** - Command translation layer. Each command file implements `Command` interface with `Concat(pm PackageManager) ([]string, error)` that generates PM-specific arguments. `types.go` defines the `PackageManager` enum (NPM, Yarn, YarnClassic, PNPM, Deno, Bun) and option structs.

**internal/native** - Go-native implementations that bypass PM executables:
- `run.go` - Executes scripts directly from package.json via shell
- `config.go` - Reads/writes .npmrc files directly
- `registry.go` - Manages npm registry configuration
- `exec.go` - Runs binaries from node_modules/.bin
- `init.go` - Creates package.json
- `fallback.go` - Handles unknown commands (scripts → binaries → system)

**internal/context** - Project detection. `detector.go` finds PM from lock files with priority order: bun.lockb → deno.lock → pnpm-lock.yaml → yarn.lock → package-lock.json. Yarn Classic vs Berry distinguished by .yarnrc.yml presence or `__metadata:` in yarn.lock.

**internal/workspace** - Monorepo support. `finder.go` globs workspace patterns from package.json `workspaces` field or pnpm-workspace.yaml. `fuzzy.go` provides interactive package selection for `-s` flag.

**internal/cli** - Cobra commands. `root.go` sets up persistent flags (-w, -s, -V, --pm, --dry-run), runs `context.Detect()` in PersistentPreRunE, and handles fallback for unknown commands.

**internal/runner** - Executes PM commands via os/exec, handles --dry-run output.

### Command Flow

```
gnpm <cmd> [args]
    │
    ├─► Known command
    │   ├─► PM delegation path
    │   │   CLI → pmcombo.Concat(pm) → runner.Run(pm, args)
    │   │                                    └─► npm/yarn/pnpm/bun/deno <args>
    │   │
    │   └─► Native path
    │       CLI → native.* → direct Go implementation
    │
    └─► Unknown command (fallback)
        CLI → native.Fallback
            ├─► 1. Check package.json scripts
            ├─► 2. Check node_modules/.bin
            └─► 3. Run as system command
```

### Fixture Testing

`fixture/` contains 14 test projects covering all PM × workspace combinations:
- Single package: npm, yarn-classic, yarn-berry, pnpm, deno, bun
- Monorepo (5 and 100 packages): npm, yarn-classic, yarn-berry, pnpm

Run `go run fixture/generate.go` to regenerate fixtures.
