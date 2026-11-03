# ftpx

基于 goftp 的 FTP 文件同步工具，提供批量文件上传能力。

## 功能列表

### sync.go

- `Sync`：FTP 同步客户端结构体。
  - `TimeOut int`：连接超时时间（秒）。
  - `Client *goftp.Client`：底层 goftp 客户端。
- `NewSync(host, port, user_name, pass_word string, time_out int) *Sync`
  - 初始化 FTP 同步客户端，通过 `goftp.DialConfig` 拨号连接。
  - 连接地址由 `net.JoinHostPort(host, port)` 拼接（host 与 port 均为字符串）。
  - 内部配置：每主机 10 个连接、指定秒数的超时、日志输出到 `os.Stderr`。
  - 拨号失败时记录日志并返回 nil，调用方需判空。
- `(*Sync) UploadFile(file_paths ...string) error`
  - 批量上传本地文件到 FTP 服务器。
  - 未传入任何文件路径时直接返回 nil。
  - 远程文件名取本地文件的 `filepath.Base`（即不保留本地目录结构）。
  - 完成后打印总耗时（依赖 logx）。

## 使用示例

```go
sync := ftpx.NewSync("ftp.example.com", "21", "user", "pass", 30)
if sync == nil {
    log.Fatal("FTP 连接失败")
}
if err := sync.UploadFile("/tmp/a.txt", "/tmp/b.txt"); err != nil {
    log.Fatal(err)
}
```

## 备注

- 外部依赖：`github.com/secsy/goftp`；另依赖本仓库 `github.com/qida/gohp/logx`。
- `UploadFile` 在打开文件或上传失败时会 panic，不适合对稳定性要求高的场景。
- 当前仅封装了上传（`Store`）操作，未提供下载、删除等能力。
- `Sync.Client` 为导出字段，可直接调用 `goftp.Client` 的原生方法扩展能力。
