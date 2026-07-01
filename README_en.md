# gocar, a cargo for Go

> A “Rust Cargo–like” scaffolding and CLI tool for Go projects, providing a clean experience for project initialization and building.

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go](https://img.shields.io/badge/go-1.25+-yellow.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Linux%20|%20macOS%20|%20Windows-blue.svg)](https://github.com/uselibrary/gocar)

**[简体中文](README.md)** | **[English](README_en.md)**

## Installation

> `git` is a prerequisite for some commands. Please make sure it is installed.

### Install via binary (recommended)

Download the prebuilt binary for your OS from the [Releases page](https://github.com/uselibrary/gocar/releases). Extract it and move it into a directory in your `$PATH`:

```
/usr/local/bin/ # Unix-like systems, e.g. Linux or macOS
C:\Program Files\ # Windows, may need to set environment variables
```

On Unix-like systems, ensure the binary has executable permissions (root required):

```
chown root:root /usr/local/bin/gocar
chmod +x /usr/local/bin/gocar
```

### Or build from source

```
git clone https://github.com/uselibrary/gocar.git
cd gocar
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gocar ./cmd/gocar
sudo mv gocar /usr/local/bin/
sudo chown root:root /usr/local/bin/gocar
sudo chmod +x /usr/local/bin/gocar
```

## Quick Start

```
# Create a standard Go application
gocar new myapp

# Enter the project directory
cd myapp

# Build the project
gocar build

# Run the project
gocar run

# Check the project
gocar check

# Clean build artifacts
gocar clean
```

## Command Reference

### Create a new project

**`gocar new <appName>`**

Create a standard Go application. `gocar new myapp` does not ask you to choose a template; it creates the default layout for most applications:

```
<appName>/
├── cmd/
│   └── <appName>/
│       └── main.go
├── internal/
├── go.mod
├── README.md
├── .gitignore
└── .git/
```

> Note: Created projects do not include `.gocar.toml` by default. Use `gocar init` to generate it manually.

> `<appName>` is the project name, used as the directory name, module name, and output executable name. The `bin/` directory is created automatically on the first `gocar build`.

> Project name rules:
>
> - Must start with a letter
> - May contain only letters, digits, underscore `_`, or hyphen `-`
> - Reserved names are not allowed: `test`, `main`, `init`, `internal`, `vendor`

### Build / Compile

**`gocar build [--release] [--target <os>/<arch>] [--with-cgo] [--help]`**

Build the executable:

- `gocar build` builds a Debug binary (default)
- `gocar build --release` builds a Release binary (enables `CGO_ENABLED=0`, `ldflags="-s -w"`, and `-trimpath`)
- `gocar build --profile <name>` builds with `[profile.<name>]` from `.gocar.toml`
- `gocar build --target <os>/<arch>` cross-compiles for the specified platform
- `gocar build --release --target <os>/<arch>` cross-compiles in Release mode for the specified platform
- `gocar build --with-cgo` forces CGO to be enabled (sets `CGO_ENABLED=1`)
- `gocar build --help` shows help information

Build behavior:

| Mode                    | Equivalent command                                           |
| ----------------------- | ------------------------------------------------------------ |
| debug (default)         | `go build -o bin/debug/<os>-<arch>/<appName> <entry>` |
| --release               | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/release/<os>-<arch>/<appName> <entry>` |
| --target                | `GOOS=<os> GOARCH=<arch> go build -o bin/debug/<os>-<arch>/<appName> <entry>` |
| --release --target      | `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build -ldflags="-s -w" -trimpath -o bin/release/<os>-<arch>/<appName> <entry>` |
| --with-cgo              | `CGO_ENABLED=1 go build -o bin/debug/<os>-<arch>/<appName> <entry>` |
| --release --with-cgo    | `CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o bin/release/<os>-<arch>/<appName> <entry>` |

Examples:

```
# Debug build (default)
gocar build

# Release build (enables CGO_ENABLED=0, ldflags="-s -w" and trimpath)
gocar build --release

# Build with a custom profile
gocar build --profile ci

# Cross-compile for a target OS/arch, e.g. Linux AMD64
gocar build --target linux/amd64

# Release cross-compile for Windows AMD64 (enables CGO_ENABLED=0, ldflags="-s -w" and trimpath)
gocar build --release --target windows/amd64

# Force enable CGO
gocar build --with-cgo

# Release build with CGO enabled
gocar build --release --with-cgo

# Show help
gocar build --help
```

### Common commands

**`gocar run [args...]`**

Run the current project directly (uses `go run`).

Examples:

```
# Run the project
gocar run

# Pass arguments to the app
gocar run --port 8080
```

**`gocar clean`**

Remove build artifacts in the `bin/` directory.

Example:

```
gocar clean
# Cleaned build artifacts for '<appName>'
```

**`gocar test [OPTIONS] [packages...]`**

Run project tests. By default this is equivalent to `go test ./...`.

Common options:
- `--coverage` runs tests with coverage (`go test -cover`)
- `--race` enables the race detector
- `--bench <pattern>` runs matching benchmarks
- Unrecognized test arguments are passed through to `go test`; you can also use `--` as an explicit separator

Examples:

```
gocar test
gocar test ./internal/...
gocar test --coverage
gocar test --bench .
gocar test -run TestConfig
gocar test -- -run TestConfig
```

**`gocar check [--no-test] [--race]`**

Run project checks. By default this runs `go fmt ./...`, `go vet ./...`, and `go test ./...`.

Examples:

```
gocar check
gocar check --race
gocar check --no-test
```

**`gocar commands`**

List built-in commands and custom commands defined in `.gocar.toml`.

**`gocar doctor`**

Check Go/Git, project detection, and `.gocar.toml` validation.

**`gocar help`**

Show help information.

**`gocar version`**

Show version information.

### Package management

**`gocar add <package>...`**

Add or update dependencies:

- `gocar add <package>` add the specified dependency
- `gocar update <package>` update the specified dependency
- `gocar update` update all dependencies and run `gocar tidy`
- `gocar tidy` tidy `go.mod` and `go.sum`
- `gocar add` is equivalent to `go get <package>...` and updates `go.mod` and `go.sum`

Dependency behavior:

| Command                     | Equivalent                      |
| --------------------------- | ------------------------------- |
| `gocar add <package>...`    | `go get <package>...`           |
| `gocar update [package]...` | `go get -u [package]...`        |
| `gocar update`              | `go get -u ./... & go mod tidy` |
| `gocar tidy`                | `go mod tidy`                   |

Examples:

```
# Add a dependency
gocar add github.com/gin-gonic/gin

# Update all dependencies
gocar update

# Update a specific dependency
gocar update github.com/gin-gonic/gin

# Tidy dependencies
gocar tidy
# Successfully tidied go.mod
```

### Configuration File

**`gocar init`**

Generate a `.gocar.toml` configuration file in the current project. Settings in the config file take priority over gocar's auto-detection.

Example:
```bash
# Generate config file in existing project
gocar init
# Created .gocar.toml in /path/to/project
```

**Configuration file structure:**

```toml
# gocar project configuration file

# project configuration

[project]
# Project name, uses directory name if empty
name = ""

# Project version
# version = "1.0.0"

# Build configuration
[build]
# Build entry path (relative to project root)
# standard layout defaults to "cmd/<appName>"
entry = "cmd/myapp"

# Output directory
output = "bin"

# Additional ldflags, appended to profile ldflags
# For example: "-X main.version=1.0.0"
ldflags = ""

# Build tags
# tags = ["jsoniter", "sonic"]

# Additional environment variables
# extra_env = ["GOPROXY=https://goproxy.cn"]

# Run configuration
[run]
# Run entry path, uses build.entry if empty
entry = ""

# Default run arguments
# args = ["-config", "config.yaml"]

# Debug build configuration
# Usage: gocar build (default)
[profile.debug]
# ldflags = ""              # Debug has no ldflags by default
# gcflags = "all=-N -l"     # Disable optimizations for easier debugging
# trimpath = false          # Retain path information
# cgo_enabled = true        # Follow system default
# race = false              # Race detection (may significantly reduce performance)
# Release build configuration
# Usage: gocar build --release
[profile.release]
ldflags = "-s -w"           # Strip symbol table and debug information
# gcflags = ""              # Compiler flags
trimpath = true             # Remove path information
cgo_enabled = false         # Disable CGO to generate static binary
# race = false              # Race detection

# Custom build profile
# Usage: gocar build --profile ci
# [profile.ci]
# trimpath = true
# race = true

# Custom commands
# Format: command_name = "shell command to execute"
# Usage: gocar <command_name>
# Commands are executed in the project root directory
[commands]
# Code checking
vet = "go vet ./..."

# Code formatting
fmt = "go fmt ./..."

# lint = "golangci-lint run"
# doc = "godoc -http=:6060"
# proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"
```

**Configuration options:**

| Option | Description |
|--------|-------------|
| `[project].name` | Custom project name, uses directory name if empty |
| `[project].version` | **Project version**, auto-injected via `-X main.version=<version>` at build time |
| `[build].entry` | **Custom build entry path**, e.g., `cmd/myapp` instead of default `cmd/<appName>` (the project name) |
| `[build].ldflags` | Additional ldflags, appended to profile ldflags |
| `[build].tags` | Build tags list |
| `[build].extra_env` | Additional environment variables |
| `[run].entry` | Run entry path, uses `build.entry` if empty |
| `[run].args` | Default run arguments |
| `[profile.debug]` | Debug build mode parameters |
| `[profile.release]` | Release build mode parameters |
| `[commands]` | Custom command mappings |

**Profile options:**

| Option | Description | Debug Default | Release Default |
|--------|-------------|---------------|------------------|
| `ldflags` | Linker flags | `""` | `"-s -w"` |
| `gcflags` | Compiler flags | `""` | `""` |
| `trimpath` | Remove path info | `false` | `true` |
| `cgo_enabled` | Enable CGO | `nil` (system) | `false` |
| `race` | Race detection | `false` | `false` |

### Custom Commands

After defining commands in the `[commands]` section of `.gocar.toml`, you can execute them directly. `test` and `check` are built-in commands, so you usually do not need to define them yourself:

```bash
# Code checking
gocar vet

# Code formatting
gocar fmt

# Custom lint
gocar lint
```

Command output is displayed in real-time to the terminal. You can define any custom commands, for example:

```toml
[commands]
lint = "golangci-lint run"
doc = "godoc -http=:6060"
proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"
dev = "air"  # Hot reload
```

#### Overriding Built-in Commands

Custom commands can override most built-in commands, giving you full control over your project's build and run workflow:

| Command Type | Commands | Can Override |
|--------------|----------|-------------|
| Protected | `new`, `init` | ❌ No |
| Project | `build`, `run`, `clean`, `add`, `update`, `tidy`, `test`, `check`, `commands`, `doctor` | ✅ Yes |

> **Protected commands** (`new`, `init`) cannot be overridden because `new` runs before project creation (no config file exists yet), and `init` generates the config file itself.

Example: Override built-in `build` and `clean` commands

```toml
[commands]
# Use Makefile for building
build = "make build"

# Custom clean logic
clean = "make clean && rm -rf dist/"

# Use docker-compose to run
run = "docker-compose up"
```

When running `gocar build`, if a `build` command is defined in the config file, the custom command will be executed instead of the built-in build logic.

------

The `main.go` template content for a new project is:

```
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hello, gocar! A golang project scaffolding tool for <appName>.")
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
```

------

## License

MIT License
