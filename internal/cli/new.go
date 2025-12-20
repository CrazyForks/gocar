package cli

import (
	"fmt"
	"os"
	"strings"

	"gocar/internal/config"
	"gocar/internal/project"
)

// NewCommand new 命令
type NewCommand struct{}

// Run 执行 new 命令
func (c *NewCommand) Run(args []string) error {
	if len(args) < 1 {
		fmt.Println("Error: Missing project name")
		fmt.Println("Usage: gocar new <name> [--mode simple|project|<template>]")
		os.Exit(1)
	}

	// Check for help
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Print(c.Help())
		return nil
	}

	appName := args[0]

	// Validate project name
	if err := project.ValidateProjectName(appName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	mode := "simple" // default mode

	// Parse --mode flag
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--mode":
			if i+1 < len(args) {
				mode = args[i+1]
				i++ // skip next arg
			} else {
				fmt.Println("Error: --mode requires a value")
				os.Exit(1)
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Printf("Error: Unknown option '%s'\n", args[i])
				fmt.Println("Run 'gocar new --help' for usage.")
				os.Exit(1)
			}
		}
	}

	// 检查是否是内置模式
	if mode == "simple" || mode == "project" {
		fmt.Printf("Creating new %s project: %s\n", mode, appName)

		creator := project.NewCreator(appName, mode)
		if err := creator.Create(); err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSuccessfully created project '%s'\n", appName)
		fmt.Printf("\nTo get started:\n")
		fmt.Printf("    cd %s\n", appName)
		fmt.Printf("    gocar build\n")
		fmt.Printf("    gocar run\n")

		return nil
	}

	// 尝试从全局配置加载模板
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		fmt.Printf("Error loading global config: %v\n", err)
		os.Exit(1)
	}

	tpl, ok := globalCfg.GetTemplate(mode)
	if !ok {
		fmt.Printf("Error: Unknown mode or template '%s'\n", mode)
		fmt.Println("\nBuilt-in modes: simple, project")

		// 显示可用模板
		templates := globalCfg.ListTemplates()
		if len(templates) > 0 {
			fmt.Println("\nAvailable templates from global config:")
			for name, t := range templates {
				desc := t.Description
				if desc == "" {
					desc = "(no description)"
				}
				fmt.Printf("  %-12s  %s\n", name, desc)
			}
		} else {
			fmt.Println("\nNo custom templates defined.")
			fmt.Println("Run 'gocar config init' to create global config with example templates.")
		}
		os.Exit(1)
	}

	// 使用模板创建项目
	fmt.Printf("Creating project '%s' from template '%s' (base: %s)\n", appName, mode, tpl.Mode)

	creator := project.NewCreatorWithTemplate(appName, tpl)
	if err := creator.Create(); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully created project '%s' from template '%s'\n", appName, mode)
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
    gocar new <name> [--mode simple|project|<template>]

OPTIONS:
    --mode <mode>    Project mode or template name
                     Built-in: 'simple' (default), 'project'
                     Or use a template name from global config

EXAMPLES:
    gocar new myapp                   Create a simple project
    gocar new myapp --mode project    Create a project-mode project
    gocar new myapi --mode api        Create from 'api' template

TEMPLATES:
    Custom templates can be defined in ~/.gocar/config.toml
    Run 'gocar config init' to create config with example templates
    Run 'gocar config list' to see available templates
`
	return helpText
}
