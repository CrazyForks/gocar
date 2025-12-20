package cli

import (
	"fmt"
	"os"

	"gocar/internal/config"
)

// ConfigCommand config 命令
type ConfigCommand struct{}

// Run 执行 config 命令
func (c *ConfigCommand) Run(args []string) error {
	// 如果没有子命令，显示帮助
	if len(args) < 1 {
		fmt.Print(c.Help())
		return nil
	}

	subCmd := args[0]

	switch subCmd {
	case "help", "--help", "-h":
		fmt.Print(c.Help())
		return nil
	case "init":
		return c.initConfig()
	case "path":
		return c.showPath()
	case "list":
		return c.listTemplates()
	case "edit":
		return c.editConfig()
	default:
		fmt.Printf("Error: Unknown subcommand '%s'\n", subCmd)
		fmt.Println("Run 'gocar config --help' for usage.")
		os.Exit(1)
	}

	return nil
}

// initConfig 初始化全局配置
func (c *ConfigCommand) initConfig() error {
	if config.GlobalConfigExists() {
		configPath, _ := config.GetGlobalConfigPath()
		fmt.Printf("Global config already exists at: %s\n", configPath)
		fmt.Println("Use 'gocar config edit' to modify it.")
		return nil
	}

	if err := config.SaveGlobalConfig(); err != nil {
		fmt.Printf("Error creating global config: %v\n", err)
		os.Exit(1)
	}

	configPath, _ := config.GetGlobalConfigPath()
	fmt.Printf("Created global config at: %s\n", configPath)
	fmt.Println("\nYou can now:")
	fmt.Println("  - Edit the config to add custom templates")
	fmt.Println("  - Use 'gocar new <name> --mode <template>' to create projects from templates")
	fmt.Println("  - Run 'gocar config list' to see available templates")

	return nil
}

// showPath 显示配置文件路径
func (c *ConfigCommand) showPath() error {
	configPath, err := config.GetGlobalConfigPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Global config path: %s\n", configPath)

	if config.GlobalConfigExists() {
		fmt.Println("Status: exists")
	} else {
		fmt.Println("Status: not created")
		fmt.Println("Run 'gocar config init' to create it.")
	}

	return nil
}

// listTemplates 列出所有模板
func (c *ConfigCommand) listTemplates() error {
	if !config.GlobalConfigExists() {
		fmt.Println("No global config found.")
		fmt.Println("Run 'gocar config init' to create one with example templates.")
		return nil
	}

	cfg, err := config.LoadGlobalConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	templates := cfg.ListTemplates()
	if len(templates) == 0 {
		fmt.Println("No templates defined.")
		fmt.Println("Edit your config file to add templates.")
		return nil
	}

	fmt.Println("Available templates:")
	fmt.Println()
	for name, tpl := range templates {
		desc := tpl.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  %-12s  %s (base: %s)\n", name, desc, tpl.Mode)
	}
	fmt.Println()
	fmt.Println("Usage: gocar new <name> --mode <template>")

	return nil
}

// editConfig 打开配置文件编辑
func (c *ConfigCommand) editConfig() error {
	if !config.GlobalConfigExists() {
		fmt.Println("No global config found.")
		fmt.Println("Run 'gocar config init' to create one first.")
		return nil
	}

	configPath, _ := config.GetGlobalConfigPath()
	fmt.Printf("Global config location: %s\n", configPath)
	fmt.Println("Please open this file in your preferred editor.")

	return nil
}

// Help 返回帮助信息
func (c *ConfigCommand) Help() string {
	return `gocar config - Manage global gocar configuration

USAGE:
    gocar config <SUBCOMMAND>

SUBCOMMANDS:
    init     Create global config file (~/.gocar/config.toml)
    path     Show global config file path
    list     List available project templates
    edit     Show config file location for editing

DESCRIPTION:
    The global config file allows you to:
    
    - Define custom project templates
    - Set default author and license
    - Create reusable project structures

    Templates can be used with: gocar new <name> --mode <template>
    
    Projects created from templates will automatically include a .gocar.toml
    configuration file with the template's settings.

EXAMPLES:
    gocar config init              Create global config with example templates
    gocar config list              List all available templates
    gocar config path              Show config file location
    gocar new myapi --mode api     Create project using 'api' template
`
}
