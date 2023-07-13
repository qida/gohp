package logx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// LogDingding .
type LogDingding struct {
	secret       string
	access_token string
	level        *Level
	fields       []zap.Field
}

var dingurl = "https://oapi.dingtalk.com/robot/send"

// Write 实现io.Writer
// 发送给钉钉群
func (ld *LogDingding) Write(p []byte) (int, error) {
	ts := time.Now().UnixNano() / 1e6
	string_to_sign := fmt.Sprintf("%d\n%s", ts, ld.secret)

	h := hmac.New(sha256.New, []byte(ld.secret))

	_, _ = h.Write([]byte(string_to_sign))

	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	var m map[string]interface{}
	var msg map[string]interface{}
	err := json.Unmarshal(p, &msg)

	if err == nil {
		var fields string
		for _, f := range ld.fields {
			fields += fmt.Sprintf("\r\n###### %s: %s\r\n", f.Key, f.String)
		}
		m = map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": fmt.Sprintf("%s: %s", msg["level"], msg["file"]),
				"text": fmt.Sprintf(
					`##### **%v:%v%v**

---

###### 级别: %s

%s

###### 文件: %s

###### 时间: %s
	`, msg["func"], msg["error"], msg["msg"], msg["level"], fields, msg["file"], msg["time"]),
			},
		}
	} else {
		m = map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": string(p),
			},
		}
	}
	/* 效果
	Thrift服务端口: 22725
	级别: INFO
	位置: 标签中心 市侧
	文件: initialize/thrift.go:40
	时间: 2022-09-21 17:27:58.857
	*/
	byt, _ := json.Marshal(m)
	params := url.Values{}
	params.Add("access_token", ld.access_token)
	params.Add("timestamp", strconv.FormatInt(ts, 10))
	params.Add("sign", sign)
	resp, err := http.Post(dingurl+"?"+params.Encode(), "application/json", bytes.NewReader(byt))
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	// _, _ = io.ReadAll(resp.Body)
	return 0, nil
}

func getDingdingWriter(secret, access_token string, fields ...zap.Field) (io.Writer, error) {
	return &LogDingding{
		secret:       secret,
		access_token: access_token,
		fields:       fields,
	}, nil
}
