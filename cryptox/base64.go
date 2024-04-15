package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"

	uuid "github.com/nu7hatch/gouuid"
)

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