package logx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/qida/gohp/rocketmq"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type LogRocketMQ struct {
	producer *rocketmq.MqProducer
	Topic    string
	Group    string
	Address  []string
	Level    *Level
	fields   []zap.Field
}

func (l *LogRocketMQ) Write(p []byte) (int, error) {
	var err error

	//fmt.Println("Body1 ======>:", string(p))

	MqMsg := rocketmq.NewMQMessage(l.Topic, "logs", "logger", p, 0)
	_, err = l.producer.SendAsyncMessage(MqMsg, func(ctx context.Context, result *primitive.SendResult, e error) {
		if e != nil {
			fmt.Println("MQ 异步发送消息 err: ", e)
			err = e
			return
		}
	})

	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func getRocketMQWriter(topic string, group string, SrvAddr []string, fields ...zap.Field) (io.Writer, error) {
	lr := &LogRocketMQ{
		Topic:  topic,
		Group:  group,
		fields: fields,
	}
	lr.producer = rocketmq.NewMQProducer(topic, group, SrvAddr)
	return lr, nil
}

/*************************初始化MQ***************************/
/*************************Rocket_MQ*************************/

type LogRocket struct {
	Topic string
	Group string
	Tag   string
	Key   string
	*rocketmq.MqProducer
	LogBody
}

// NewProbeLog 初始化 数据探针
func NewProbeLog(m map[string]interface{}) (*LogRocket, error) {

	v, ok := m["rocketmq"]
	if !ok {
		return nil, errors.New("查询 rocketmq 失败")
	}

	lrs, err := toRocketMQ(v)
	if err != nil {
		return nil, err
	}

	if len(lrs) <= 0 {
		return nil, errors.New("mq 配置是空值")
	}

	mq := rocketmq.NewMQProducer(lrs[0].Topic, lrs[0].Group, lrs[0].Address)

	return &LogRocket{
		Topic:      lrs[0].Topic,
		Group:      lrs[0].Group,
		Key:        lrs[0].Group,
		Tag:        lrs[0].Group,
		MqProducer: mq,
		LogBody:    LogBody{},
	}, nil
}

func (l *LogRocket) MQDebug(title, trans, action string, fields ...LogField) {
	l.sendLogMQ("DEBUG", title, trans, action, "logs", fields...)
}

func (l *LogRocket) MQInfo(title, trans, action string, fields ...LogField) {
	l.sendLogMQ("INFO", title, trans, action, "logs", fields...)
}

func (l *LogRocket) MQWarn(title, trans, action string, fields ...LogField) {
	l.sendLogMQ("WARN", title, trans, action, "logs", fields...)
}

func (l *LogRocket) MQError(title, trans, action string, fields ...LogField) {
	l.sendLogMQ("ERROR", title, trans, action, "logs", fields...)
}

func (l *LogRocket) MQPanic(title, trans, action string, fields ...LogField) {
	l.sendLogMQ("Panic", title, trans, action, "logs", fields...)
}

func (l *LogRocket) SetTagKey(tag, key string) *LogRocket {
	l.Tag = tag
	l.Key = key
	return l
}

func (l *LogRocket) sendLogMQ(level, title, trans, action, measurement string, fields ...LogField) {

	if l == nil {
		fmt.Println("数据采集 没有被初始化")
		return
	}

	fileName, line, funcName := "???", 0, "???"

	pc, fileName, line, ok := runtime.Caller(2)
	if !ok {
		fileName, line, funcName = "???", 0, "???"
	} else {
		funcName = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcName = filepath.Ext(funcName)            // .foo
		funcName = strings.TrimPrefix(funcName, ".") // foo
	}

	l.LogBody = LogBody{
		LogHeader: LogHeader{
			TransId:        0,
			Level:          level,
			Time:           time.Now().Local().Format(time.RFC3339),
			File:           fmt.Sprintf("%s:%d", fileName, line),
			ServiceCluster: "",
			ServiceID:      "",
			ServiceName:    "",
			Title:          title,
			Trans:          trans,
			Action:         action,
			Method:         funcName,
			DataType:       "",
		},
		Measurement:  measurement,
		AppID:        "",
		RawData:      "",
		ElapsedTime:  0,
		Fields:       fields,
		ClientIP:     "",
		StartTime:    "",
		EndTime:      "",
		TimeUsed:     0,
		TimeUsedType: "",
	}

	b, err := json.Marshal(&l.LogBody)
	if err != nil {
		fmt.Println("json解析异常")
		return
	}

	MqMsg := rocketmq.NewMQMessage(l.Topic, l.Tag, l.Key, b, 0)
	_, err = l.MqProducer.SendAsyncMessage(MqMsg, func(ctx context.Context, result *primitive.SendResult, e error) {
		if e != nil {
			fmt.Println("MQ 异步发送消息 err: ", e)
			return
		}
	})

	if err != nil {
		fmt.Println("MQ 异步发送消息 err: ", err)
		return
	}
}
