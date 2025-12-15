# gocar

一个"类 Cargo"的 Go 项目脚手架与命令行工具，提供简洁的项目初始化和构建体验。

## 安装

```bash
go install github.com/yourname/gocar@latest
```

或从源码构建：

```bash
git clone https://github.com/yourname/gocar.git
cd gocar
go build -o gocar main.go
```

## 快速开始

```bash
# 创建新项目（简洁模式）
gocar new myapp

# 进入项目目录
cd myapp

# 构建项目
gocar build

# 运行项目
gocar run
```

## 命令

### `gocar new <name> [--mode simple|project]`

创建新的 Go 项目。

**参数：**
- `<name>` - 项目名称，同时作为目录名和输出的可执行文件名
- `--mode` - 项目模式，可选 `simple`（默认）或 `project`

**示例：**
```bash
# 创建简洁模式项目（默认）
gocar new myapp

# 创建项目模式项目
gocar new myserver --mode project
```

### `gocar build [--release]`

构建当前项目。

**参数：**
- `--release` - 使用 release 模式构建（优化体积）

**构建行为：**

| 模式 | 命令等价 |
|------|----------|
| Debug（默认） | `go build -o bin/<appName> ./main.go` |
| Release | `CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/<appName> ./main.go` |

> 项目模式下入口为 `./cmd/server/main.go`

**示例：**
```bash
# Debug 构建
gocar build

# Release 构建（更小的二进制文件）
gocar build --release
```

### `gocar run [args...]`

构建并运行当前项目。

**示例：**
```bash
# 运行项目
gocar run

# 传递参数给应用
gocar run --port 8080
```

### `gocar help`

显示帮助信息。

### `gocar version`

显示版本信息。

---

## 项目模式

### 简洁模式（Simple）

适用于小型项目、脚本、CLI 工具等。

```
myapp/
├── go.mod
├── main.go
├── README.md
├── bin/
├── .gitignore
└── .git/
```

**main.go 模板：**
```go
package main

import "fmt"

func main() {
    fmt.Println("hello from gocar")
}
```

### 项目模式（Project）

适用于大型项目、Web 服务、微服务等，遵循 Go 标准项目布局。

```
myapp/
├── cmd/
│   └── server/
│       └── main.go      # 应用入口
├── internal/            # 私有代码（不可被外部导入）
├── pkg/                 # 公共库代码
├── test/                # 集成测试
├── bin/                 # 构建输出
├── go.mod
├── README.md
├── .gitignore
└── .git/
```

**cmd/server/main.go 模板：**
```go
package main

import "fmt"

func main() {
    fmt.Println("server started")
}
```

**目录说明：**
- `cmd/` - 应用程序入口
- `internal/` - 私有代码，仅本模块内部使用（Go 编译器强制）
- `pkg/` - 可被外部项目导入的公共库
- `test/` - 集成测试、端到端测试

---

## 特性

- ✅ **自动 Git 初始化** - 创建项目时自动执行 `git init` 并生成 `.gitignore`
- ✅ **智能项目检测** - 自动识别 simple/project 模式
- ✅ **Release 优化** - 使用 `-ldflags="-s -w" -trimpath` 减小二进制体积
- ✅ **跨平台支持** - Windows 自动添加 `.exe` 后缀

---

## .gitignore 模板

自动生成的 `.gitignore` 包含：

```gitignore
# Binaries
<appName>
bin/
*.exe
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/

# OS files
.DS_Store
Thumbs.db
```

---

## License

MIT License
