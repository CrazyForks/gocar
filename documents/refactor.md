# gocar Refactor Notes

当前结构已经接近目标分层：

```text
gocar/
├── cmd/gocar/main.go
├── internal/
│   ├── build/
│   ├── cli/
│   ├── config/
│   ├── project/
│   └── util/
├── README.md
└── go.mod
```

## 保持的设计约束

- `gocar new myapp` 是唯一的新建项目主路径。
- 新项目始终使用 `cmd/<appName>/main.go`。
- `internal/` 默认创建，构建产物目录只在构建时创建。
- `.gocar.toml` 只描述项目名称、构建、运行、profile 和自定义命令，不描述模板或模式。

## 后续可重构点

1. CLI 声明集中化
   - 统一命令名称、help、示例和保护命令列表。

2. 配置校验集中化
   - 为 `build.entry`、`run.entry`、`build.output`、profile 名称提供更细的错误信息。

3. 子进程执行统一化
   - 统一 `go build`、`go run`、自定义命令和依赖命令的输出、退出码、环境变量处理。
