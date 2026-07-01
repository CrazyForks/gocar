# gocar TODO

当前产品语义：`gocar new myapp` 创建一个标准 Go 应用，不提供模板选择。

默认结构：

```text
myapp/
├── cmd/myapp/main.go
├── internal/
├── go.mod
├── README.md
└── .gitignore
```

## 优先事项

1. 增强 `doctor`
   - 检查 `go env GOPATH/GOMODCACHE/GOPROXY`
   - 检查配置中的 `build.entry` 是否存在
   - 检查自定义命令是否覆盖内置命令

2. 增强 `check`
   - 支持 `--fix` 执行格式化
   - 支持跳过 `vet` 或指定包

3. 增强发布体验
   - 支持 `gocar build --profile release`
   - 支持常用目标矩阵构建

4. 自动化测试
   - 为 `new/build/run/check/doctor` 增加端到端测试
   - 在 CI 中运行 Linux/macOS/Windows matrix

## 非目标

- 不增加多模板系统。
