# osx

系统信息与 SSH/网络工具集，提供 SSH 远程命令执行、交互式终端、网卡上传速率统计以及 CPU/主板 ID 获取功能。

## 功能列表

### ssh.go — SSH 客户端

- `Cli`：SSH 客户端结构，导出字段为 IP、Username、Password、Port、LastResult。
- `NewSSH(ip, username, password string, port ...int) *Cli`：创建 SSH 客户端对象，默认端口 22。
- `(*Cli) RunSH(cmd string) error`：自动连接并远程执行命令，打印合并输出；执行失败时调用 `log.Fatal` 直接退出进程。
- `(*Cli) RunTerminal(stdout, stderr io.Writer) error`：申请 xterm-256color PTY 并接管本地 stdin，进入交互式远程终端（退出时恢复终端原状态）。

### network.go — 网速统计

- `NetSpeed(speed chan float64)`：每秒采样一次各网卡 IO 计数，计算并打印每个网卡的上传速率（KB/s），将所有网卡的上传总速率（字节/秒）写入 channel；需在 goroutine 中运行。

### os_linux.go — 硬件信息（Linux）

- `GetCPUID() string`：通过 `sudo dmidecode -t 4` 命令获取 CPU ID。
- `GetBaseBoardID() string`：获取主板 ID（当前实现返回空字符串）。

### os_windows.go — 硬件信息（Windows）

- `GetCPUID() string`：通过 `wmic cpu get processorid` 获取 CPU ID，失败返回 "unkonw"。
- `GetBaseBoardID() string`：通过 `wmic baseboard get serialnumber` 获取主板序列号。

## 使用示例

远程执行命令并进入交互式终端：

```go
cli := osx.NewSSH("192.168.1.10", "root", "password", 22)
if err := cli.RunSH("uptime"); err != nil {
    log.Fatal(err)
}
// 交互式终端（需要在真实 TTY 环境下运行）
// cli.RunTerminal(os.Stdout, os.Stderr)
```

## 备注

- 外部依赖：`golang.org/x/crypto/ssh`、`golang.org/x/term`、`github.com/shirou/gopsutil/v3`。
- SSH 连接的 HostKeyCallback 直接返回 nil（不校验主机公钥），仅适合可信内网环境。
- 硬件信息函数通过构建标签分别提供 linux / windows 实现。
- Linux 下 `GetCPUID` 依赖 `sudo dmidecode`（需要 sudo/root 权限）；Windows 下依赖系统自带 `wmic` 命令。
- `RunTerminal` 会将本地终端设为 raw 模式，只适合命令行程序中使用。
