# cmdx

系统命令执行与彩色日志输出工具，提供带 stdout/stderr 分离返回的命令执行封装。

## 功能列表

### cmd.go — 命令执行

- `ExecCmd(cmdName string, args ...string) (string, string, error)`：执行系统命令，返回 stdout、stderr（字符串）及错误。
- `ExecCmdDir(dir, cmdName string, args ...string)`：在指定目录下执行系统命令，返回字符串形式的 stdout/stderr。
- `ExecCmdBytes(cmdName string, args ...string)`：执行系统命令，返回字节形式的 stdout/stderr。
- `ExecCmdDirBytes(dir, cmdName string, args ...string)`：在指定目录下执行系统命令，返回字节形式的输出。

### cmd.go — 彩色日志

- 颜色常量：`Gray`、`Red`、`Green`、`Yellow`、`Blue`、`Magenta`（ANSI 颜色码）与 `EndColor`（颜色结束符）。
- `ColorLog(format string, a ...interface{})`：向 stdout 打印彩色日志。
- `ColorLogS(format string, a ...interface{}) string`：生成彩色日志字符串，着色规则：
  - 级别 `[TRAC]`/`[ERRO]`/`[WARN]`/`[SUCC]` 分别着蓝色/红色/品红/绿色，其他级别原样输出（级别须位于行首）。
  - 高亮内容用 `# 内容 #` 包裹（灰色）。
  - 路径用 `( 路径 )` 包裹（黄色）。
  - 错误内容用 `[ 错误 ]` 包裹（红色）。
  - Windows 下自动剥离着色标记，仅保留纯文本。

## 使用示例

```go
// 执行系统命令（提炼自 cmd_test.go）
stdout, stderr, err := cmdx.ExecCmd("go", "help", "get")
fmt.Println(stdout, stderr, err)

// 打印彩色日志
cmdx.ColorLog("[SUCC] build ok, output( %s )", "app.exe")
```

## 备注

- 包使用 `go1.3` 构建标签，要求 Go 1.3 及以上版本。
- 彩色输出依赖 ANSI 转义序列，Windows 平台自动降级为纯文本（仅去除标记）。
- 各命令执行函数基于 `os/exec`，命令不存在或退出码非 0 时通过 error 返回。
