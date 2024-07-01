package httpx

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// 读文件内容从multipart
func ReadFileDataFromMultipart(f multipart.File) (data *[]byte, err error) {
	d := make([]byte, 0)
	buf := make([]byte, 512)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		d = append(d, buf[:n]...)
	}
	data = &d
	return
}

// http发送二进制数据
func NewFileDataUploadRequest(url, path string, data *[]byte, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	file := bytes.NewReader(*data)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}
