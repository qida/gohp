# sms

短信验证码工具包，聚合阿里云、七牛、运满友科技三个短信渠道，并提供手机号格式校验与发送频率限制。

## 功能列表

### common.go（通用）

- `CheckRegexMobile(mobile)`：校验手机号格式（正则 `^1\d{10}$`，去除首尾空格），不合法时返回错误。
- `RequestRegLimit` / `LimitMax`：包级发送频次计数器（默认上限 2 次）。
- `TimerAddPoolLimit()`：每秒自动将计数器恢复 1（由包 `init()` 后台协程启动）。

### sms_ali.go（阿里云短信）

- `SMSAli`：阿里云短信客户端，封装 `dysmsapi.Client`。
- `NewSMSAli(accessKeyId, accessKeySecret)`：初始化客户端（RegionId 固定 cn-hangzhou），失败返回 nil。
- `Send(mobile)`：校验手机号后生成 4 位随机验证码并调用阿里云接口发送；内置签名 `UPM`、模板 `SMS_175475271`，返回验证码与错误。

### sms_manyou.go（运满友科技短信）

- `ManYou`：短信平台账号结构体（Account、Password）。
- `VerificationParam`：平台响应结构（codetype、sendmsgid）。
- `SendManYouSMS(man_you, mobile)`：校验手机号、检查频率限制后，通过 HTTP 接口（sdk.shmyt.cn）发送 6 位验证码，返回验证码与平台响应。

### sms_qiniu.go（七牛短信）

- `QiniuSMS`：短信参数结构体（SignatureID、TemplateID）。
- `InitQiniuSMS(accessKey, secretKey)`：初始化七牛短信 Manager（包级变量）。
- `SendQiniuSMS(qiniuSMS, mobile, code)`：校验手机号、检查频率限制后发送模板短信；`code` 为空时自动生成 4 位验证码。

## 使用示例

```go
// 阿里云
ali := sms.NewSMSAli("accessKeyId", "accessKeySecret")
code, err := ali.Send("13800138000")

// 七牛（需先初始化）
sms.InitQiniuSMS("ak", "sk")
err = sms.SendQiniuSMS(sms.QiniuSMS{SignatureID: "签名ID", TemplateID: "模板ID"}, "13800138000", "")
```

## 备注

- 阿里云依赖：`github.com/alibabacloud-go/dysmsapi-20170525/v4`、`darabonba-openapi/v2`、`tea`。
- 七牛依赖：`github.com/qiniu/api.v7/v7`（auth、sms）。
- 运满友渠道无第三方 SDK 依赖，直接 HTTP 调用，注意账号密码出现在 URL 中。
- 频率限制为进程内计数（每秒恢复 1 次、上限 2 次），多实例部署时不共享。
