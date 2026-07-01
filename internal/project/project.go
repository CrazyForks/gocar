package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// Info 项目信息
type Info struct {
	Root string // 项目根目录
	Name string // 项目名称
	Mode string // 项目布局: "standard"
}

// Detector 项目检测器
type Detector struct{}

// NewDetector 创建项目检测器
func NewDetector() *Detector {
	return &Detector{}
}

// Detect 检测项目信息
func (d *Detector) Detect() (*Info, error) {
	// 查找项目根目录
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	root := cwd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			return nil, fmt.Errorf("not in a Go module (go.mod not found)")
		}
		root = parent
	}

	// 检测项目布局
	mode := d.detectMode(root)
	if mode == "" {
		return nil, fmt.Errorf("cannot detect project layout: cmd/*/main.go does not exist")
	}

	return &Info{
		Root: root,
		Name: filepath.Base(root),
		Mode: mode,
	}, nil
}

// detectMode 检测项目布局
func (d *Detector) detectMode(root string) string {
	// Standard mode: any cmd/*/main.go exists.
	cmdGlob := filepath.Join(root, "cmd", "*", "main.go")
	matches, err := filepath.Glob(cmdGlob)
	if err == nil && len(matches) > 0 {
		return ModeStandard
	}

	return ""
}

// DetectProject 便捷函数：检测当前项目
func DetectProject() (projectRoot, appName, projectMode string, err error) {
	detector := NewDetector()
	info, err := detector.Detect()
	if err != nil {
		return "", "", "", err
	}
	return info.Root, info.Name, info.Mode, nil
}
