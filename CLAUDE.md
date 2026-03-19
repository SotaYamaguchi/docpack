# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`docpack` is a CLI tool for organizing meeting documents before and after customer project meetings. It automates file renaming and collection based on project-specific prefixes.

## Development Commands

All commands are run from the `mtg/` directory:

```bash
cd mtg

# Build
make build          # Creates ./mtg binary

# Testing
make test           # Run tests with coverage
go test -v          # Run tests verbose
go test -run TestFunctionName  # Run single test

# Linting
make lint           # Run golangci-lint
golangci-lint run   # Direct invocation

# Installation (local)
make install        # Install to ~/bin/mtg with config
make uninstall      # Remove binary and config
```

## Architecture

### Single-File CLI (`main.go`)

All functionality is in `mtg/main.go` (~340 lines). No modules or packages.

**Subcommands:**
- `prep` - Rename files (main→date) and collect for pre-meeting
- `memo` - Rename files (main→date_MTG後) and collect for post-meeting
- `list` - Show configured projects from config.json
- `completion` - Generate zsh completion script

**Core flow:**
1. Parse subcommand and flags (`-project`, `-prefix`, `-dir`)
2. Resolve prefix from project name via `~/.config/mtg/config.json`
3. Execute file operations:
   - `renameFiles()` - Replace "main" in filenames with current date (YYYYMMDD format)
   - `collectFiles()` - Move files matching prefix pattern to destination folder

**Configuration:**
- User config: `~/.config/mtg/config.json` (maps project names to file prefixes)
- Sample: `mtg/config.sample.json` (committed to repo)
- **Important:** Actual `config.json` contains customer information and is `.gitignore`d

### Testing Strategy

Tests in `main_test.go` use temporary directories and files:
- `loadConfig` - JSON parsing and validation
- `resolvePrefix` - Project name to prefix resolution
- `renameFiles` - File renaming with date/suffix
- `collectFiles` - File moving to destination folder

Test files must handle cleanup with `defer func() { _ = os.RemoveAll(tmpDir) }()` pattern to satisfy errcheck linter.

## CI/CD

GitHub Actions (`.github/workflows/test.yml`) runs three jobs:
1. **Lint** - golangci-lint v2.11 (config: `.golangci.yml`)
2. **Test** - Unit tests with race detector, coverage displayed in Step Summary
3. **Build** - Binary compilation check

## pre-commit Hooks

`.pre-commit-config.yaml` runs on commit:
- golangci-lint (only on changed Go files)
- Standard checks (trailing whitespace, EOF, YAML syntax, file size)

Setup: `pre-commit install`

## Key Constraints

- **No external dependencies** - Uses only Go standard library
- **Stateless** - No database, all config from JSON file
- **File-based** - Operates on filesystem directly using glob patterns
- **Date format** - Always YYYYMMDD (time.Now().Format("20060102"))
- **Japanese output** - All user messages and help text in Japanese
