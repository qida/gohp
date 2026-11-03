# websocket

WebSocket 连接管理中枢：按网络地址与设备订阅两个维度管理客户端，支持定向广播与级联转发（连接上游 WebSocket 并把其消息转发给本端订阅者）。

## 功能列表

### websocket.go — 连接管理与广播
- `ClientMangerCenter` / `NewWsClientMangerCenter()`：管理中枢（全局单例）
- `Start()`：事件循环——处理 Register（注册双 map、回写 "OK"、可选建立上游 Link）、Unregister（注销并关闭连接）、Broadcaster（按订阅 Key 定向推送，3s 超时踢掉慢客户端）
- `Broadcast(sub_device, message)`：向指定订阅及全体订阅（`all_-1`）各广播一次
- `ClientWS` / `NewClientWS(conn, sub_device)`：客户端连接封装，创建后自动启动读循环与两个写循环（`MsgConn`→本端、`MsgLink`→上游）
- `ClientWS.StartLink(t)`：连接 `Url` 指向的上游 WebSocket（失败 20s 重试），上游消息经 `Broadcast` 转发给订阅者
- `SubDevice`（DeviceId/IdScheme/ZoneId）与 `GetKeySub()`：订阅维度，KeySub 格式为 `设备ID_方案ID`
- `ListSub`/`CountSub`/`ListNet`/`CountNet`：列出/统计订阅与在线连接
- `ClientWS.Close()`：把自身送入 Unregister 注销

### websocket_in.go — 并发安全映射
- `MapWSNet`：KeyNet（远端地址）→ 客户端
- `MapWSSub`：KeySub（订阅键）→ KeyNet → 客户端
- `MapLink`、`MapWSAll`：两级映射结构（当前未在管理中枢中使用）
- 以上类型均带 `sync.RWMutex`，提供 `Add` / `Delete`

## 使用示例

```go
go websocket.NewWsClientMangerCenter().Start() // 启动事件循环

// HTTP 升级得到 conn 后注册客户端
client := websocket.NewClientWS(conn, websocket.SubDevice{DeviceId: "dev1", IdScheme: 1})
websocket.NewWsClientMangerCenter().Register <- client

// 向所有订阅 dev1_1 的客户端广播
websocket.NewWsClientMangerCenter().Broadcast(
    websocket.SubDevice{DeviceId: "dev1", IdScheme: 1}, []byte(`{"msg":"hello"}`))
```

## 备注

- 依赖 github.com/gorilla/websocket
- 使用前必须先调用 `Start()`，否则 Register/Unregister/Broadcast 会阻塞（channel 无缓冲）
- 客户端若需级联上游，注册前设置 `ClientWS.Url` 字段
