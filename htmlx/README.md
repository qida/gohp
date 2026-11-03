# htmlx

HTML 字符串处理工具包，提供 HTML 与 JS 格式互转、编解码、标签清理等纯函数。

## 功能列表

### html.go

- `Html2JS(data []byte) []byte`：将 HTML 内容转换为 JS 字符串格式，处理规则：
  - `\` 转义为 `\\`；换行 `\n` 转义为 `\n` 字面量；删除 `\r`；`"` 转义为 `\"`。
  - `<table>` 转义为 `&lt;table>`（避免与 JS 脚本标签冲突）。
- `HtmlEncode(str string) string`：HTML 转义，等价于 `html.EscapeString`（如 `<` → `&lt;`）。
- `HtmlDecode(str string) string`：HTML 反转义，等价于 `html.UnescapeString`（如 `&amp;` → `&`）。
- `StripTags(src string) string`：去除 HTML 中的 style/script 块、HTML 注释与各类标签，并将 2 个以上连续空白压缩为单个换行，最终去除首尾空白，返回纯文本。
- `Nl2br(str string) string`：将 `\n` 换行符替换为 `<br/>`，用于纯文本在网页中展示。

## 使用示例

```go
// 来自 html_test.go 的用例
js := string(htmlx.Html2JS([]byte(`<div id="button" class="btn">Click me</div>` + "\n\r")))
// 结果: `<div id=\"button\" class=\"btn\">Click me</div>\n`

// 标签清理
text := htmlx.StripTags("<style>a{}</style><!--c--><p>hello</p>")
// 结果: "hello"

// 编解码与换行转换
s := htmlx.HtmlEncode("<a>")  // &lt;a&gt;
htmlx.HtmlDecode(s)            // <a>
htmlx.Nl2br("a\nb")            // a<br/>b
```

## 备注

- 无外部依赖，仅使用标准库 `html`、`regexp`、`strings`。
- 源码源自 2013 年 go-committers 的代码，遵循 Apache-2.0 许可证。
- 所有函数均为无副作用纯函数，可并发安全调用。
- `StripTags` 通过正则一次性处理 style/script 块、注释与标签三类内容。
- 附带 `html_test.go` 覆盖 `Html2JS` 的转换与基准测试。
