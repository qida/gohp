package httpx

import (
	"context"
	"fmt"

	"goph/logx"

	"github.com/go-resty/resty/v2"
)

type ClientHttp struct {
	client *resty.Client
}

func NewClientHttp() *ClientHttp {
	client := resty.New().SetContentLength(true)
	client.Header.Set("User-Agent", "LabelCenter HTTP CLIENT")
	client.Header.Set("Content-Type", "application/json;charset=utf-8")
	return &ClientHttp{
		client: client,
	}
}
func (t *ClientHttp) Debug() *ClientHttp {
	t.client.SetDebug(true)
	return t
}

func (t *ClientHttp) PostBody(url string, req, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	r, _err := t.client.R().
		EnableTrace().
		SetBody(req).
		SetResult(resp).
		Post(url)
	if _err != nil {
		logx.Errorf(" 错误:%+v", _err)
		return
	}
	fmt.Println("Trace Info:", r.Request.TraceInfo())
	return
}
func (t *ClientHttp) Post(ctx context.Context, url string, req map[string]string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	_, _err = t.client.R().
		SetQueryParams(req).
		SetResult(resp).
		Post(url)
	return
}
func (t *ClientHttp) Get(ctx context.Context, url string, resp interface{}, header map[string]string) (_err error) {
	defer func() {
		t.client.SetCloseConnection(true)
	}()
	for k, v := range header {
		t.client.Header.Set(k, v)
	}
	_, _err = t.client.R().
		SetResult(resp).
		Get(url)
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
	_, _err = t.client.R().
		SetResult(resp).
		Get(url)
	return
}