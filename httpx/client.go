package httpx

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type FileLogger struct {
	file *os.File
}

func (l *FileLogger) Errorf(format string, v ...interface{}) {
	l.file.WriteString(fmt.Sprintf(format, v...))
}
func (l *FileLogger) Warnf(format string, v ...interface{}) {
	l.file.WriteString(fmt.Sprintf(format, v...))
}
func (l *FileLogger) Debugf(format string, v ...interface{}) {
	l.file.WriteString(fmt.Sprintf(format, v...))
}

type ClientHttp struct {
	client *resty.Client
}

var logger *FileLogger
var once sync.Once

func NewClientHttp() *ClientHttp {
	client := resty.New().SetContentLength(true).
		SetHeader("User-Agent", "HTTP Client").
		SetHeader("Content-Type", "application/json;charset=utf-8").
		SetTimeout(time.Second * 10).
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(20 * time.Second)
	return &ClientHttp{
		client: client,
	}
}
func (t *ClientHttp) Debug(debug bool) *ClientHttp {
	t.client.SetDebug(debug)
	return t
}

func (t *ClientHttp) LogFile() *ClientHttp {
	once.Do(func() {
		if _, err := os.Stat("./log"); os.IsNotExist(err) {
			os.Mkdir("./log", os.ModePerm)
		}
		// 创建日志文件
		logFile, err := os.OpenFile("./log/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		// defer logFile.Close()
		logger = &FileLogger{file: logFile}
	})
	t.client.SetLogger(logger)
	//t.client.SetLogger(log.New(logFile, "", log.LstdFlags))
	return t
}

func (t *ClientHttp) SetTimeout(time_out_second int) *ClientHttp {
	t.client.SetTimeout(time.Duration(time_out_second * int(time.Second)))
	return t
}

func (t *ClientHttp) PostBody(url string, req, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	var r *resty.Response
	r, _err = t.client.R().
		SetBody(req).
		SetResult(resp).
		// SetError(_err).
		Post(url)
	if _err != nil {
		return
	}
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) Post(url string, req map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	var r *resty.Response
	r, _err = t.client.R().
		SetFormData(req).
		SetResult(resp).
		Post(url)
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) UploadFile(url string, name_param string, file_path string, params map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	// 设置请求头
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	// 执行文件上传
	r, _err := t.client.R().
		SetFile(name_param, file_path).
		SetFormData(params).
		SetResult(resp).
		Post(url)
	if _err != nil {
		return
	}
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) UploadFileFromBytes(url string, name_param string, file_name string, file_data []byte, params map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	// 设置请求头
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	// 创建一个bytes.Reader，用于读取fileData
	reader := bytes.NewReader(file_data)
	// 执行文件上传
	r, _err := t.client.R().
		SetFileReader(name_param, file_name, reader).
		SetFormData(params).
		SetResult(resp).
		Post(url)
	if _err != nil {
		_err = fmt.Errorf("服务器解析错误 [%s] Result:[%+v] Error:%+v", url, r.Result(), _err)
		return
	}
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) Get(url string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	var r *resty.Response
	r, _err = t.client.R().
		SetResult(resp).
		Get(url)
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) GetParams(url string, req map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	var r *resty.Response
	r, _err = t.client.R().
		SetQueryParams(req).
		SetResult(resp).
		Get(url)
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 [%s] %d %s", url, r.StatusCode(), r.Status())
		return
	}
	return
}



