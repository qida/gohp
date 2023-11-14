package logx

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger .
type Logger struct {
	options    []zap.Option
	cores      []zapcore.Core
	log        *zap.Logger
	sugar      *zap.SugaredLogger
	Kafka      string
	File       string
	level      Level
	SourceItem string
}

// Message
type Message struct {
	Msg         string      `json:"msg"`
	Level       string      `json:"level"`
	TimeKey     string      `json:"time"`
	CallerKey   string      `json:"file"`
	ServiceName string      `json:"serviceName"`
	Location    string      `json:"location"`
	Other       interface{} `json:"detail"`
}

type Level int8

// ZapLevel .
func (level Level) ZapLevel() zapcore.Level {
	return zapcore.Level(int8(level))
}

// Log level
const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	MessageKey:  "msg",
	LevelKey:    "level",
	EncodeLevel: zapcore.CapitalLevelEncoder,
	TimeKey:     "time",
	EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	},
	CallerKey:      "file",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	//EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendInt64(int64(d) / 1000000)
	//},
}

var defaultEncoder = zapcore.NewJSONEncoder(defaultEncoderConfig)

// Options .
type Options func(l *Logger)

var defaultMu sync.Mutex
var defaultLogger *Logger
var once sync.Once

func init() {
	once.Do(func() {
		defaultLogger, _ = Default()
	})
}

// Default .
func Default() (*Logger, error) {
	return NewLog(
		WithLevel(DebugLevel),
		WithCaller(),
	)
}

func DefaultWithMap(m map[string]interface{}) error {
	l, err := NewWithMap(m)
	if err != nil {
		return err
	}
	defaultMu.Lock()
	defaultLogger = l
	defaultMu.Unlock()
	return nil
}

func InitLogZapMQ(path string) {

	jsonStr, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(jsonStr, &m)
	if err != nil {
		panic(err)
	}

	defaultLogger, err = NewWithMap(m)
	if err != nil {
		panic(err)
	}
}

// NewWithMap new with map
func NewWithMap(m map[string]interface{}) (*Logger, error) {
	opts := make([]Options, 0)

	for k, v := range m {
		_ = v
		switch k {
		case "console":
			csl, ok := v.(bool)
			if csl && ok {
				opts = append(opts, WithConsole())
			}
		case "dingding":
			lk, err := toDingding(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithDingding(lk.secret, lk.access_token, *lk.level, lk.fields...))
		case "mail":
			lms, err := toMail(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithMail(lms))
		case "kafka":
			lks, err := toKafka(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithKafka(lks))
		case "rocketmq":
			lrs, err := toRocketMQ(v)
			if err != nil {
				continue
			}

			opts = append(opts, WithRocketMQ(lrs))

		case "logs":
			lfs, err := toLogs(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithLogs(lfs))
		case "file":
			lfs, err := toFile(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithFile(lfs))
		case "level":
			level, err := toLevel(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithLevel(level))
		case "caller":
			caller, ok := v.(bool)
			if caller && ok {
				opts = append(opts, WithCaller())
			}
		case "skip":
			skip, ok := v.(int)
			if ok {
				opts = append(opts, WithSkip(skip))
			}
		case "stacktrace":
			level, err := toLevel(v)
			if err != nil {
				continue
			}
			opts = append(opts, WithStacktrace(level))
		case "initial":
			initialMap, ok := v.(map[string]string)
			if !ok {
				continue
			}
			key, ok := initialMap["key"]
			if !ok {
				continue
			}
			value, ok := initialMap["value"]
			if !ok {
				continue
			}
			opts = append(opts, WithInitialFields(key, value))
		default:
			// fmt.Println("没有这个参数: ", k)
		}
	}
	return NewLog(opts...)
}

// New  Logger
func NewLog(opts ...Options) (*Logger, error) {
	l := &Logger{
		cores:   make([]zapcore.Core, 0),
		options: make([]zap.Option, 0),
		level:   DebugLevel,
	}
	for _, o := range opts {
		o(l)
	}

	//cfg := zap.NewDevelopmentEncoderConfig()
	//cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	//}
	//cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//l.cores = append(l.cores,
	//	zapcore.NewCore(
	//		zapcore.NewConsoleEncoder(cfg),
	//		zapcore.AddSync(colorable.NewColorableStdout()), l.level.ZapLevel(),
	//	),
	//)
	core := zapcore.NewTee(l.cores...)

	l.options = append(l.options, zap.AddCallerSkip(1))

	l.log = zap.New(core, l.options...)
	l.sugar = l.log.Sugar()
	return l, nil
}

func WithConsole() Options {
	return func(l *Logger) {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(l.level.ZapLevel())
		//cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		build, err := cfg.Build()
		if err != nil {
			return
		}

		l.cores = append(l.cores,
			build.Core(),
			//zapcore.NewCore(
			//	zapcore.NewConsoleEncoder(cfg),
			//	zapcore.AddSync(colorable.NewColorableStdout()),
			//	l.level.ZapLevel(),
			//),
		)
	}
}

// WithRocketMQ 启用RocketMQ
func WithRocketMQ(rocketmqLogger []LogRocketMQ) Options {
	return func(l *Logger) {
		for _, mq := range rocketmqLogger {
			if mq.Level == nil {
				level := ErrorLevel
				mq.Level = &level
			}
			writer, err := getRocketMQWriter(mq.Topic, mq.Group, mq.Address, mq.fields...)
			if err != nil {
				continue
			}
			core := zapcore.NewCore(
				defaultEncoder,
				//zapcore.NewMultiWriteSyncer(getRocketMQMultiWriter(mq.Topic, mq.Group, mq.Address, 50, mq.fields...)...),
				zapcore.AddSync(writer),
				mq.Level.ZapLevel(),
			)
			if len(mq.fields) > 0 {
				core = core.With(mq.fields)
			}

			l.cores = append(l.cores, core)
		}
	}
}

// WithKafka 启用kafka后，日志将打入kafka
func WithKafka(kafkaLogger []LogKafka) Options {
	return func(l *Logger) {
		for _, k := range kafkaLogger {
			if k.Level == nil {
				level := ErrorLevel
				k.Level = &level
			}
			writer, err := getKafkaWriter(k.Topic, k.Address, k.fields...)
			if err != nil {
				continue
			}
			core := zapcore.NewCore(
				defaultEncoder,
				zapcore.AddSync(writer),
				k.Level.ZapLevel(),
			)
			if len(k.fields) > 0 {
				core = core.With(k.fields)
			}
			l.cores = append(l.cores, core)
		}
	}
}

// WithMail 发送email
// 默认Level为DPanic
func WithMail(mailLogger []Mail) Options {
	return func(l *Logger) {
		for _, m := range mailLogger {
			if m.Level == nil {
				level := DPanicLevel
				m.Level = &level
			}
			writer, err := getMailWriter(m.Level, m.Port, m.From, m.To, m.Subject, m.Stmp, m.Password, m.fields...)
			if err != nil {
				continue
			}
			core := zapcore.NewCore(
				defaultEncoder,
				zapcore.AddSync(writer),
				m.Level.ZapLevel(),
			)
			if len(m.fields) > 0 {
				core = core.With(m.fields)
			}
			l.cores = append(l.cores, core)
		}
	}
}

// WithFile 开启写入文件
func WithFile(fileLogger []FileLogger) Options {
	return func(l *Logger) {
		for _, f := range fileLogger {
			if f.Level == nil {
				level := InfoLevel
				f.Level = &level
			}
			writer, _ := getFileWriter(f)
			l.cores = append(l.cores,
				zapcore.NewCore(
					defaultEncoder,
					zapcore.AddSync(writer),
					f.Level.ZapLevel(),
				),
			)
		}
	}
}

// WithLogs 开启写入单个文件
func WithLogs(fileLogger FilesLogger) Options {
	return func(l *Logger) {
		writer, _ := getFilesWriter(fileLogger)
		var currentDefaultEncoder = defaultEncoder
		if fileLogger.IsUTC {
			defaultEncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(t.UTC().Format("2006-01-02 15:04:05.000"))
			}
			currentDefaultEncoder = zapcore.NewJSONEncoder(defaultEncoderConfig)
		}

		l.cores = append(l.cores,
			zapcore.NewCore(
				currentDefaultEncoder,
				zapcore.AddSync(writer),
				l.level.ZapLevel(),
			),
		)
	}
}

// WithLogs 开启写入单个文件 并且制定enCider
func WithLogsEncoder(fileLogger FilesLogger, enCoder zapcore.Encoder) Options {
	return func(l *Logger) {
		writer, err := getFilesWriter(fileLogger)
		if err != nil {
			return
		}
		l.cores = append(l.cores,
			zapcore.NewCore(
				enCoder,
				zapcore.AddSync(writer),
				l.level.ZapLevel(),
			),
		)
	}
}

func WithDingding(secret, access_token string, level Level, fields ...zap.Field) Options {
	return func(l *Logger) {
		writer, err := getDingdingWriter(secret, access_token, fields...)
		if err != nil {
			return
		}
		core := zapcore.NewCore(
			defaultEncoder,
			zapcore.AddSync(writer), level.ZapLevel(),
		)
		if len(fields) > 0 {
			core = core.With(fields)
		}
		l.cores = append(l.cores, core)
	}
}

// WithCaller 添加文件名、行号、函数名等
func WithCaller() Options {
	return func(l *Logger) {
		l.options = append(l.options, zap.AddCaller())
	}
}

func WithStacktrace(level Level) Options {
	return func(l *Logger) {
		l.options = append(l.options, zap.AddStacktrace(level.ZapLevel()))
	}
}

// WithLevel 指定日志级别
func WithLevel(level Level) Options {
	return func(l *Logger) {
		l.level = level
	}
}

func WithSkip(skip int) Options {
	return func(l *Logger) {
		l.options = append(l.options, zap.AddCallerSkip(skip))
	}
}

// WithInitialFields 固定字段
func WithInitialFields(key, value string) Options {
	return func(l *Logger) {
		l.options = append(l.options, zap.Fields(zap.String(key, value)))
	}
}

func WithSourceItem(item string) Options {
	return func(l *Logger) {
		l.SourceItem = item
	}
}

// AddSkip 临时增加跳转层级
func AddSkip(skip int) *Logger {
	return defaultLogger.AddSkip(skip)
}

// AddSkip 临时增加跳转层级
func (l *Logger) AddSkip(skip int) *Logger {
	cp := *l
	cp.log = l.log.WithOptions(zap.AddCallerSkip(skip))
	return &cp
}

// Stacktrace 临时开启堆栈输出
func Stacktrace() *Logger {
	return defaultLogger.Stacktrace()
}

// Stacktrace 临时开启堆栈输出
func (l *Logger) Stacktrace() *Logger {
	cp := *l
	cp.log = l.log.WithOptions(zap.AddStacktrace(zapcore.DebugLevel))
	return &cp
}

func Log(lvl Level, msg string, fields ...zap.Field) {
	defaultLogger.log.Log(lvl.ZapLevel(), msg, fields...)
}

func (l *Logger) Log(lvl Level, msg string, fields ...zap.Field) {
	l.log.Log(lvl.ZapLevel(), msg, fields...)
}

// Debug .
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.log.Debug(msg, fields...)
}

// Debug .
// 不同写入位置可能区分级别也不一样
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg)
}

// Debugf .
func Debugf(msg string, args ...interface{}) {
	defaultLogger.sugar.Debugf(msg, args...)
}

// Debugf .
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.sugar.Debugf(msg, args...)
}

// Info .
func Info(msg string, fields ...zap.Field) {
	defaultLogger.log.Info(msg, fields...)
}

// Info .
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

// Infof .
func Infof(msg string, args ...interface{}) {
	defaultLogger.sugar.Infof(msg, args...)
}

// Infof .
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.sugar.Infof(msg, args...)
}

// Warn .
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.log.Warn(msg, fields...)
}

// Warn .
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

// Warnf .
func Warnf(msg string, args ...interface{}) {
	defaultLogger.sugar.Warnf(msg, args...)
}

// Wranf .
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.sugar.Warnf(msg, args...)
}

// Error .
func Error(msg string, fields ...zap.Field) {
	defaultLogger.log.Error(msg, fields...)
}

// Error .
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

// Errorf .
func Errorf(msg string, args ...interface{}) {
	defaultLogger.sugar.Errorf(msg, args...)
}

// Errorf .
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.sugar.Errorf(msg, args...)
}

// DPanic .
func DPanic(msg string, fields ...zap.Field) {
	defaultLogger.log.DPanic(msg, fields...)
}

// DPanic .
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.log.DPanic(msg, fields...)
}

// DPanicf .
func DPanicf(msg string, args ...interface{}) {
	defaultLogger.sugar.DPanicf(msg, args...)
}

// DPanicf .
func (l *Logger) DPanicf(msg string, args ...interface{}) {
	l.sugar.DPanicf(msg, args...)
}

// Fatal .
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.log.Fatal(msg, fields...)
}

// Fatal .
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

// Fatalf .
func Fatalf(msg string, args ...interface{}) {
	defaultLogger.sugar.Fatalf(msg, args...)
}

// Fatalf .
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.sugar.Fatalf(msg, args...)
}

/***************************loggerMQ -- 数据采集***********************************/

// MQInfoElapsedTime 统计耗时
func (l *Logger) MQInfoElapsedTime(title string, fields ...zap.Field) {
	//fields = append(fields, zap.String("trans", trans), zap.String("method", method), zap.Duration("costTime", costTime))
	l.log.Info(title, fields...)
}
