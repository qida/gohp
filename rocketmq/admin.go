package rocketmq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type RockerMQAdmin struct {
	NameSrvAddr  []string
	BrokerAddr   string
	Admin        admin.Admin
	Topic        string
	defaultRetry int
}

// NewMQAdmin 创建MQAdmin
func NewMQAdmin(Topic, BrokerAddr string, nameSrvAddr []string) *RockerMQAdmin {

	Admin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)))
	if err != nil {
		panic(err)
	}
	return &RockerMQAdmin{
		NameSrvAddr:  nameSrvAddr,
		BrokerAddr:   BrokerAddr,
		Admin:        Admin,
		Topic:        Topic,
		defaultRetry: 0,
	}
}

// CloseMQAdmin 关闭MQAdmin
func (r *RockerMQAdmin) CloseMQAdmin() error {
	err := r.Admin.Close()
	if err != nil {
		fmt.Println("Create topic error:", err)
	}

	return err
}

// CreateTopic 创建topic
func (r *RockerMQAdmin) CreateTopic(TopicName string) error {
	if TopicName == "" {
		TopicName = r.Topic
	}

	err := r.Admin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(TopicName),
		admin.WithBrokerAddrCreate(r.BrokerAddr))
	if err != nil {
		fmt.Println("Create topic error:", err)
	}
	return err
}

// DeleteTopic 删除topic
func (r *RockerMQAdmin) DeleteTopic(TopicName string) error {
	err := r.Admin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete(TopicName),
		//admin.WithBrokerAddrDelete(brokerAddr),
		//admin.WithNameSrvAddr(nameSrvAddr),
	)

	return err
}
