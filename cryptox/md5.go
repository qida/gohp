package cryptox

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func GetMD5Srting(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMD5Srtings(data []string) string {
	if len(data) == 0 {
		return ""
	}
	h := md5.New()
	bytes, _ := json.Marshal(data)
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

func GetMD5Byte(byts []byte) string {
	h := md5.New()
	h.Write(byts)
	return hex.EncodeToString(h.Sum(nil))
}

func CheckMD5(data string, sum string) bool {
	out := GetMD5Srting(data)
	return sum == out
}
