package httpx

import (
	"context"
	"fmt"
	"time"

	"github.com/qida/gohp/logx"

	"github.com/go-resty/resty/v2"
)

type ClientHttp struct {
	client *resty.Client
}

func NewClientHttp() *ClientHttp {
	client := resty.New().SetContentLength(true)
	client.Header.Set("User-Agent", "LabelCenter HTTP CLIENT")
	client.Header.Set("Content-Type", "application/json;charset=utf-8")
	client.SetTimeout(time.Second * 10)

	return &ClientHttp{
		client: client,
	}
}
func (t *ClientHttp) Debug(debug bool) *ClientHttp {
	t.client.SetDebug(debug)
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
		logx.Errorf(" 错误:%+v", _err)
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
