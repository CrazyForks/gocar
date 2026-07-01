package cli

import (
	"fmt"
	"os/exec"

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
			fmt.Printf("OK  %s: valid\n", config.ConfigFileName)
			fmt.Printf("  profiles: %v\n", cfg.ListProfiles())
		}
	}

	if !ok {
		return fmt.Errorf("doctor found issues")
	}
	fmt.Println("All checks passed")
	return nil
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
