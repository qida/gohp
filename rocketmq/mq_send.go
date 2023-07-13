package rocketmq

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type SendMsg struct {
	IP       string
	Port     string
	Topic    string
	Group    string
	Producer rocketmq.Producer
}

//var Topic = "broker-a"
//var Group = "my_service"

func ConnRockerMQ(ip, port, topic, group string) *SendMsg {

	send := new(SendMsg)
	send.IP = ip
	send.Port = port
	send.Topic = topic
	send.Group = group

	addr, err := primitive.NewNamesrvAddr(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		panic(err)
	}

	send.Producer, err = rocketmq.NewProducer(
		producer.WithGroupName(group),
		producer.WithNameServer(addr),
		producer.WithCreateTopicKey(topic),
		producer.WithRetry(5))
	if err != nil {
		panic(err)
	}

	err = send.Producer.Start()
	if err != nil {
		panic(err)
	}

	return send
}

// SendSync 发送异步消息
func (p *SendMsg) SendSync(dataMsg []byte) {
	// 发送异步消息
	res, err := p.Producer.SendSync(context.Background(), primitive.NewMessage(p.Topic, dataMsg))
	if err != nil {
		fmt.Printf("send sync message error:%s\n", err)
	} else {
		fmt.Printf("send sync message success. result=%s\n", res.String())
	}

}

// SendAsync 发送消息后回调
// TODO 有高并发bug
func (p *SendMsg) SendAsync(dataMsg []byte) {
	var err error
	err = p.Producer.SendAsync(
		context.Background(),
		func(ctx context.Context, result *primitive.SendResult, err error) {
			if err != nil {
				fmt.Printf("receive message error:%v\n", err)
			} else {
				fmt.Printf("send message success. result=%s\n", result.String())
			}
		},
		primitive.NewMessage(p.Topic, dataMsg))
	if err != nil {
		fmt.Printf("send async message error:%s\n", err)
	}
}

// SendAsyncBatch // 批量发送消息
func (p *SendMsg) SendAsyncBatch(msgs []*primitive.Message) {

	for i := 0; i < len(msgs); i++ {
		msgs = append(msgs, primitive.NewMessage(p.Topic, []byte("batch send message. num:"+strconv.Itoa(i))))
	}
	res, err := p.Producer.SendSync(context.Background(), msgs...)
	if err != nil {
		fmt.Printf("batch send sync message error:%s\n", err)
	} else {
		fmt.Printf("batch send sync message success. result=%s\n", res.String())
	}

}

// DelaySendSync // 延迟发送消息
// delay send message
func (p *SendMsg) DelaySendSync(dataMsg []byte) {
	// 发送延迟消息
	msg := primitive.NewMessage(p.Topic, dataMsg)
	msg.WithDelayTimeLevel(3)
	res, err := p.Producer.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("delay send sync message error:%s\n", err)
	} else {
		fmt.Printf("delay send sync message success. result=%s\n", res.String())
	}
}

// TagSendSync 发送带有tag的消息
func (p *SendMsg) TagSendSync(tag string, dataMsg []byte) {
	msg1 := primitive.NewMessage(p.Topic, dataMsg)
	msg1.WithTag(tag)
	res, err := p.Producer.SendSync(context.Background(), msg1)
	if err != nil {
		fmt.Printf("send tag sync message error:%s\n", err)
	} else {
		fmt.Printf("send tag sync message success. result=%s\n", res.String())
	}
}

// Close 关闭
func (p *SendMsg) Close() {
	var err error
	err = p.Producer.Shutdown()
	if err != nil {
		panic(err)
	}
}
