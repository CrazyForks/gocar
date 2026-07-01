package build

import "runtime"

// Config 构建配置
type Config struct {
	Release    bool   // 是否为发布模式
	Profile    string // 构建 profile 名称
	TargetOS   string // 目标操作系统
	TargetArch string // 目标架构
	WithCGO    bool   // 是否启用 CGO
}

// NewConfig 创建默认构建配置
func NewConfig() *Config {
	return &Config{
		Release:    false,
		Profile:    "debug",
		TargetOS:   runtime.GOOS,
		TargetArch: runtime.GOARCH,
		WithCGO:    false,
	}
}

// SetTarget 设置目标平台
func (c *Config) SetTarget(os, arch string) {
	c.TargetOS = os
	c.TargetArch = arch
}

// IsCurrentPlatform 检查是否为当前平台
func (c *Config) IsCurrentPlatform() bool {
	return c.TargetOS == runtime.GOOS && c.TargetArch == runtime.GOARCH
}

// BuildMode 返回构建模式字符串
func (c *Config) BuildMode() string {
	if c.Profile != "" {
		return c.Profile
	}
	if c.Release {
		return "release"
	}
	return "debug"
}
