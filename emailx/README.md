# emailx

SMTP 邮件发送工具包，基于 `jordan-wright/email` 封装纯文本与 HTML 邮件的发送。

## 功能列表

### email.go

- `Email`：SMTP 客户端结构体：
  - `Host`：SMTP 服务器地址，如 `smtp.gmail.com`。
  - `Port`：端口，如 `587`。
  - `Auth`：`smtp.PlainAuth` 认证对象。
- `Mail`：邮件内容结构体：
  - `From`：发件人地址。
  - `To`：收件人地址列表。
  - `Subject`：邮件主题。
  - `Text`：纯文本正文（[]byte，可选）。
  - `HTML`：超文本正文（[]byte，可选）。
- `NewEmail(host, port, username, password)`：创建 SMTP 客户端，使用 PlainAuth 认证。
- `SendEmail(mail)`：组装并发送邮件，Text 与 HTML 均为可选（非空才设置）；发送失败时返回包装后的错误（`发送邮件失败: ...`），不再直接 `log.Fatal` 终止程序。

## 使用示例

```go
e := emailx.NewEmail("smtp.gmail.com", "587", "xxx@gmail.com", "password")
err := e.SendEmail(emailx.Mail{
    From:    "xxx@gmail.com",
    To:      []string{"to@example.com"},
    Subject: "测试邮件",
    Text:    []byte("纯文本内容"),
    HTML:    []byte("<h1>Hello</h1>"),
})
if err != nil {
    fmt.Println("发送失败:", err)
}
```

## 备注

- 外部依赖：`github.com/jordan-wright/email`，标准库 `net`、`net/smtp`。
- 需要可用的 SMTP 服务器、账号与密码（或授权码）。
- 错误处理约定：`SendEmail` 失败返回 error（`%w` 包装原始错误），由调用方决定重试或记录日志。
