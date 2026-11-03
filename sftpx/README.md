# sftpx

基于 SSH 的 SFTP 客户端封装（派生自 melbahja/goph），提供密码/私钥/agent 多种认证、远程命令执行、文件上传下载（含进度回调）以及 known_hosts 管理。

## 功能列表

### sftpx.go — 客户端与文件传输

- `Client`：SSH 客户端封装，内嵌 `*ssh.Client` 并携带 `Config`。
- `Config`：连接配置（Auth、User、Addr、Port、Timeout、HostKeyCallback、BannerCallback）。
- `DefaultTimeout`：默认连接超时（20 秒）。
- `New(user, addr, auth)`：建立 SSH 连接，远程主机公钥必须在 known_hosts 中。
- `NewUnknown(user, addr, auth)`：跳过 known_hosts 校验建立连接（不安全，慎用）。
- `NewConn(config *Config)` / `Dial(proto, config)`：按自定义配置建立连接。
- `(*Client) Run(cmd string)` / `RunContext(ctx, name)`：远程执行命令，返回合并输出，后者支持 context 取消。
- `(*Client) Command / CommandContext`：构造类似 `os/exec` 的远程命令对象（可设置环境变量）。
- `(*Client) NewSftp(opts...)`：基于当前 SSH 连接创建 `*sftp.Client`。
- `(*Client) Upload(localPath, remotePath)`：上传本地文件到远程服务器。
- `(*Client) UploadWithProgress(localPath, remotePath, progress chan float64)`：按 2MB 分块上传，通过 channel 持续推送进度百分比（0-100），出错时发送 -1，正常结束后关闭 channel。
- `(*Client) Download(remotePath, localPath)`：下载远程文件到本地。
- `(*Client) Close()`：关闭连接。

### auth.go — 认证方式

- `Auth`：SSH 认证方法集合（`[]ssh.AuthMethod`）。
- `Password(pass)`：密码认证。
- `KeyboardInteractive(pass)`：键盘交互认证（密码作为应答，兼作密码认证回退）。
- `Key(prvFile, passphrase)`：私钥文件认证（支持带口令私钥）。
- `RawKey(privateKey, passphrase)`：私钥内容（字符串）认证。
- `HasAgent()` / `UseAgent()`：检测并使用 ssh-agent 认证（仅 Unix 系统）。
- `GetSigner / GetSignerForRawKey`：解析私钥获取 `ssh.Signer`。

### cmd.go — 远程命令对象

- `Cmd`：类似 `os/exec.Cmd` 的远程命令结构，包含 Path、Args、Env、Session、Context。
- `(*Cmd) CombinedOutput / Output / Run / Start / String`：执行远程命令的各种方式；context 取消时向远端发送 SIGINT。

### hosts.go — known_hosts 管理

- `DefaultKnownHosts()` / `KnownHosts(file)`：生成 known_hosts 校验回调。
- `CheckKnownHost(host, remote, key, knownFile)`：检查主机是否在 known_hosts 中；已存在但公钥不匹配时返回错误（可能存在中间人攻击）。
- `AddKnownHost(host, remote, key, knownFile)`：将主机公钥追加写入 known_hosts 文件。
- `DefaultKnownHostsPath()`：返回默认路径 `~/.ssh/known_hosts`。

## 使用示例

```go
auth := sftpx.Password("your-password")
client, err := sftpx.New("user", "example.com", auth)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 上传文件
if err := client.Upload("/tmp/local.txt", "/remote/path/local.txt"); err != nil {
    log.Fatal(err)
}

// 远程执行命令
out, err := client.Run("ls -l /remote/path")
fmt.Println(string(out))
```

## 备注

- 外部依赖：`github.com/pkg/sftp`、`golang.org/x/crypto/ssh`、`github.com/pkg/errors`。
- `New()` 要求远程主机公钥已存在于 known_hosts，否则连接失败；可先用 `AddKnownHost` 添加。
- `UseAgent()` 仅支持 Unix 系统（依赖 `SSH_AUTH_SOCK` 环境变量）。
