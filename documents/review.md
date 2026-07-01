# gocar Review Checklist

## 新建项目

- `gocar new myapp` 应生成：

```text
myapp/
├── cmd/myapp/main.go
├── internal/
├── go.mod
├── README.md
└── .gitignore
```

- 不应生成额外的顶层源码布局目录。
- 额外参数应返回清晰错误。

## 构建与运行

- 新项目中 `gocar build` 应构建 `cmd/<appName>`。
- 新项目中 `gocar run` 应运行 `cmd/<appName>`。
- 构建产物应输出到 `bin/<profile>/<os>-<arch>/<appName>`。

## 配置

- `.gocar.toml` 不应包含模板或模式字段。
- `[project].name` 可覆盖应用名。
- `[build].entry` 可覆盖构建入口。
- `[run].entry` 可覆盖运行入口。
- `[profile.<name>]` 可定义构建 profile。

## 命令

- `gocar commands` 应列出内置命令和自定义命令。
- `gocar doctor` 应检查 Go、Git、项目检测和配置合法性。
- `gocar check` 应执行 `go fmt ./...`、`go vet ./...`、`go test ./...`。

## 本地验证

```bash
GOCACHE=/tmp/gocar-go-build go test ./...
GOCACHE=/tmp/gocar-go-build go run ./cmd/gocar doctor
```
