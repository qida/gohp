package cryptox

import (
	"bytes"

	uuid "github.com/nu7hatch/gouuid"
)

func GetRandomKey() string {
	t, _ := uuid.NewV4()
	r := GetSha1(t.String())
	return r
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
