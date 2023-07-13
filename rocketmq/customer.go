package rocketmq

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type MqCustomer struct {
	pushConsumer   rocketmq.PushConsumer //推流模式
	pullConsumer   rocketmq.PullConsumer //拉流模式
	subTopic       []string
	groupName      string
	isBroadCasting bool
	QueuePull      primitive.MessageQueue
	MessageList    chan *primitive.MessageExt
}

// NewPushCustomer 创建一个MQ消费者
func NewPushCustomer(groupName string, NameSrvAddr []string, isBroadCasting bool) *MqCustomer {
	consumerModel := consumer.Clustering //集群消费模式
	if isBroadCasting {
		consumerModel = consumer.BroadCasting //广播模式
	}
	newPushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(NameSrvAddr), // consumer.WithNsResolver(primitive.NewPassthroughResolver(NamesrvAddr)),
		consumer.WithGroupName(groupName),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), // 选择消费时间(首次/当前/根据时间)
		consumer.WithConsumerModel(consumerModel))                      // 消费模式(集群消费:消费完其他人不能再读取/广播消费：所有人都能读)
	if err != nil {
		panic(err)
		// return nil
	}

	err = newPushConsumer.Start()
	if err != nil {
		panic(err)
		// return nil
	}

	mq := &MqCustomer{
		pushConsumer:   newPushConsumer,
		groupName:      groupName,
		isBroadCasting: isBroadCasting,
		MessageList:    make(chan *primitive.MessageExt, 1024),
	}
	//订阅
	return mq
}

// Subscribe 订阅者模式
func (mq *MqCustomer) Subscribe(topicName, tag string) {
	//tag = "broker-a"
	mq.subTopic = append(mq.subTopic, topicName)
	//订阅
	err := mq.pushConsumer.Subscribe(
		topicName,
		consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: tag, // 可以 TagA || TagB
		},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			ConcurrentlyCtx, _ := primitive.GetConcurrentlyCtx(ctx)
			log.Printf(" 开始读取数据: %v\n", ConcurrentlyCtx)
			for i := range msgs {
				select {
				case <-ctx.Done():
					return 0, nil
				case mq.MessageList <- msgs[i]:
					//fmt.Printf("订阅的回调: %v\n", msgs[i].Topic)
					//fmt.Println("Body", string(msgs[i].Body))
					continue
				}
			}
			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		fmt.Printf("Subscribe error:%s\n", err)
	}
}

// ShutdownPushConsumer  关闭消费者
func (mq *MqCustomer) ShutdownPushConsumer() error {

	err := mq.pushConsumer.Shutdown()
	if err != nil {
		fmt.Println("Shutdown Consumer error: ", err)
	}

	return err
}
