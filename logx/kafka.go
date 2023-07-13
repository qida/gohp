package logx

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"io"
)

// LogKafka 写入kafka
type LogKafka struct {
	producer sarama.SyncProducer
	Topic    string
	Address  []string
	Level    *Level
	fields   []zap.Field
}

// Write 实现io.Writer接口
func (lk *LogKafka) Write(p []byte) (int, error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err := lk.producer.SendMessage(msg)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// getKafkaWriter .
func getKafkaWriter(topic string, address []string, fields ...zap.Field) (io.Writer, error) {
	kl := &LogKafka{
		Topic:  topic,
		fields: fields,
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	var err error
	kl.producer, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		return nil, err
	}
	return kl, nil
}
