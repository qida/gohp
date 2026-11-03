# mqtt

消息队列封装包：本目录提供 Kafka 消费/发送封装；子目录 `rocketmq/` 提供 RocketMQ 生产者、消费者与管理端封装（见 `rocketmq/README.md`）。

## 功能列表

### kafka.go — Kafka 封装（基于 segmentio/kafka-go）
- `NewKAFKA(config kafka.ReaderConfig) *server`：创建 Kafka 服务实例
- `server.Listen(time_reconn)`：启动消费循环（goroutine 运行，不阻塞；从当前时间前 1 小时的偏移开始读），异常退出后按 `time_reconn` 间隔自动重连
- `server.OnNewClient(callback)`：注册消费启动回调
- `server.OnNewMessage(callback)`：注册收到消息回调
- `server.OnClientConnectionClosed(callback)`：注册消费异常退出回调
- `ClientKafka.Send(key, msg)`：向配置的 topic 同步写入消息
- `ClientKafka.Writer`：暴露底层 `*kafka.Writer`（LeastBytes 分区均衡）
- 收到 SIGINT/SIGTERM 系统信号时消费循环优雅退出

### 子包 rocketmq
- 生产者（同步/异步/批量/Tag/延迟消息）、Push 消费者（集群/广播模式）、Topic 管理端，详见 `rocketmq/README.md`

## 使用示例

```go
package main

import (
	"fmt"
	"time"

	"github.com/qida/gohp/mqtt"
	"github.com/segmentio/kafka-go"
)

func main() {
	s := mqtt.NewKAFKA(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "test",
		GroupID: "g1",
	})
	s.OnNewMessage(func(c *mqtt.ClientKafka, message kafka.Message) {
		fmt.Println("recv:", string(message.Value))
	})
	s.Listen(3 * time.Second) // 消费异常退出后 3 秒重连
}
```

## 备注

- 外部依赖：github.com/segmentio/kafka-go，以及本仓库 `logx` 包
- 需要可访问的 Kafka broker；`Listen` 内部已捕获 panic，仅记录日志
