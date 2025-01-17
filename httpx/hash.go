package httpx

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func HashHttp(file *multipart.FileHeader) (_hash string, _err error) {
	src, _err := file.Open()
	if _err != nil {
		return
	}
	defer src.Close()
	h := sha256.New()
	if _, _err = io.Copy(h, src); _err != nil {
		return
	}
	_hash = hex.EncodeToString(h.Sum(nil))
	return
}

func HashFile(file *multipart.File) (_hash string, _err error) {
	h := sha256.New()
	if _, _err = io.Copy(h, *file); _err != nil {
		return
	}
	_hash = hex.EncodeToString(h.Sum(nil))
	return
}
