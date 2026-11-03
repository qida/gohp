# auth

令牌（token）工具包：基于 XXTEA 对称加密生成与解析携带用户信息的 token，主要用于 token 类功能。

导入：`import "github.com/qida/gohp/auth"`

## 功能列表

### auth.go — Token 生成与解析
- `Auth`：认证结构（`User` 载荷、`XxTea` 加密器）
- `NewAuth(key)`：使用密钥创建 Auth
- `Encrypt(uid, s)`：将用户对象加密为 token，格式为 "uid:Base64URL密文"（附带随机 salt 与时间戳）
- `Decrypt(uid_auth)`：解析 token，还原出 `Auth`（含 User 载荷）

### xxxtea.go — XXTEA 对称加密
- `XxTea`：XXTEA 加密器（底层为 `github.com/xxtea/xxtea-go`）
- `NewXxTea(key)`：创建 XxTea
- `EncryptStr(src)` / `DecryptStr(src)`：字符串加密 / 解密（密文经 Base64 URL 编码）
- `Encrypt(debug, src)` / `Decrypt(debug, src)`：字节版加密 / 解密，`debug=true` 时直接透传不加密

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/auth"
)

func main() {
	a := auth.NewAuth("your-secret-key")

	type User struct {
		Name string
	}
	token := a.Encrypt(1001, User{Name: "张三"})
	fmt.Println(token) // 1001:xxxx

	res, err := a.Decrypt(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.User) // map[Name:张三]
}
```

## 备注
- token 结构为 `uid:Base64URL(XXTEA(JSON(Auth)))`；`Decrypt` 也兼容只传密文部分（不含 "uid:" 前缀）。
- `Encrypt` 每次生成都会附带 0-999 的随机 salt 与当前时间，同一用户每次产出的 token 不同。
