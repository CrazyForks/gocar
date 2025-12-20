package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gocar/internal/config"
	"gocar/internal/util"
)

// Creator 项目创建器
type Creator struct {
	Name     string
	Mode     string
	Template *config.TemplateConfig // 模板配置（可选）
}

// NewCreator 创建项目创建器
func NewCreator(name, mode string) *Creator {
	return &Creator{
		Name:     name,
		Mode:     mode,
		Template: nil,
	}
}

// NewCreatorWithTemplate 使用模板创建项目创建器
func NewCreatorWithTemplate(name string, tpl *config.TemplateConfig) *Creator {
	return &Creator{
		Name:     name,
		Mode:     tpl.Mode,
		Template: tpl,
	}
}

// Create 创建项目
func (c *Creator) Create() error {
	// Check if directory already exists
	if _, err := os.Stat(c.Name); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", c.Name)
	}

	var err error
	if c.Template != nil {
		// 使用模板创建
		err = c.createFromTemplate()
	} else if c.Mode == "simple" {
		err = c.createSimpleProject()
	} else {
		err = c.createProjectMode()
	}

	if err != nil {
		return err
	}

	// Initialize git
	if err := util.InitGit(c.Name); err != nil {
		fmt.Printf("Warning: Failed to initialize git: %v\n", err)
	}

	return nil
}

// createSimpleProject 创建简单项目
func (c *Creator) createSimpleProject() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create main.go
	if err := util.WriteFile(filepath.Join(c.Name, "main.go"), SimpleMainTemplate(c.Name)); err != nil {
		return err
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), SimpleReadmeTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), GitignoreTemplate(c.Name)); err != nil {
		return err
	}

	return nil
}

// createProjectMode 创建项目模式
func (c *Creator) createProjectMode() error {
	// Create directories
	dirs := []string{
		c.Name,
		filepath.Join(c.Name, "cmd", "server"),
		filepath.Join(c.Name, "internal"),
		filepath.Join(c.Name, "pkg"),
		filepath.Join(c.Name, "test"),
		filepath.Join(c.Name, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	if err := util.RunCommandSilent(c.Name, "go", "mod", "init", c.Name); err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	// Create cmd/server/main.go
	if err := util.WriteFile(filepath.Join(c.Name, "cmd", "server", "main.go"), ProjectMainTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitkeep files for empty directories
	emptyDirs := []string{
		filepath.Join(c.Name, "internal", ".gitkeep"),
		filepath.Join(c.Name, "pkg", ".gitkeep"),
		filepath.Join(c.Name, "test", ".gitkeep"),
	}
	for _, f := range emptyDirs {
		if err := util.WriteFile(f, ""); err != nil {
			return err
		}
	}

	// Create README.md
	if err := util.WriteFile(filepath.Join(c.Name, "README.md"), ProjectReadmeTemplate(c.Name)); err != nil {
		return err
	}

	// Create .gitignore
	if err := util.WriteFile(filepath.Join(c.Name, ".gitignore"), GitignoreTemplate(c.Name)); err != nil {
		return err
	}

	return nil
}

// createFromTemplate 从模板创建项目
func (c *Creator) createFromTemplate() error {
	// 首先创建基础项目结构
	var err error
	if c.Mode == "simple" {
		err = c.createSimpleProject()
	} else {
		err = c.createProjectMode()
	}

	if err != nil {
		return err
	}

	// 创建模板中定义的额外目录
	for _, dir := range c.Template.Dirs {
		dirPath := filepath.Join(c.Name, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		// 创建 .gitkeep 文件
		if err := util.WriteFile(filepath.Join(dirPath, ".gitkeep"), ""); err != nil {
			return err
		}
	}

	// 创建模板中定义的额外文件
	for filePath, content := range c.Template.Files {
		fullPath := filepath.Join(c.Name, filePath)
		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
		}
		if err := util.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	// 创建 .gocar.toml 配置文件（模板项目自动包含配置文件）
	if err := c.createTemplateConfig(); err != nil {
		return err
	}

	return nil
}

// createTemplateConfig 创建模板项目的配置文件
func (c *Creator) createTemplateConfig() error {
	// 合并模板中的命令与默认命令
	commands := map[string]string{
		"vet":  "go vet ./...",
		"fmt":  "go fmt ./...",
		"test": "go test -v ./...",
	}

	// 添加模板中定义的命令
	for name, cmd := range c.Template.Commands {
		commands[name] = cmd
	}

	// 生成配置文件内容
	content := c.generateTemplateConfigContent(commands)

	return util.WriteFile(filepath.Join(c.Name, config.ConfigFileName), content)
}

// generateTemplateConfigContent 生成模板配置文件内容
func (c *Creator) generateTemplateConfigContent(commands map[string]string) string {
	entry := "."
	if c.Mode == "project" {
		entry = "cmd/server"
	}

	// 构建命令部分
	cmdSection := ""
	for name, cmd := range commands {
		cmdSection += fmt.Sprintf("%s = %q\n", name, cmd)
	}

	return fmt.Sprintf(`# gocar 项目配置文件
# 文档: https://github.com/uselibrary/gocar

# 项目配置
[project]
# 项目模式: "simple" (单文件) 或 "project" (标准目录结构)
# 留空则自动检测
mode = "%s"

# 项目名称，留空则使用目录名
name = "%s"

# 构建配置
[build]
# 构建入口路径 (相对于项目根目录)
# simple 模式默认为 ".", project 模式默认为 "cmd/server"
entry = "%s"

# 输出目录
output = "bin"

# 额外的 ldflags，会追加到默认 ldflags 之后
# 例如: "-X main.version=1.0.0"
ldflags = ""

# 构建标签
# tags = ["jsoniter", "sonic"]

# 额外的环境变量
# extra_env = ["GOPROXY=https://goproxy.cn"]

# 运行配置
[run]
# 运行入口路径，留空则使用 build.entry
entry = ""

# 默认运行参数
# args = ["-config", "config.yaml"]

# Debug 构建配置
# 使用: gocar build (默认)
[profile.debug]
# ldflags = ""              # Debug 默认无 ldflags
# gcflags = "all=-N -l"     # 禁用优化，方便调试
# trimpath = false          # 保留路径信息
# cgo_enabled = true        # 跟随系统默认
# race = false              # 竞态检测 (会显著降低性能)

# Release 构建配置
# 使用: gocar build --release
[profile.release]
ldflags = "-s -w"           # 裁剪符号表和调试信息
# gcflags = ""              # 编译器参数
trimpath = true             # 移除编译路径信息
cgo_enabled = false         # 禁用 CGO 以生成静态二进制
# race = false              # 竞态检测

# 自定义命令
# 格式: 命令名 = "要执行的 shell 命令"
# 使用: gocar <命令名>
# 命令会在项目根目录下执行
[commands]
%s`, c.Mode, c.Name, entry, cmdSection)
}
