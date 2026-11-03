# plugin

第三方平台插件工具包，当前包含阿里云直播推流/播放地址的鉴权生成。

## 功能列表

### live_ali.go（阿里云直播）

- `LiveUrl`：直播地址结构体，字段包括：
  - `AppName`：直播应用名（作为流路径的一部分）。
  - `PullUrl` / `PushUrl`：拉流域名 / 推流域名。
  - `Key`：URL 鉴权密钥。
  - `PushRtmp`：生成的 RTMP 推流地址。
  - `PlayRtmp` / `PlayM3u8` / `PlayFlv`：生成的 RTMP / HLS(m3u8) / FLV 播放地址。
  - `TimeUpdate`：地址刷新时间。
- `NewLiveUrl(app_name, pull_url, key)`：初始化直播配置（应用名、拉流域名、鉴权 Key），不涉及网络请求。
- `RefreshUrl(id_live)`：按直播 ID 重新生成全部地址，规则：
  - 有效期固定为当前时间起 24 小时。
  - 采用阿里云 URL 鉴权格式 `auth_key=timestamp-rand-uid-md5hash`（rand、uid 固定为 0）。
  - 对推流（rtmp）、播放（rtmp / m3u8 / flv）分别按 URI 计算 MD5 签名后拼接到查询参数。

## 使用示例

```go
var live plugin.LiveUrl
live.NewLiveUrl("liveapp", "pull.example.com", "鉴权Key")
live.PushUrl = "push.example.com" // NewLiveUrl 不设置推流域名，需自行赋值
if err := live.RefreshUrl(1001); err == nil {
    fmt.Println("推流地址:", live.PushRtmp)
    fmt.Println("播放地址:", live.PlayFlv, live.PlayM3u8)
}
```

## 备注

- 无第三方依赖，仅使用标准库（crypto/md5、encoding/hex、fmt、time）。
- 需要已开通阿里云直播服务，并配置好推流/拉流域名与 URL 鉴权 Key。
- 地址有效期 24 小时，过期后需再次调用 `RefreshUrl` 重新生成。
- `NewLiveUrl` 与 `RefreshUrl` 均只做本地计算，不发起网络请求。
- m3u8 / flv 播放地址的签名 URI 含对应扩展名，与 rtmp 的签名不同。
