package cli

import (
	"fmt"

	"gocar/internal/project"
	"gocar/internal/util"
)

// FmtCommand fmt 命令
type FmtCommand struct{}

// Run 执行 fmt 命令
func (c *FmtCommand) Run(args []string) error {
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

	fmt.Printf("Formatting '%s'...\n", appName)
	if err := util.RunCommand(projectRoot, "go", append([]string{"fmt"}, packages...)...); err != nil {
		return fmt.Errorf("go fmt failed: %w", err)
	}
	fmt.Println("Format passed")
	return nil
}

// Help 返回帮助信息
func (c *FmtCommand) Help() string {
	return `gocar fmt - Format Go code

USAGE:
    gocar fmt [packages...]

DESCRIPTION:
    Runs go fmt. This command may modify files.

EXAMPLES:
    gocar fmt
    gocar fmt ./internal/...
`
}
