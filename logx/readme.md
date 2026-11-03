# logx

基于 zap 封装的日志包，支持控制台、文件、钉钉、邮件、Kafka、RocketMQ 多种输出目标。

## 功能列表

### logx.go — 核心封装
- `Loger`：日志器，内嵌 `*zap.Logger`；`Level` 与 `LevelDebug`~`LevelPanic` 级别常量
- `Config`（含 `DingdingConfig`/`MailConfig`/`KafkaConfig`/`RocketMQConfig`）：日志配置，`mode` 可选 console/file/dingding/mail/kafka/rocketmq，其他值时同时输出控制台+文件
- `DefaultWithMap(cfg *Config) error`：按配置初始化全局日志（未初始化时默认 Info 级控制台输出）
- `GetLoger() *Loger`：获取全局日志实例
- `Debug/Info/Error/Panic/Fatal/DPanic` 及 `f`/`ln` 后缀变体：包级日志函数

### dingtalk.go — 钉钉机器人
- `Robot` / `NewRobot(token, secret)`：钉钉群机器人，支持加签
- `SendMessage(msg)`：发送任意消息体；`SendTextMessage`/`SendMarkdownMessage`/`SendLinkMessage`：发送 text/markdown/link 消息

### dingding.go — 钉钉日志 writer
- `LogDingding`：实现 `io.Writer`，把 zap 日志格式化为钉钉 markdown 消息发送（复用 dingtalk.go 的 `Robot`）

### mail.go — 邮件日志 writer
- `Mail`：实现 `io.Writer`，通过 gomail 异步发送日志邮件

### kafka.go — Kafka 日志 writer
- `LogKafka`：实现 `io.Writer`，通过 sarama 同步生产者把日志写入 Kafka

### rocketmq.go / rocket.go / rocket_mq.go — RocketMQ 日志
- `LogRocketmq`：实现 `io.Writer`，通过官方客户端同步发送日志
- `LogRocket` / `NewProbeLog(m map[string]interface{})`：数据探针日志，`MQDebug/MQInfo/MQWarn/MQError/MQPanic` 分级发送、`SetTagKey` 设置 Tag/Key（异步发送，复用 mqtt/rocketmq 的 MqProducer）
- `NewRocketMQCore(enc, topic, producer, enabler) zapcore.Core`：自定义 zap Core，将日志写入 RocketMQ

### gorm.go — GORM 适配
- `GormLogger` / `NewGormLogger(l, cfg)`：把 zap 日志适配为 gorm `logger.Interface`，支持慢查询告警、可忽略 record not found

### file.go — 文件日志
- `FileLogger`：lumberjack 文件日志配置（大小/份数/天数/压缩），内置 cron 每日轮转

### log_MQ_type.go — MQ 日志结构
- `LogHeader`/`LogBody`/`LogQuery`/`Base`/`LogStatistic`/`LogField`：日志采集与查询数据结构
- `LogAny(key, val)`：构造 `LogField`

## 使用示例

```go
package main

import "github.com/qida/gohp/logx"

func main() {
    err := logx.DefaultWithMap(&logx.Config{
        Mode:     "file",
        Level:    "info",
        Filename: "logs/app.log",
        MaxSize:  100, MaxBackups: 3, MaxAge: 7, Compress: true,
    })
    if err != nil {
        panic(err)
    }
    logx.Infof("服务启动 addr=%s", ":8080")
}
```

## 备注

- 外部依赖：zap、lumberjack、sarama(IBM)、rocketmq-client-go、gomail、robfig/cron、gorm、goccy/go-json、spf13/cast、pkg/errors，以及本仓库 `mqtt/rocketmq`、`slice` 包
- dingding 模式需可访问钉钉 OpenAPI，mail 模式需可访问 SMTP 服务器，kafka/rocketmq 模式需对应集群可用
- mail 每条日志发一封邮件，建议配合较高 `level`（如 error）使用
