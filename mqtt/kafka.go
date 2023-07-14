package mqtt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goph/logx"

	kafka "github.com/segmentio/kafka-go"
)

type ClientKafka struct {
	Writer *kafka.Writer

	Server *server
}

type server struct {
	config kafka.ReaderConfig

	onNewClientCallback      func(c *ClientKafka)
	onClientConnectionClosed func(c *ClientKafka, err error)
	onNewMessage             func(c *ClientKafka, message kafka.Message)
}

func (c *ClientKafka) listen() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	c.Server.onNewClientCallback(c)
	reader := kafka.NewReader(c.Server.config)
	reader.SetOffsetAt(context.Background(), time.Now().Add(-1*time.Hour))
	defer reader.Close()
	for {
		select {
		case <-signals:
			logx.Error("结束信号，退出")
			return
		default:
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				c.Server.onClientConnectionClosed(c, err)
				return
			}
			c.Server.onNewMessage(c, message)
		}
	}
}
func (s *server) Listen(time_reconn time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("panic: %v", r)
		}
	}()
	client := &ClientKafka{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(s.config.Brokers...),
			Topic:    s.config.Topic,
			Balancer: &kafka.LeastBytes{},
		},
		Server: s,
	}
	defer client.Writer.Close()
	go func() {
		for {
			client.listen()
			time.Sleep(time_reconn)
		}
	}()
	// time.Sleep(1 * time.Second)
	// client.Send("时间", time.Now().Format("2006-01-02 15:04:05"))
}

func (s *server) OnNewClient(callback func(c *ClientKafka)) {
	s.onNewClientCallback = callback
}

func (s *server) OnNewMessage(callback func(c *ClientKafka, message kafka.Message)) {
	s.onNewMessage = callback
}

func (s *server) OnClientConnectionClosed(callback func(c *ClientKafka, err error)) {
	s.onClientConnectionClosed = callback
}

func NewKAFKA(config kafka.ReaderConfig) *server {
	server := &server{
		config: config,
	}
	server.OnNewClient(func(c *ClientKafka) {})
	server.OnNewMessage(func(c *ClientKafka, message kafka.Message) {})
	server.OnClientConnectionClosed(func(c *ClientKafka, err error) {})

	return server
}

func (c *ClientKafka) Send(key string, msg string) (_err error) {
	_err = c.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(msg),
		},
	)
	if _err != nil {
		logx.Errorf("Kafka failed to write messages:%+v", _err)
	}
	return
}
