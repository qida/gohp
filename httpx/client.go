package httpx

import (
	"context"
	"fmt"
	"log"
	"os"
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

func init() {
	// 创建日志文件
	// logFile, err := os.OpenFile("request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logFile, err := os.OpenFile("./log/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// defer logFile.Close()
	logger = &FileLogger{file: logFile}
}

func NewClientHttp() *ClientHttp {
	client := resty.New().SetContentLength(true).
		SetHeader("User-Agent", "HTTP CLIENT").
		SetHeader("Content-Type", "application/json;charset=utf-8").
		SetTimeout(time.Second * 10).
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second)
	return &ClientHttp{
		client: client,
	}
}
func (t *ClientHttp) Debug(debug bool) *ClientHttp {
	t.client.SetDebug(debug)
	return t
}
func (t *ClientHttp) LogFile() *ClientHttp {
	t.client.SetLogger(logger)
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
		_err = fmt.Errorf("服务器异常 %d %s", r.StatusCode(), r.Status())
		return
	}
	return
}
func (t *ClientHttp) Post(ctx context.Context, url string, req map[string]string, resp interface{}, header map[string]string) (_err error) {
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
		Post(url)
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 %d %s", r.StatusCode(), r.Status())
		return
	}
	return
}
func (t *ClientHttp) Get(ctx context.Context, url string, resp interface{}, header map[string]string) (_err error) {
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
		_err = fmt.Errorf("服务器异常 %d %s", r.StatusCode(), r.Status())
		return
	}
	return
}

func (t *ClientHttp) GetParams(ctx context.Context, url string, req map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	t.client.SetQueryParams(req)
	var r *resty.Response
	r, _err = t.client.R().
		SetResult(resp).
		Get(url)
	if r.StatusCode() != 200 {
		_err = fmt.Errorf("服务器异常 %d %s", r.StatusCode(), r.Status())
		return
	}
	return
}
