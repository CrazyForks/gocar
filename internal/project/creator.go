package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gocar/internal/util"
)

const (
	// ModeStandard is the default application layout for new projects.
	ModeStandard = "standard"
)

// Creator 项目创建器
type Creator struct {
	Name string
}

// NewCreator 创建项目创建器
func NewCreator(name string) *Creator {
	return &Creator{
		Name: name,
	}
}

// Create 创建项目
func (c *Creator) Create() error {
	// Check if directory already exists
	if _, err := os.Stat(c.Name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", c.Name)
	}

	if err := c.createStandardProject(); err != nil {
		return err
	}

	// Initialize git
	if err := util.InitGit(c.Name); err != nil {
		fmt.Printf("Warning: Failed to initialize git: %v\n", err)
	}

	return nil
}

// createStandardProject 创建标准 Go 应用。
func (c *Creator) createStandardProject() error {
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "cmd", c.Name),
		filepath.Join(c.Name, "internal"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	if err := util.WriteFile(filepath.Join(c.Name, "cmd", c.Name, "main.go"), c.mainTemplate()); err != nil {
		return err
	}

	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), c.readmeTemplate()); err != nil {
		return err
	}

	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), c.gitignoreTemplate()); err != nil {
		return err
	}

	return nil
}

// mainTemplate 生成 main.go 内容
func (c *Creator) mainTemplate() string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, gocar! A golang project scaffolding tool for %s.")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
`, c.Name)
}

// readmeTemplate 生成 README.md 内容
func (c *Creator) readmeTemplate() string {
	return fmt.Sprintf(`# %s

A Go project created with gocar.

## Build

`+"```bash"+`
# Debug build (current platform)
gocar build

# Release build (current platform)
gocar build --release

# Cross-compile for Linux on AMD64
gocar build --target linux/amd64
`+"```"+`

## Run

`+"```bash"+`
gocar run
`+"```"+`

## Check

`+"```bash"+`
gocar check
`+"```"+`

## Output Structure

`+"```"+`
bin/
├── debug/
│   └── <os>-<arch>/
│       └── %s
└── release/
    └── <os>-<arch>/
        └── %s
`+"```"+`

Build artifacts are organized by:
- **Build mode**: debug or release
- **Target platform**: OS and architecture (e.g., linux-amd64, darwin-arm64)

Examples:
- Debug build for current platform: `+"`./bin/debug/linux-amd64/%s`"+`
- Release build for Windows: `+"`./bin/release/windows-amd64/%s.exe`"+`
`, c.Name, c.Name, c.Name, c.Name, c.Name)
}

// gitignoreTemplate 生成 .gitignore 内容
func (c *Creator) gitignoreTemplate() string {
	return fmt.Sprintf(`# Binaries
%s
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
`, c.Name)
}
