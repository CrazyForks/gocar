package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// ConfigFileName 配置文件名
const ConfigFileName = ".gocar.toml"

// GocarConfig gocar 配置结构
type GocarConfig struct {
	Project  ProjectConfig     `toml:"project"`
	Build    BuildConfig       `toml:"build"`
	Run      RunConfig         `toml:"run"`
	Profile  ProfilesConfig    `toml:"profile"`
	Commands map[string]string `toml:"commands"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name    string `toml:"name"`    // 项目名称，为空时使用目录名
	Version string `toml:"version"` // 项目版本号，构建时自动注入到 main.version
}

// ProfilesConfig 构建配置档案
type ProfilesConfig struct {
	Profiles map[string]ProfileConfig
}

// ProfileConfig 单个构建档案配置
type ProfileConfig struct {
	Ldflags    string `toml:"ldflags"`     // ldflags 参数
	Gcflags    string `toml:"gcflags"`     // 编译器参数
	Trimpath   *bool  `toml:"trimpath"`    // 是否移除路径信息
	CgoEnabled *bool  `toml:"cgo_enabled"` // 是否启用 CGO
	Race       bool   `toml:"race"`        // 是否启用竞态检测
}

// BuildConfig 构建配置
type BuildConfig struct {
	Entry    string   `toml:"entry"`     // 构建入口路径
	Output   string   `toml:"output"`    // 输出目录
	Ldflags  string   `toml:"ldflags"`   // 额外的 ldflags
	Tags     []string `toml:"tags"`      // 构建标签
	ExtraEnv []string `toml:"extra_env"` // 额外的环境变量
}

// RunConfig 运行配置
type RunConfig struct {
	Entry string   `toml:"entry"` // 运行入口路径
	Args  []string `toml:"args"`  // 默认运行参数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *GocarConfig {
	trueVal := true
	falseVal := false
	return &GocarConfig{
		Project: ProjectConfig{
			Name:    "",
			Version: "",
		},
		Build: BuildConfig{
			Entry:    "",
			Output:   "bin",
			Ldflags:  "",
			Tags:     []string{},
			ExtraEnv: []string{},
		},
		Run: RunConfig{
			Entry: "",
			Args:  []string{},
		},
		Profile: ProfilesConfig{
			Profiles: map[string]ProfileConfig{
				"debug": {
					Ldflags:    "",
					Gcflags:    "",
					Trimpath:   &falseVal,
					CgoEnabled: nil, // nil 表示跟随系统默认
					Race:       false,
				},
				"release": {
					Ldflags:    "-s -w",
					Gcflags:    "",
					Trimpath:   &trueVal,
					CgoEnabled: &falseVal,
					Race:       false,
				},
			},
		},
		Commands: map[string]string{
			"vet": "go vet ./...",
			"fmt": "go fmt ./...",
		},
	}
}

// DefaultConfigTemplate 返回默认配置文件模板
func DefaultConfigTemplate(projectName string) string {
	entry := "cmd/app"
	if projectName != "" {
		entry = "cmd/" + projectName
	}

	return fmt.Sprintf(`# gocar 项目配置文件
# 文档: https://github.com/uselibrary/gocar

# 项目配置
[project]
# 项目名称，留空则使用目录名
name = "%s"

# 项目版本号
# version = "1.0.0"

# 构建配置
[build]
# 构建入口路径 (相对于项目根目录)
# standard 布局默认为 "cmd/<appName>"
entry = "%s"

# 输出目录
output = "bin"

# 额外的 ldflags，会追加到 profile 的 ldflags 之后
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

# 自定义构建配置
# 使用: gocar build --profile ci
# [profile.ci]
# trimpath = true
# race = true

# 自定义命令
# 格式: 命令名 = "要执行的 shell 命令"
# 使用: gocar <命令名>
# 命令会在项目根目录下执行
#
# 自定义命令可以覆盖以下内置命令: build, run, clean, add, update, tidy, test, check, commands, doctor
# 保护命令 (new, init) 不可被覆盖
[commands]
# 代码检查
vet = "go vet ./..."

# 代码格式化
fmt = "go fmt ./..."

# lint = "golangci-lint run"
# doc = "godoc -http=:6060"
# proto = "protoc --go_out=. --go-grpc_out=. ./proto/*.proto"

# 覆盖内置命令示例 (取消注释以启用):
# build = "make build"
# run = "docker-compose up"
# clean = "make clean && rm -rf dist/"
`, projectName, entry)
}

// Load 从指定目录加载配置
func Load(projectRoot string) (*GocarConfig, error) {
	configPath := filepath.Join(projectRoot, ConfigFileName)

	// 使用内置默认配置作为基础
	baseConfig := DefaultConfig()

	// 如果项目配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return baseConfig, nil
	}

	// 加载项目配置
	projectConfig, err := decodeProjectConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", ConfigFileName, err)
	}

	// 合并: 项目配置覆盖默认配置
	finalConfig := mergeProjectConfig(baseConfig, projectConfig)

	return finalConfig, nil
}

func decodeProjectConfig(configPath string) (*GocarConfig, error) {
	var raw struct {
		Project  ProjectConfig            `toml:"project"`
		Build    BuildConfig              `toml:"build"`
		Run      RunConfig                `toml:"run"`
		Profile  map[string]ProfileConfig `toml:"profile"`
		Commands map[string]string        `toml:"commands"`
	}

	if _, err := toml.DecodeFile(configPath, &raw); err != nil {
		return nil, err
	}

	return &GocarConfig{
		Project:  raw.Project,
		Build:    raw.Build,
		Run:      raw.Run,
		Profile:  ProfilesConfig{Profiles: raw.Profile},
		Commands: raw.Commands,
	}, nil
}

// mergeProjectConfig 将项目配置合并到基础配置（项目配置优先）
func mergeProjectConfig(base *GocarConfig, project *GocarConfig) *GocarConfig {
	if project.Project.Name != "" {
		base.Project.Name = project.Project.Name
	}
	if project.Project.Version != "" {
		base.Project.Version = project.Project.Version
	}

	// Build 配置
	if project.Build.Entry != "" {
		base.Build.Entry = project.Build.Entry
	}
	if project.Build.Output != "" {
		base.Build.Output = project.Build.Output
	}
	if project.Build.Ldflags != "" {
		base.Build.Ldflags = project.Build.Ldflags
	}
	if len(project.Build.Tags) > 0 {
		base.Build.Tags = project.Build.Tags
	}
	if len(project.Build.ExtraEnv) > 0 {
		base.Build.ExtraEnv = project.Build.ExtraEnv
	}

	// Run 配置
	if project.Run.Entry != "" {
		base.Run.Entry = project.Run.Entry
	}
	if len(project.Run.Args) > 0 {
		base.Run.Args = project.Run.Args
	}

	// Profile 配置
	for name, profile := range project.Profile.Profiles {
		if strings.TrimSpace(name) == "" {
			continue
		}
		base.Profile.Profiles[name] = mergeProfile(base.Profile.Profiles[name], profile)
	}

	// Commands - 项目命令覆盖全局命令
	for name, cmd := range project.Commands {
		base.Commands[name] = cmd
	}

	return base
}

func mergeProfile(base ProfileConfig, project ProfileConfig) ProfileConfig {
	if project.Ldflags != "" {
		base.Ldflags = project.Ldflags
	}
	if project.Gcflags != "" {
		base.Gcflags = project.Gcflags
	}
	if project.Trimpath != nil {
		base.Trimpath = project.Trimpath
	}
	if project.CgoEnabled != nil {
		base.CgoEnabled = project.CgoEnabled
	}
	if project.Race {
		base.Race = project.Race
	}
	return base
}

// Exists 检查配置文件是否存在
func Exists(projectRoot string) bool {
	configPath := filepath.Join(projectRoot, ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// Save 保存配置到文件
func Save(projectRoot, projectName string) error {
	configPath := filepath.Join(projectRoot, ConfigFileName)
	content := DefaultConfigTemplate(projectName)
	return os.WriteFile(configPath, []byte(content), 0644)
}

// GetBuildEntry 获取构建入口路径
func (c *GocarConfig) GetBuildEntry() string {
	return c.GetBuildEntryForApp("")
}

// GetBuildEntryForApp 获取构建入口路径，可用 defaultName 补全标准布局入口。
func (c *GocarConfig) GetBuildEntryForApp(defaultName string) string {
	if c.Build.Entry != "" {
		return c.Build.Entry
	}

	if c.Project.Name != "" {
		return "cmd/" + c.Project.Name
	}
	if defaultName != "" {
		return "cmd/" + defaultName
	}
	return "cmd/app"
}

// GetRunEntry 获取运行入口路径
func (c *GocarConfig) GetRunEntry() string {
	return c.GetRunEntryForApp("")
}

// GetRunEntryForApp 获取运行入口路径，可用 defaultName 补全标准布局入口。
func (c *GocarConfig) GetRunEntryForApp(defaultName string) string {
	if c.Run.Entry != "" {
		return c.Run.Entry
	}
	return c.GetBuildEntryForApp(defaultName)
}

// GetBuildOutputRoot 获取构建输出根目录
func (c *GocarConfig) GetBuildOutputRoot() string {
	output := strings.TrimSpace(c.Build.Output)
	if output == "" {
		return "bin"
	}
	return filepath.Clean(output)
}

// ResolveBuildOutputDir 解析并校验构建输出目录（用于清理等破坏性操作）
func (c *GocarConfig) ResolveBuildOutputDir(projectRoot string) (string, error) {
	outputRoot := c.GetBuildOutputRoot()
	if outputRoot == "." {
		return "", fmt.Errorf("invalid [build].output: '.' is not allowed")
	}

	buildOutputDir := outputRoot
	if !filepath.IsAbs(buildOutputDir) {
		buildOutputDir = filepath.Join(projectRoot, buildOutputDir)
	}

	absOutputDir, err := filepath.Abs(buildOutputDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve build output directory: %w", err)
	}
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root: %w", err)
	}

	if filepath.Dir(absOutputDir) == absOutputDir {
		return "", fmt.Errorf("invalid [build].output: filesystem root is not allowed")
	}
	if absOutputDir == absProjectRoot {
		return "", fmt.Errorf("invalid [build].output: project root is not allowed")
	}

	return absOutputDir, nil
}

// GetProjectName 获取项目名称
func (c *GocarConfig) GetProjectName(defaultName string) string {
	if c.Project.Name != "" {
		return c.Project.Name
	}
	return defaultName
}

// GetCommand 获取自定义命令
func (c *GocarConfig) GetCommand(name string) (string, bool) {
	cmd, ok := c.Commands[name]
	return cmd, ok
}

// RunCustomCommand 执行自定义命令
func (c *GocarConfig) RunCustomCommand(projectRoot, name string, extraArgs []string) error {
	cmdStr, ok := c.Commands[name]
	if !ok {
		return fmt.Errorf("command '%s' not defined in %s", name, ConfigFileName)
	}

	// 如果有额外参数，追加到命令后面
	if len(extraArgs) > 0 {
		quotedArgs := make([]string, 0, len(extraArgs))
		for _, arg := range extraArgs {
			quotedArgs = append(quotedArgs, shellQuoteArg(arg))
		}
		cmdStr = cmdStr + " " + strings.Join(quotedArgs, " ")
	}

	// 使用 shell 执行命令
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func shellQuoteArg(arg string) string {
	if arg == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
}

// ListCommands 列出所有自定义命令
func (c *GocarConfig) ListCommands() map[string]string {
	return c.Commands
}

// GetProfile 获取指定名称的构建配置
func (c *GocarConfig) GetProfile(name string) (*ProfileConfig, bool) {
	if name == "" {
		name = "debug"
	}
	profile, ok := c.Profile.Profiles[name]
	if !ok {
		return nil, false
	}
	return &profile, true
}

// GetProfileForBuild 获取当前构建使用的 profile。
func (c *GocarConfig) GetProfileForBuild(profileName string, release bool) (*ProfileConfig, string, bool) {
	if profileName == "" {
		if release {
			profileName = "release"
		} else {
			profileName = "debug"
		}
	}
	profile, ok := c.GetProfile(profileName)
	return profile, profileName, ok
}

// ListProfiles 列出所有 profile 名称。
func (c *GocarConfig) ListProfiles() []string {
	names := make([]string, 0, len(c.Profile.Profiles))
	for name := range c.Profile.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Validate 校验配置中的关键约束。
func (c *GocarConfig) Validate(projectRoot string) error {
	if _, err := c.ResolveBuildOutputDir(projectRoot); err != nil {
		return err
	}
	if len(c.Profile.Profiles) == 0 {
		return fmt.Errorf("at least one build profile is required")
	}
	for name := range c.Profile.Profiles {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("profile name cannot be empty")
		}
	}
	for name, cmd := range c.Commands {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("custom command name cannot be empty")
		}
		if strings.TrimSpace(cmd) == "" {
			return fmt.Errorf("custom command %q cannot be empty", name)
		}
	}

	return nil
}
