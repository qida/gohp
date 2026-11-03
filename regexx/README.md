# regexx

常用正则校验工具，在包加载时预编译邮箱与 URL 正则，提供常见格式的字符串校验函数。

## 功能列表

### regex.go
- `IsEmail(email string) bool`：校验是否为邮箱地址
  - 宽松规则（`(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`）
  - 基本校验可覆盖约 99% 场景
- `IsEmailRFC(email string) bool`：按 RFC 2822 严格规则校验邮箱地址
  - 更严格地校验本地部分与域名字符集
- `IsUrl(url string) bool`：校验是否为 URL 链接
  - 支持 ftp/http/https 协议，可包含用户名、端口、路径等
  - 简单校验可覆盖约 99% 场景

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/regexx"
)

func main() {
	fmt.Println(regexx.IsEmail("user@example.com"))    // true
	fmt.Println(regexx.IsEmail("invalid-email"))       // false
	fmt.Println(regexx.IsEmailRFC("user@example.com")) // true
	fmt.Println(regexx.IsUrl("https://example.com/a?b=1")) // true
	fmt.Println(regexx.IsUrl("ftp://127.0.0.1:21"))    // true
}
```

## 备注

- 三个正则均在 `init()` 中通过 `regexp.MustCompile` 一次性编译，多次调用无需重复编译，适合高频校验场景。
- `IsEmail` 与 `IsEmailRFC` 的区别在于规则严格程度，按业务需要选择。
- 包内正则模式均为私有常量，不支持外部自定义。
