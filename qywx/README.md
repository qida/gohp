# qywx

企业微信 API 工具包，负责获取与自动刷新企业微信应用的 AccessToken，并定义应用消息相关结构体。

## 功能列表

### qywx.go

- `QYWX`：企业微信客户端结构体，字段包括：
  - `AccessToken`：当前有效的访问令牌。
  - `TimeExpires`：Token 过期时间（创建时按 7200 秒计算）。
  - `CorpId` / `CorpSecret`：企业 ID 与应用 Secret，用于过期后自动重新获取。
- `NewQywx(corpid, corpsecret)`：调用企业微信 `gettoken` 接口（`https://qyapi.weixin.qq.com/cgi-bin/gettoken`）获取 AccessToken 并创建实例；接口出错或 errcode 非 0 时返回错误。
- `GetAccessToken()`：返回当前 AccessToken；若已过期则自动用保存的 CorpId/CorpSecret 重新获取并更新缓存后返回。
- `MsgAccessToken`：gettoken 接口响应结构体（errcode、errmsg、access_token、expires_in）。
- `MsgMarkdown`：markdown 应用消息结构体，字段含：
  - `Touser` / `Toparty` / `Totag`：接收消息的用户/部门/标签（多个用 `|` 分隔）。
  - `Msgtype`：消息类型（markdown）。
  - `Agentid`：应用 ID。
  - `Markdown.Content`：markdown 内容。
  - `EnableDuplicateCheck`：是否开启重复消息检查。
- `MsgRestult`：消息发送结果结构体（errcode、errmsg、invaliduser、invalidparty、invalidtag 无效的接收者列表）。

## 使用示例

```go
qy, err := qywx.NewQywx("corpid", "corpsecret")
if err != nil {
    return
}
token := qy.GetAccessToken() // 过期后会自动刷新
```

过期判断基于 `TimeExpires` 字段（创建时间 + 7200 秒），超时后下次调用触发刷新。

## 备注

- 内部依赖本仓库 `httpx`（HTTP 客户端，Debug 模式开启）。
- 需要企业微信企业的 CorpId 与应用的 CorpSecret。
- Token 有效期按 7200 秒处理；`GetAccessToken` 在过期后惰性刷新，刷新失败时返回旧 Token。
- 本包当前只封装了 Token 获取与自动刷新；`MsgMarkdown`、`MsgRestult` 仅为结构体定义，发送消息需自行调用企业微信 `message/send` 接口。
