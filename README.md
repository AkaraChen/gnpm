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
gnpm add lodash       # npm install lodash / yarn add lodash / pnpm add lodash
gnpm add -D typescript
gnpm remove lodash
gnpm run build
gnpm test
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

| Command | Description |
|---------|-------------|
| `gnpm install` | Install all dependencies |
| `gnpm add <pkg>` | Add a package |
| `gnpm add -D <pkg>` | Add a dev dependency |
| `gnpm remove <pkg>` | Remove a package |
| `gnpm update` | Update packages |
| `gnpm ci` | Clean install (frozen lockfile) |

### Scripts

| Command | Description |
|---------|-------------|
| `gnpm run <script>` | Run a script from package.json |
| `gnpm test` | Run test script |
| `gnpm exec <cmd>` | Execute a binary from node_modules/.bin |
| `gnpm dlx <pkg>` | Download and execute a package |

### Configuration

| Command | Description |
|---------|-------------|
| `gnpm config list` | List all config |
| `gnpm config get <key>` | Get a config value |
| `gnpm config set <key> <value>` | Set a config value |
| `gnpm registry [url]` | Get or set registry |

### Other

| Command | Description |
|---------|-------------|
| `gnpm init` | Initialize package.json |
| `gnpm create <template>` | Create project from template |
| `gnpm publish` | Publish to npm |
| `gnpm why <pkg>` | Show why a package is installed |
| `gnpm view <pkg>` | Open package on npm |
| `gnpm use <pm>@<version>` | Switch PM version via corepack |

## Flags

| Flag | Description |
|------|-------------|
| `-w, --workspace` | Run in workspace root |
| `-s, --select` | Fuzzy select a workspace package |
| `--pm <pm>` | Override detected package manager |
| `--dry-run` | Print command without executing |
| `-v, --verbose` | Verbose output |

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

