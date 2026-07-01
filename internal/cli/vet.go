package cli

import (
	"fmt"

	"gocar/internal/project"
	"gocar/internal/util"
)

// VetCommand vet 命令
type VetCommand struct{}

// Run 执行 vet 命令
func (c *VetCommand) Run(args []string) error {
	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h":
			fmt.Print(c.Help())
			return nil
		}
	}

	projectRoot, appName, _, err := project.DetectProject()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	packages := args
	if len(packages) == 0 {
		packages = []string{"./..."}
	}

	fmt.Printf("Vetting '%s'...\n", appName)
	if err := util.RunCommand(projectRoot, "go", append([]string{"vet"}, packages...)...); err != nil {
		return fmt.Errorf("go vet failed: %w", err)
	}
	fmt.Println("Vet passed")
	return nil
}

// Help 返回帮助信息
func (c *VetCommand) Help() string {
	return `gocar vet - Run go vet

USAGE:
    gocar vet [packages...]

EXAMPLES:
    gocar vet
    gocar vet ./internal/...
`
}
