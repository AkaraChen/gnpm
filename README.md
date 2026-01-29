# gnpm

A fast, unified CLI that wraps npm, yarn, pnpm, deno, and bun. Write one command, run anywhere.

## Why?

Different projects use different package managers. Instead of remembering `npm install`, `yarn add`, `pnpm add`, `bun add`... just use `gnpm add`.

gnpm detects your package manager from lock files and translates commands automatically.

## Install

```bash
go install github.com/AkaraChen/gnpm/cmd/gnpm@latest
```

## Usage

```bash
# Works in any project - gnpm detects the package manager
gnpm install          # npm install / yarn install / pnpm install / bun install
gnpm i react          # Add react (i is alias for install)
gnpm add -D typescript
gnpm remove lodash
gnpm run build
gnpm test

# Unknown commands are resolved automatically
gnpm build            # Runs "build" script from package.json
gnpm eslint .         # Runs eslint from node_modules/.bin
gnpm ls               # Falls back to system command
```

## Detection

gnpm detects your package manager from lock files:

| Lock File | Package Manager |
|-----------|-----------------|
| `bun.lockb` | Bun |
| `deno.lock` | Deno |
| `pnpm-lock.yaml` | pnpm |
| `yarn.lock` + `.yarnrc.yml` | Yarn (Berry) |
| `yarn.lock` + `.yarnrc` | Yarn Classic |
| `package-lock.json` | npm |

## Commands

### Package Management

| Command | Aliases | Description |
|---------|---------|-------------|
| `gnpm install` | `i` | Install all dependencies |
| `gnpm install <pkg>` | `i`, `a`, `add` | Add a package |
| `gnpm install -D <pkg>` | | Add a dev dependency |
| `gnpm remove <pkg>` | `rm`, `un`, `uninstall` | Remove a package |
| `gnpm update` | `up`, `upgrade` | Update packages |
| `gnpm ci` | | Clean install (frozen lockfile) |

### Scripts & Execution

| Command | Aliases | Description |
|---------|---------|-------------|
| `gnpm run <script>` | `r` | Run a script from package.json |
| `gnpm test` | `t` | Run test script |
| `gnpm exec <cmd>` | `x`, `npx`, `dlx` | Execute binary (local or download) |

The `exec` command first looks for binaries in `node_modules/.bin`, then falls back to downloading and executing (like npx/dlx).

### Configuration

| Command | Aliases | Description |
|---------|---------|-------------|
| `gnpm config list` | | List all config |
| `gnpm config get <key>` | | Get a config value |
| `gnpm config set <key> <value>` | | Set a config value |
| `gnpm registry [url]` | `reg` | Get or set registry |

### Project Setup

| Command | Aliases | Description |
|---------|---------|-------------|
| `gnpm create` | `c`, `init` | Initialize package.json |
| `gnpm create <template>` | `c`, `init` | Create project from template |
| `gnpm scaffold` | `sc` | Scaffold using create-akrc |

The `create` command without arguments creates a package.json (like `npm init`). With a template name, it scaffolds a new project.

### Other

| Command | Aliases | Description |
|---------|---------|-------------|
| `gnpm publish` | `pub` | Publish to npm |
| `gnpm why <pkg>` | | Show why a package is installed |
| `gnpm view <pkg>` | `v`, `info`, `show` | Open package on npm |
| `gnpm use <pm>@<version>` | | Switch PM version via corepack |

## Aliases Quick Reference

| Full Command | Short Aliases |
|--------------|---------------|
| `install` | `i`, `a`, `add` |
| `remove` | `rm`, `un`, `uninstall` |
| `update` | `up`, `upgrade` |
| `run` | `r` |
| `test` | `t` |
| `exec` | `x`, `npx`, `dlx` |
| `create` | `c`, `init` |
| `registry` | `reg` |
| `publish` | `pub` |
| `view` | `v`, `info`, `show` |
| `scaffold` | `sc` |

## Default Command Fallback

When you run an unknown command, gnpm tries to resolve it:

1. **Scripts** - Check if it's a script in package.json
2. **Binaries** - Check if it's in node_modules/.bin
3. **System** - Run as a system command

```bash
gnpm build      # → runs "build" script if defined
gnpm eslint .   # → runs ./node_modules/.bin/eslint if installed
gnpm ls         # → runs system ls command
```

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Run in workspace root |
| `-s, --select` | Fuzzy select a workspace package |
| `--pm <pm>` | Override detected package manager |
| `--dry-run` | Print command without executing |
| `-V, --verbose` | Verbose output |

## Monorepo Support

gnpm supports npm/yarn/pnpm workspaces:

```bash
# Run in workspace root
gnpm install -w

# Fuzzy select a package to run command in
gnpm run build -s
```

## License

MIT
