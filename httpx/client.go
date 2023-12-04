package httpx

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type ClientHttp struct {
	client *resty.Client
}

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	// 设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// logger.SetOutput(os.Stdout)
	logFile, err := os.OpenFile("./log/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("failed to log to file.")
		return
	}
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	writers := []io.Writer{
		logFile,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger.SetOutput(fileAndStdoutWriter)
	//设置最低loglevel
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("init logger success")
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
