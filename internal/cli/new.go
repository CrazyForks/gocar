package cli

import (
	"fmt"

	"gocar/internal/project"
)

// NewCommand new 命令
type NewCommand struct{}

// Run 执行 new 命令
func (c *NewCommand) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing project name (usage: gocar new <name>)")
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Print(c.Help())
		return nil
	}

	appName := args[0]

	// Validate project name
	if err := project.ValidateProjectName(appName); err != nil {
		return err
	}

	if len(args) > 1 {
		return fmt.Errorf("unknown option '%s' (run 'gocar new --help' for usage)", args[1])
	}

	fmt.Printf("Creating new Go application: %s\n", appName)

	creator := project.NewCreator(appName)
	if err := creator.Create(); err != nil {
		return fmt.Errorf("error creating project: %w", err)
	}

	fmt.Printf("\nSuccessfully created project '%s'\n", appName)
	fmt.Printf("\nTo get started:\n")
	fmt.Printf("    cd %s\n", appName)
	fmt.Printf("    gocar build\n")
	fmt.Printf("    gocar run\n")

	return nil
}

// Help 返回帮助信息
func (c *NewCommand) Help() string {
	helpText := `gocar new - Create a new Go project

USAGE:
    gocar new <name>

DESCRIPTION:
    Creates a standard Go application using cmd/<name>/main.go and internal/.

EXAMPLES:
    gocar new myapp
`
	return helpText
}
