package cli

import (
	"fmt"
	"sort"

	"gocar/internal/config"
	"gocar/internal/project"
)

// CommandsCommand commands 命令
type CommandsCommand struct{}

// Run 执行 commands 命令
func (c *CommandsCommand) Run(args []string) error {
	for _, arg := range args {
		switch arg {
		case "help", "--help", "-h":
			fmt.Print(c.Help())
			return nil
		default:
			return fmt.Errorf("unknown option '%s' (run 'gocar commands --help' for usage)", arg)
		}
	}

	fmt.Println("Built-in commands:")
	for _, name := range builtInCommandNames() {
		info, _ := builtInCommandInfo(name)
		fmt.Printf("  %-12s %s\n", name, info.Description)
	}

	projectRoot, _, _, err := project.DetectProject()
	if err != nil {
		fmt.Println("\nCustom commands: unavailable outside a Go module")
		return nil
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", config.ConfigFileName, err)
	}

	names := make([]string, 0, len(cfg.Commands))
	for name := range cfg.Commands {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("\nCustom commands:")
	if len(names) == 0 {
		fmt.Println("  (none)")
		return nil
	}
	for _, name := range names {
		override := ""
		if isBuiltInCommandName(name) && !isProtectedCommand(name) {
			override = " (overrides built-in)"
		}
		fmt.Printf("  %s%s = %s\n", name, override, cfg.Commands[name])
	}

	return nil
}

func isBuiltInCommandName(name string) bool {
	for _, builtIn := range builtInCommandNames() {
		if name == builtIn {
			return true
		}
	}
	return false
}

// Help 返回帮助信息
func (c *CommandsCommand) Help() string {
	return `gocar commands - List built-in and custom commands

USAGE:
    gocar commands

DESCRIPTION:
    Lists built-in commands and commands defined in .gocar.toml.
`
}
