# rocketmq

Apache RocketMQ 客户端封装，提供生产者、Push 消费者、Topic 管理端与统一消息定义。

## 功能列表

### producer.go — 生产者封装
- `MqProducer` / `NewMQProducer(topicName, groupName, NameSrvAddr)`：创建并启动生产者（自动创建 topic key、重试 5 次，内置 ants 协程池）
- `SendSyncMessage(msg)`：同步发送单条消息
- `SendSyncbatchMessage(msgs)`：同步批量发送
- `SendAsyncMessage(msg, callback)`：异步发送并回调（发送失败时以 3s 超时 ctx 回调错误）
- `SendAsyncMessageBody(msg)`：异步发送消息体（仅 Topic+Body）
- `Shutdown()`：关闭生产者并释放协程池

### message.go — 消息定义
- `MQMessage` / `NewMQMessage(topic, tag, key, body, delayTimeLevel)`：统一消息结构，支持 Tag、Key、延迟级别

### customer.go — Push 消费者
- `MqCustomer` / `NewPushCustomer(groupName, NameSrvAddr, isBroadCasting)`：创建 Push 消费者（集群/广播消费模式，从首次偏移开始消费）
- `Subscribe(topicName, tag)`：按 Tag 订阅 topic，收到的消息送入 `MessageList` channel（缓冲 1024）
- `ShutdownPushConsumer()`：关闭消费者

### admin.go — Topic 管理
- `RockerMQAdmin` / `NewMQAdmin(topic, brokerAddr, nameSrvAddr)`：创建管理端
- `CreateTopic(name)` / `DeleteTopic(name)`：创建/删除 topic
- `CloseMQAdmin()`：关闭管理端

### mq_send.go — 简易发送器
- `SendMsg` / `ConnRockerMQ(ip, port, topic, group)`：连接 NameServer 并启动生产者
- `SendSync` / `SendAsync` / `SendAsyncBatch`：同步/异步/批量发送
- `DelaySendSync`：发送延迟消息（DelayTimeLevel=3）
- `TagSendSync(tag, body)`：发送带 Tag 消息
- `Close()`：关闭生产者

## 使用示例

```go
p := rocketmq.NewMQProducer("test_topic", "my_group", []string{"127.0.0.1:9876"})
defer p.Shutdown()

msg := rocketmq.NewMQMessage("test_topic", "tag_a", "key_1", []byte("hello"), 0)
p.SendSyncMessage(msg)
```

## 备注

- 外部依赖：github.com/apache/rocketmq-client-go/v2、github.com/panjf2000/ants/v2
- 需要可访问的 RocketMQ NameServer；`NewMQProducer`/`NewPushCustomer` 初始化失败时直接 panic
- `SendMsg.SendAsync` 源码注释标注存在高并发 bug，慎用
