package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// GlobalConfigFileName 全局配置文件名
const GlobalConfigFileName = "config.toml"

// GlobalConfigDir 全局配置目录名
const GlobalConfigDir = ".gocar"

// GlobalConfig 全局配置结构
type GlobalConfig struct {
	Templates map[string]TemplateConfig `toml:"templates"`
	Defaults  DefaultsConfig            `toml:"defaults"`
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Description string            `toml:"description"` // 模板描述
	Mode        string            `toml:"mode"`        // 基础模式: simple 或 project
	Dirs        []string          `toml:"dirs"`        // 额外创建的目录
	Files       map[string]string `toml:"files"`       // 额外创建的文件 (路径 -> 内容)
	Commands    map[string]string `toml:"commands"`    // 预设的自定义命令
	Build       BuildConfig       `toml:"build"`       // 构建配置
	Run         RunConfig         `toml:"run"`         // 运行配置
}

// DefaultsConfig 默认配置
type DefaultsConfig struct {
	Author  string `toml:"author"`  // 默认作者
	License string `toml:"license"` // 默认许可证
}

// GetGlobalConfigDir 获取全局配置目录路径
func GetGlobalConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, GlobalConfigDir), nil
}

// GetGlobalConfigPath 获取全局配置文件路径
func GetGlobalConfigPath() (string, error) {
	configDir, err := GetGlobalConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, GlobalConfigFileName), nil
}

// GlobalConfigExists 检查全局配置文件是否存在
func GlobalConfigExists() bool {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(configPath)
	return err == nil
}

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig() (*GlobalConfig, error) {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultGlobalConfig(), nil
	}

	config := &GlobalConfig{
		Templates: make(map[string]TemplateConfig),
	}
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", configPath, err)
	}

	return config, nil
}

// DefaultGlobalConfig 返回默认全局配置
func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Templates: make(map[string]TemplateConfig),
		Defaults: DefaultsConfig{
			Author:  "",
			License: "MIT",
		},
	}
}

// SaveGlobalConfig 保存全局配置
func SaveGlobalConfig() error {
	configDir, err := GetGlobalConfigDir()
	if err != nil {
		return err
	}

	// 创建配置目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, GlobalConfigFileName)
	content := GlobalConfigTemplate()
	return os.WriteFile(configPath, []byte(content), 0644)
}

// GlobalConfigTemplate 返回全局配置文件模板
func GlobalConfigTemplate() string {
	return `# gocar 全局配置文件
# 位置: ~/.gocar/config.toml
# 文档: https://github.com/uselibrary/gocar

# 默认设置
[defaults]
# 默认作者
author = ""

# 默认许可证
license = "MIT"

# 项目模板
# 使用方式: gocar new <name> --mode <template_name>
# 
# 模板会继承基础模式 (simple/project) 的结构，并添加额外的目录和文件
# 使用模板创建的项目会自动包含 .gocar.toml 配置文件

# 示例: Web API 模板
[templates.api]
description = "Web API project with common structure"
mode = "project"  # 基础模式: simple 或 project

# 额外创建的目录
dirs = [
    "api",
    "configs",
    "scripts",
]

# 预设的自定义命令
[templates.api.commands]
dev = "go run cmd/server/main.go -env=dev"
lint = "golangci-lint run ./..."

# 示例: CLI 工具模板
[templates.cli]
description = "CLI tool project"
mode = "simple"

dirs = [
    "cmd",
]

[templates.cli.commands]
install = "go install ."

# 示例: 库模板
[templates.lib]
description = "Go library project"
mode = "simple"

dirs = [
    "examples",
]

[templates.lib.commands]
test = "go test -v -cover ./..."
bench = "go test -bench=. ./..."
`
}

// GetTemplate 获取指定模板
func (c *GlobalConfig) GetTemplate(name string) (*TemplateConfig, bool) {
	tpl, ok := c.Templates[name]
	if !ok {
		return nil, false
	}
	return &tpl, true
}

// ListTemplates 列出所有模板
func (c *GlobalConfig) ListTemplates() map[string]TemplateConfig {
	return c.Templates
}
