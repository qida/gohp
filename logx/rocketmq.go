package logx

import (
	"context"
	"io"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

// LogRocketmq 写入RocketMQ
type LogRocketmq struct {
	producer rocketmq.Producer
	Topic    string
	fields   []zap.Field
}

// Write 实现io.Writer接口
func (lr *LogRocketmq) Write(p []byte) (int, error) {
	msg := &primitive.Message{
		Topic: lr.Topic,
		Body:  p,
	}
	_, err := lr.producer.SendSync(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// getRocketmqWriter .
func getRocketmqWriter(topic string, nameServer []string, group string, fields ...zap.Field) (io.Writer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(primitive.NamesrvAddr(nameServer)),
		producer.WithGroupName(group),
		producer.WithRetry(2),
	)
	if err != nil {
		return nil, err
	}
	if err := p.Start(); err != nil {
		return nil, err
	}
	return &LogRocketmq{producer: p, Topic: topic, fields: fields}, nil
}
