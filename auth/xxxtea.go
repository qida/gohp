// 用于对称加密功能
package auth

import (
	"encoding/base64"
	"strings"

	"github.com/xxtea/xxtea-go/xxtea"
)

type XxTea struct {
	Key string
}

func NewXxTea(key string) *XxTea {
	return &XxTea{Key: key}
}

// 加密
func (x *XxTea) EncryptStr(src []byte) (dst string) {
	dst = base64.URLEncoding.EncodeToString(xxtea.Encrypt(src, []byte(x.Key)))
	return
}

// 解密
func (x *XxTea) DecryptStr(src string) (dst string) {
	src = strings.TrimSpace(src)
	d, _ := base64.URLEncoding.DecodeString(src)
	dst = string(xxtea.Decrypt(d, []byte(x.Key)))
	return
}

// 加密
func (x *XxTea) Encrypt(debug bool, src []byte) (dst []byte) {
	if debug {
		dst = src
	} else {
		dst = xxtea.Encrypt(src, []byte(x.Key))
	}
	return
}

// 解密
func (x *XxTea) Decrypt(debug bool, src []byte) (dst []byte) {
	if debug {
		dst = src
	} else {
		dst = xxtea.Decrypt(src, []byte(x.Key))
	}
	return
}
