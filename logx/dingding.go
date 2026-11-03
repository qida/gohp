package logx

import (
	"encoding/json"
	"fmt"
	"io"

	"go.uber.org/zap"
)

// LogDingding .
type LogDingding struct {
	robot  *Robot
	level  *Level
	fields []zap.Field
}

// Write 实现io.Writer
// 将 zap 日志内容格式化为钉钉消息并发送到钉钉群（复用 dingtalk.go 中的 Robot）
func (ld *LogDingding) Write(p []byte) (int, error) {
	var msg map[string]interface{}
	if err := json.Unmarshal(p, &msg); err == nil {
		var fields string
		for _, f := range ld.fields {
			fields += fmt.Sprintf("\r\n###### %s: %s\r\n", f.Key, f.String)
		}
		msg = map[string]interface{}{
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
		msg = map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": string(p),
			},
		}
	}

	if err := ld.robot.SendMessage(msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

func getDingdingWriter(secret, access_token string, fields ...zap.Field) (io.Writer, error) {
	return &LogDingding{
		robot:  NewRobot(access_token, secret),
		fields: fields,
	}, nil
}
