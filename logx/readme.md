# logx

基于 [zap](https://github.com/uber-go/zap) 封装的日志包，支持控制台、文件、钉钉、邮件、Kafka、RocketMQ 多种输出目标。

## 特性

- 包级函数直接调用：`logx.Info`、`logx.Infof`、`logx.Errorf` 等，无需持有实例
- `init()` 默认创建一个 Info 级别的控制台 logger，未调用 `Init` 也能安全使用
- 通过 `Init` 按配置切换输出目标与级别
- 文件模式基于 lumberjack，支持按大小切割、按份数/天数清理、gzip 压缩
- 钉钉/邮件/Kafka/RocketMQ 作为日志落地目标

## 初始化

```go
import "github.com/qida/gohp/logx"

// 通常在程序启动时调用一次，传入从配置文件解析出的 *logx.Config
if err := logx.Init(&cfg.Loger); err != nil {
    panic(err)
}

// 也可以拿到底层 *Loger（内嵌 *zap.Logger）做进一步定制
l := logx.GetLoger()
```

## 日志输出

```go
logx.Debug("调试信息")
logx.Info("启动成功", zap.String("addr", addr))
logx.Errorf("请求失败: %v", err)
logx.Infof("Swagger文档地址: http://localhost%s%s", addr, "/swagger")
logx.Fatal("致命错误", zap.Error(err))
```

可用包级函数（与 zap 风格一致）：

| 类型 | 函数 |
| --- | --- |
| `args ...any` | `Debug` `Info` `Error` `Panic` `Fatal` `DPanic` |
| `format` 格式化 | `Debugf` `Infof` `Errorf` `Panicf` `Fatalf` `DPanicf` |
| `ln` 换行 | `Debugln` `Infoln` `Errorln` `Panicln` `Fatalln` `DPanicln` |

## 配置项

`Config` 结构对应 yaml 中的 `loger` 节点：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `mode` | string | 输出模式，见下表 |
| `level` | string | 日志级别：`debug` `info` `warn` `error` `fatal`，非法值回退为 `info` |
| `filename` | string | 文件路径（`file` / 默认模式使用） |
| `max_size` | int | 单文件最大大小，单位 MB（文件切割） |
| `max_backups` | int | 最大保留旧文件份数 |
| `max_age` | int | 最大保留天数 |
| `compress` | bool | 是否 gzip 压缩旧文件 |
| `dingding` | object | 钉钉配置，见下 |
| `mail` | object | 邮件配置，见下 |
| `kafka` | object | Kafka 配置，见下 |
| `rocketmq` | object | RocketMQ 配置，见下 |

`mode` 取值：

| mode | 输出目标 | 说明 |
| --- | --- | --- |
| `console` | 控制台 | 仅 stdout |
| `file` | 文件 | 仅落盘，按 `filename`/`max_size` 等切割 |
| `dingding` | 钉钉群机器人 | 每条日志发送一条钉钉消息 |
| `mail` | 邮件 | 每条日志发送一封邮件 |
| `kafka` | Kafka | 每条日志作为消息写入指定 topic |
| `rocketmq` | RocketMQ | 每条日志作为消息写入指定 topic |
| 其他/留空 | 文件 + 控制台 | 同时落盘与输出到 stdout |

子配置字段：

```yaml
dingding:
  secret: ""          # 钉钉机器人加签密钥
  access_token: ""    # 机器人 access_token

mail:
  from: ""            # 发件人邮箱
  to: ""              # 收件人邮箱
  subject: ""         # 邮件主题
  smtp: ""            # SMTP 服务器地址
  port: 465           # SMTP 端口
  password: ""        # 发件人邮箱密码/授权码

kafka:
  topic: ""           # 目标 topic
  address: []         # Kafka broker 地址列表

rocketmq:
  topic: ""           # 目标 topic
  name_server: []     # NameServer 地址列表
  group: ""           # 生产者 group
```

## 配置示例

### 控制台输出（开发环境常用）

```yaml
loger:
  mode: "console"
  level: "debug"
```

### 文件输出

```yaml
loger:
  mode: "file"
  level: "info"
  filename: "logs/app.log"
  max_size: 100        # 单文件 100MB
  max_backups: 3       # 保留 3 份历史
  max_age: 7           # 保留 7 天
  compress: true       # 历史文件 gzip 压缩
```

### 文件 + 控制台（默认模式，生产环境推荐）

`mode` 留空或设为未识别值时，日志同时写入文件与控制台：

```yaml
loger:
  mode: ""
  level: "info"
  filename: "logs/app.log"
  max_size: 100
  max_backups: 3
  max_age: 7
  compress: true
```

### 钉钉告警

适合把 `error` 及以上级别的日志推送到钉钉群：

```yaml
loger:
  mode: "dingding"
  level: "error"
  dingding:
    secret: "SECxxxxxxxxxxxxxxxx"
    access_token: "xxxxxxxxxxxxxxxx"
```

### 邮件告警

```yaml
loger:
  mode: "mail"
  level: "error"
  mail:
    from: "alert@xx.com"
    to: "ops@xx.com"
    subject: "系统告警"
    smtp: "smtp.xx.com"
    port: 465
    password: "********"
```

### Kafka

```yaml
loger:
  mode: "kafka"
  level: "info"
  kafka:
    topic: "app_log"
    address:
      - "127.0.0.1:9092"
```

### RocketMQ

```yaml
loger:
  mode: "rocketmq"
  level: "info"
  rocketmq:
    topic: "app_log"
    name_server:
      - "127.0.0.1:9876"
    group: "log_producer"
```

## 备注

- `dingding` / `mail` / `kafka` / `rocketmq` 为单目标输出，不会同时落盘。若需同时保留本地文件，请使用默认模式（`mode` 留空），或在外层用 `zapcore.NewTee` 自行组合
- `mail` 模式每条日志都会发邮件，高频日志可能造成邮件轰炸，建议配合较高的 `level`（如 `error`）使用
- 钉钉/邮件的发送为异步或同步由各 writer 实现决定，详见 `dingding.go` / `mail.go` / `kafka.go` / `rocketmq.go`
