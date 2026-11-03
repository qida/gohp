# jsonx

JSON 工具集，提供不转义 HTML 的 JSON 编码、Unicode 转义还原，以及 nginx 配置文件读取。

## 功能列表

### json.go
- `JSONMarshal(t interface{}) ([]byte, error)`：序列化为 JSON 并关闭 HTML 转义（`&`、`<`、`>` 保持原样输出）
- `TransHtmlJson(data []byte) []byte`：将 JSON 中的 `\u0026`、`\u003c`、`\u003e` 转义还原为 `&`、`<`、`>`

### nginx.go
- `Upstream`：nginx upstream 配置结构（`Server` 字段）
- `ServConf`：nginx server 配置结构（`ServerName` 字段）
- `ReadNginxConf(file_path string) error`：读取 nginx 配置文件并解析为 `ServConf`，结果通过日志输出（基于 `github.com/recoye/config`）

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/jsonx"
)

func main() {
	data := map[string]string{"url": "https://example.com?a=1&b=2"}

	b, err := jsonx.JSONMarshal(data)
	if err != nil {
		return
	}
	fmt.Println(string(b)) // {"url":"https://example.com?a=1&b=2"}（& 不被转义）
}
```

## 备注

- `JSONMarshal` 使用 `Encoder.SetEscapeHTML(false)`，适合直接输出给前端/浏览器的场景。
- `nginx.go` 依赖第三方库 `github.com/recoye/config`。
