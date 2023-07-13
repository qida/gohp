package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type MQMessage struct {
	TopicName      string
	Body           []byte
	TagName        string
	Key            string
	DelayTimeLevel int
}

func NewMQMessage(topicName, tagName, key string, body []byte, delayTimeLevel int) *MQMessage {
	return &MQMessage{
		TopicName:      topicName,
		Body:           body,
		TagName:        tagName,
		Key:            key,
		DelayTimeLevel: delayTimeLevel,
	}
}

func (mqMessage *MQMessage) toMessage() *primitive.Message {

	//fmt.Println("Body2 ======>:", string(mqMessage.Body))
	message := primitive.NewMessage(mqMessage.TopicName, mqMessage.Body)

	//设置Tag
	if mqMessage.TagName != "" {
		message.WithTag(mqMessage.TagName)
	}

	//设置key
	if mqMessage.Key != "" {
		message.WithKeys([]string{mqMessage.Key})
	}

	//设置延迟时间
	if mqMessage.DelayTimeLevel > 0 {
		message.WithDelayTimeLevel(mqMessage.DelayTimeLevel)
	}

	return message
}
