package zh

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Encode(src string) (dst string) {
	data, err := io.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GB18030.NewEncoder()))
	if err == nil {
		dst = string(data)
	}
	return
}
func Decode(src string) (dst string) {
	data, err := io.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GB18030.NewDecoder()))
	if err == nil {
		dst = string(data)
	}
	return
}
