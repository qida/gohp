# cryptox

加解密工具包：AES、MD5、SHA1、HMAC-SHA1、Base64 与随机数等常用功能的统一实现（原 crypto 包已删除并与本包合并，本包是唯一的加解密实现）。

导入：`import "github.com/qida/gohp/cryptox"`

## 功能列表

### aes.go — AES 加解密
- `AESEncodeByte(data, key)` / `AESDecodeByte(data, key)`：AES CFB 模式加解密（字节）
- `AESEncode(data, key)` / `AESDecode(data, key)`：AES CFB 模式加解密（字符串，密文经 Base64 编码）
- `AesEncrypt(origData, key)` / `AesDecrypt(crypted, key)`：AES CBC 模式加解密（PKCS7 填充）

### md5.go — MD5
- `GetMD5Srting(data)`：计算字符串 MD5（十六进制小写）
- `GetMD5Srtings(data)`：将字符串数组 JSON 序列化后计算 MD5
- `GetMD5Byte(byts)`：计算字节切片 MD5
- `CheckMD5(data, sum)`：校验字符串 MD5 是否匹配

### sha1.go — SHA1
- `GetSha1(data)`：计算字符串 SHA1（十六进制小写）
- `GetHmacSha1(data, key)`：计算 HMAC-SHA1

### base64.go — Base64
- `Base64EncodeByte(data)` / `Base64DecodeByte(data)`：字节版标准 Base64 编解码
- `Base64Encode(data)` / `Base64Decode(data)`：字符串版标准 Base64 编解码

### random.go — 安全随机数
- `RandomDigits(length)`：生成指定长度的 0-9 随机数字（可用作验证码），内部基于 crypto/rand 并规避模偏差

### crypto.go — 随机 Key 与填充
- `GetRandomKey()`：基于 UUID v4 + SHA1 生成 40 位随机密钥
- `PKCS7Padding(ciphertext, blockSize)` / `PKCS7UnPadding(origData)`：PKCS7 填充 / 去填充

### crypto_test.go — 测试
- 覆盖 MD5、SHA1、HMAC-SHA1、Base64、AES 加解密及随机 Key 的正确性

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/cryptox"
)

func main() {
	key := "astaxie12798akljzmknm.ahkjkljl;k" // 32 字节密钥
	enc, _ := cryptox.AESEncode("test skcloud crypto", key)
	fmt.Println(enc) // XGFffDlOwueH4DCVD+5Um/S/OA==

	dec, _ := cryptox.AESDecode(enc, key)
	fmt.Println(dec) // test skcloud crypto
}
```

## 备注
- `GetMD5Srting` / `GetMD5Srtings` 为历史拼写（Srting），保留以兼容旧代码。
- AES 两种模式均使用 key 的前 16 字节作为 IV；密钥长度须为 16/24/32 字节。
- `random.go` 源自 Dmitry Chestnykh 的 captcha 项目（MIT 协议）。
