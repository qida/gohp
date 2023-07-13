package rocketmq

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/panjf2000/ants/v2"
)

type MqProducer struct {
	pool     *ants.Pool
	producer rocketmq.Producer
}

// NewMQProducer 生产者
func NewMQProducer(topicName, groupName string, NameSrvAddr []string) *MqProducer {
	pl, _ := ants.NewPool(10000)

	p, err := rocketmq.NewProducer(
		producer.WithNameServer(NameSrvAddr),
		// producer.WithNsResolver(primitive.NewPassthroughResolver(NamesrvAddr)),
		producer.WithCreateTopicKey(topicName), //创建主题key
		producer.WithGroupName(groupName),      //组名
		//producer.WithRetry(defaultRetry),
		producer.WithRetry(5), // 重试次数
	)

	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}

	return &MqProducer{pool: pl, producer: p}
}

// Shutdown 关闭
func (mq *MqProducer) Shutdown() error {
	err := mq.producer.Shutdown()
	if err != nil {
		fmt.Println("CloseMQProducer error", err)
	}
	mq.pool.Release()
	return err
}

// SendSyncMessage 发送同步消息
func (mq *MqProducer) SendSyncMessage(msg *MQMessage) (bool, error) {
	res, err := mq.producer.SendSync(context.Background(), msg.toMessage())
	if err != nil {
		fmt.Printf("send sync message error:%s\n", err)
		return false, err
	} else {
		fmt.Printf("send sync message success. result=%s\n", res.String())
	}
	return true, nil
}

// SendSyncbatchMessage SendSyncMessageList 发送批量消息
func (mq *MqProducer) SendSyncbatchMessage(msgs []*MQMessage) (bool, error) {
	mqmgs := make([]*primitive.Message, 0)
	for _, msg := range msgs {
		mqmgs = append(mqmgs, msg.toMessage())
	}

	res, err := mq.producer.SendSync(context.Background(), mqmgs...)
	if err != nil {
		fmt.Printf("send sync message error:%s\n", err)
		return false, err
	} else {
		fmt.Printf("send sync message success. result=%s\n", res.String())
	}
	return true, nil
}

// SendAsyncMessageBody 发送异步消息后回调
func (mq *MqProducer) SendAsyncMessageBody(msg *MQMessage) (bool, error) {

	fmt.Println("body 1 发送异步消息后回调 ========>:", msg.Key, msg.TagName)

	fmt.Println("body 2 发送异步消息后回调 ========>:", string(msg.Body))
	err := mq.producer.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive message send error: %s\n", e)
				return
			}

		}, primitive.NewMessage(msg.TopicName, msg.Body))

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
		return false, err
	}

	return true, nil
}

// SendAsyncMessage 发送异步消息后回调
func (mq *MqProducer) SendAsyncMessage(msg *MQMessage, callback func(ctx context.Context, result *primitive.SendResult, err error)) (bool, error) {

	//err := mq.pool.Submit(func() {
	//fmt.Println("body  发送异步消息后回调 ========>:", string(msg.Body))
	e := mq.producer.SendAsync(context.Background(), callback, msg.toMessage())
	if e != nil {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		callback(ctx, nil, fmt.Errorf("mqproducer send async message error: %w", e))
		fmt.Printf("mqproducer send async message error:%s\n", e)
		return false, e
	}
	//})
	//
	//if err != nil {
	//	fmt.Printf("pool Submit func error:%s\n", err)
	//	return false, err
	//}
	return true, nil
}
