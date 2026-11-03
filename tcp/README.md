# tcp

TCP 服务端/客户端封装（fork 自 github.com/firstrow/tcp_server），支持回调式事件处理、TLS、消息中继与在线调试。

## 功能列表

### tcp_server.go — 服务端与客户端
- `New(address, relay) *server`：创建 TCP 服务端；`relay` 非空时收到的消息会转发到该地址
- `NewWithTLS(address, relay, certFile, keyFile)`：创建支持 TLS 的服务端
- `server.Listen()`：开始监听并 Accept（阻塞）；每个连接一个 goroutine，按 `\n` 分行读取
- `server.OnNewClient(callback)`：新客户端连接回调
- `server.OnNewMessage(callback)`：收到新消息回调
- `server.OnClientConnectionClosed(callback)`：连接关闭回调
- `Client.Send(message)` / `Client.SendBytes(b)`：向客户端发送文本/字节
- `Client.GetConn()` / `Client.Conn()` / `Client.Close()`：获取底层连接、关闭连接
- `server.RelayMessage`：中继消息 channel，`initRelay` 从中取消息转发到 relay 地址（断线后下次消息时重连）

### tcp_server_linux.go / tcp_server_windows.go — 平台实现
- `Listen()` 的平台相关实现：Linux 版对每个连接启用 TCP keepalive（空闲 30s、间隔 5s、探测 4 次，依赖 felixge/tcpkeepalive）；Windows 版不启用

### client_tcp_yd.go — 音频流客户端
- `TCP` / `NewTcpClient()`：流媒体 TCP 客户端，`Link(addr)` 建连并起 goroutine 把收到的数据送入 `ReceiveChan`
- `SendStart(stream_id)` / `SendStop()`：发送 startAudioStream/stopAudioStream 握手指令（JSON `Hand` 结构）
- `Close()`：关闭连接与 `ReceiveChan`

### logs.go — 调试服务
- `ServerDebug(port)`：启动调试 TCP 服务；客户端发送 `debug\r\n` 注册为调试端，其后所有消息广播给全部已注册调试端
- `DebugList`：已注册调试端列表（远端地址 → `*Client`）

## 使用示例

```go
server := tcp.New("localhost:9999", "")
server.OnNewClient(func(c *tcp.Client) { c.Send("Welcome\n") })
server.OnNewMessage(func(c *tcp.Client, message string) {
    fmt.Println("recv:", message)
})
server.OnClientConnectionClosed(func(c *tcp.Client, err error) {
    fmt.Println("closed:", err)
})
go server.Listen() // Listen 阻塞，需在 goroutine 中启动
```

## 备注

- 消息以 `\n` 为分隔符（`bufio.ReadString`），仅适用于按行协议
- TLS 模式证书加载失败仅打印日志，不阻断启动
