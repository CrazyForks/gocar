package cli

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gocar/internal/config"
	"gocar/internal/project"
)

// DoctorCommand doctor 命令
type DoctorCommand struct{}

// Run 执行 doctor 命令
func (c *DoctorCommand) Run(args []string) error {
	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h":
			fmt.Print(c.Help())
			return nil
		default:
			return fmt.Errorf("unknown option '%s' (run 'gocar doctor --help' for usage)", arg)
		}
	}

	fmt.Println("gocar doctor")

	ok := true
	ok = printToolCheck("go", true) && ok
	ok = printToolCheck("git", false) && ok
	ok = printGoEnvironment() && ok

	projectRoot, appName, projectMode, err := project.DetectProject()
	if err != nil {
		fmt.Printf("ERR project: %v\n", err)
		ok = false
	} else {
		fmt.Printf("OK  project: %s (%s mode)\n", appName, projectMode)
		fmt.Printf("  root: %s\n", projectRoot)

		cfg, err := config.Load(projectRoot)
		if err != nil {
			fmt.Printf("ERR %s: %v\n", config.ConfigFileName, err)
			ok = false
		} else if err := cfg.Validate(projectRoot); err != nil {
			fmt.Printf("ERR %s: %v\n", config.ConfigFileName, err)
			ok = false
		} else {
			if config.Exists(projectRoot) {
				fmt.Printf("OK  %s: valid\n", config.ConfigFileName)
			} else {
				fmt.Printf("OK  %s: not found, using defaults\n", config.ConfigFileName)
			}
			fmt.Printf("  profiles: %v\n", cfg.ListProfiles())
			if !printCommandOverrideWarnings(cfg) {
				ok = false
			}
			if !printEntryCheck(projectRoot, "build.entry", cfg.GetBuildEntryForApp(appName), true) {
				ok = false
			}
			runEntry := cfg.GetRunEntryForApp(appName)
			if runEntry != cfg.GetBuildEntryForApp(appName) {
				if !printEntryCheck(projectRoot, "run.entry", runEntry, true) {
					ok = false
				}
			}
		}
	}

	if !ok {
		return fmt.Errorf("doctor found issues")
	}
	fmt.Println("All checks passed")
	return nil
}

func printGoEnvironment() bool {
	versionCmd := exec.Command("go", "version")
	versionOutput, err := versionCmd.Output()
	if err != nil {
		fmt.Printf("ERR go version: %v\n", err)
		return false
	}
	fmt.Printf("OK  go.version: %s\n", strings.TrimSpace(string(versionOutput)))

	envCmd := exec.Command("go", "env", "GOPROXY", "GOMODCACHE", "CGO_ENABLED")
	envOutput, err := envCmd.Output()
	if err != nil {
		fmt.Printf("ERR go env: %v\n", err)
		return false
	}
	lines := strings.Split(strings.TrimSpace(string(envOutput)), "\n")
	keys := []string{"GOPROXY", "GOMODCACHE", "CGO_ENABLED"}
	for i, key := range keys {
		value := ""
		if i < len(lines) {
			value = lines[i]
		}
		fmt.Printf("OK  go.env.%s: %s\n", key, value)
	}
	return true
}

func printCommandOverrideWarnings(cfg *config.GocarConfig) bool {
	for name := range cfg.Commands {
		if isBuiltInCommandName(name) && !isProtectedCommand(name) {
			fmt.Printf("WARN command %q overrides built-in command\n", name)
		}
	}
	return true
}

func printEntryCheck(projectRoot, label, entry string, requireMain bool) bool {
	path := entry
	if path == "" {
		fmt.Printf("ERR %s: empty\n", label)
		return false
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(projectRoot, path)
	}
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("ERR %s: %s does not exist\n", label, entry)
		return false
	}
	fmt.Printf("OK  %s: %s\n", label, entry)
	if requireMain {
		if !entryHasMainPackage(path) {
			fmt.Printf("ERR %s: %s is not a main package\n", label, entry)
			return false
		}
		fmt.Printf("OK  %s.package: main\n", label)
	}
	return true
}

func entryHasMainPackage(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	fset := token.NewFileSet()
	if !stat.IsDir() {
		file, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
		return err == nil && file.Name != nil && file.Name.Name == "main"
	}

	pkgs, err := parser.ParseDir(fset, path, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, parser.PackageClauseOnly)
	if err != nil {
		return false
	}
	_, ok := pkgs["main"]
	return ok
}

func printToolCheck(name string, required bool) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		if required {
			fmt.Printf("ERR %s: not found\n", name)
			return false
		}
		fmt.Printf("WARN %s: not found (optional)\n", name)
		return true
	}
	fmt.Printf("OK  %s: %s\n", name, path)
	return true
}

// Help 返回帮助信息
func (c *DoctorCommand) Help() string {
	return `gocar doctor - Check project and toolchain setup

USAGE:
    gocar doctor

DESCRIPTION:
    Checks Go, Git, project detection, and .gocar.toml validation.
`
}
