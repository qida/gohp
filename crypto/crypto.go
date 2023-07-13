package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/nu7hatch/gouuid"
)

func GetRandomKey() string {
	t, _ := uuid.NewV4()
	r := GetSha1(t.String())
	return r
}

func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func GetHmacSha1(data string, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func GetMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func CheckMD5(data string, sum string) bool {
	out := GetMD5(data)
	return sum == out
}

func Base64EncodeByte(data []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func Base64DecodeByte(data []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(data))
}

func Base64Encode(data string) string {
	return string(Base64EncodeByte([]byte(data)))
}

func Base64Decode(data string) (string, error) {
	d, e := Base64DecodeByte([]byte(data))
	return string(d), e
}

func AESEncodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBEncrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESDecodeByte(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}
	var iv = key[:aes.BlockSize]
	blockMode := cipher.NewCFBDecrypter(block, iv)
	dest := make([]byte, len(string(data)))
	blockMode.XORKeyStream(dest, data)
	return dest, nil
}
func AESEncode(data string, key string) (string, error) {
	out, err := AESEncodeByte([]byte(data), []byte(key))
	return string(Base64EncodeByte(out)), err
}
func AESDecode(data string, key string) (string, error) {
	d, e := Base64DecodeByte([]byte(data))
	if e != nil {
		return data, e
	}
	out, err := AESDecodeByte(d, []byte(key))
	return string(out), err
}
