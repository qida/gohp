package logx

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// L 全局日志实例
var l *Loger

// Loger 日志封装
type Loger struct {
	*zap.Logger
}

// Level 日志级别
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	LevelPanic Level = "panic"
)

// Config 日志配置
type Config struct {
	Mode       string         `mapstructure:"mode" yaml:"mode"` // 日志模式，可选值：console, file, dingding, mail, kafka, rocketmq
	Level      string         `mapstructure:"level" yaml:"level"`
	Filename   string         `mapstructure:"filename" yaml:"filename"`
	MaxSize    int            `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int            `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int            `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool           `mapstructure:"compress" yaml:"compress"`
	Dingding   DingdingConfig `mapstructure:"dingding" yaml:"dingding"`
	Mail       MailConfig     `mapstructure:"mail" yaml:"mail"`
	Kafka      KafkaConfig    `mapstructure:"kafka" yaml:"kafka"`
	RocketMQ   RocketMQConfig `mapstructure:"rocketmq" yaml:"rocketmq"`
}

// DingdingConfig 钉钉日志配置
type DingdingConfig struct {
	Secret      string `mapstructure:"secret" yaml:"secret"`
	AccessToken string `mapstructure:"access_token" yaml:"access_token"`
}

// MailConfig 邮件日志配置
type MailConfig struct {
	From     string `mapstructure:"from" yaml:"from"`
	To       string `mapstructure:"to" yaml:"to"`
	Subject  string `mapstructure:"subject" yaml:"subject"`
	Smtp     string `mapstructure:"smtp" yaml:"smtp"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
}

// KafkaConfig Kafka日志配置
type KafkaConfig struct {
	Topic   string   `mapstructure:"topic" yaml:"topic"`
	Address []string `mapstructure:"address" yaml:"address"`
}

// RocketMQConfig RocketMQ日志配置
type RocketMQConfig struct {
	Topic      string   `mapstructure:"topic" yaml:"topic"`
	NameServer []string `mapstructure:"name_server" yaml:"name_server"`
	Group      string   `mapstructure:"group" yaml:"group"`
}

func init() {
	l = &Loger{Logger: zap.New(zapcore.NewCore(getEncoder(),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)}
}

// Init 初始化日志
// mode: "console" - 仅输出到控制台, "file" - 仅输出到文件, 其他值 - 同时输出到控制台和文件
func Init(cfg *Config) error {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoder := getEncoder()
	var core zapcore.Core

	switch cfg.Mode {
	case "console":
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	case "file":
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		})
		core = zapcore.NewCore(encoder, writeSyncer, level)

	case "dingding":
		w, err := getDingdingWriter(cfg.Dingding.Secret, cfg.Dingding.AccessToken)
		if err != nil {
			return err
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(w), level)
	case "mail":
		lvl := Level(cfg.Level)
		w, err := getMailWriter(&lvl, cfg.Mail.Port, cfg.Mail.From, cfg.Mail.To, cfg.Mail.Subject, cfg.Mail.Smtp, cfg.Mail.Password)
		if err != nil {
			return err
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(w), level)
	case "kafka":
		w, err := getKafkaWriter(cfg.Kafka.Topic, cfg.Kafka.Address)
		if err != nil {
			return err
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(w), level)
	case "rocketmq":
		w, err := getRocketmqWriter(cfg.RocketMQ.Topic, cfg.RocketMQ.NameServer, cfg.RocketMQ.Group)
		if err != nil {
			return err
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(w), level)
	default:
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		})
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		)
	}
	l = &Loger{Logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))}
	return nil
}

func GetLoger() *Loger {
	return l
}

// getEncoder 获取编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Debug Debug级别格式化日志
func Debug(args ...any) {
	l.Sugar().Debug(args...)
}

// Info Info级别格式化日志
func Info(args ...any) {
	l.Sugar().Info(args...)
}

// Error Error级别格式化日志
func Error(args ...any) {
	l.Sugar().Error(args...)
}

// Panic Panic级别格式化日志
func Panic(args ...any) {
	l.Sugar().Panic(args...)
}

// Fatal Fatal级别格式化日志
func Fatal(args ...any) {
	l.Sugar().Fatal(args...)
}

// DPanic Panic级别格式化日志
func DPanic(args ...any) {
	l.Sugar().DPanic(args...)
}

// Debugf Debug级别格式化日志
func Debugf(format string, args ...any) {
	l.Sugar().Debugf(format, args...)
}

// Infof Info级别格式化日志
func Infof(format string, args ...any) {
	l.Sugar().Infof(format, args...)
}

// Errorf Error级别格式化日志
func Errorf(format string, args ...any) {
	l.Sugar().Errorf(format, args...)
}

// Panicf Panic级别格式化日志
func Panicf(format string, args ...any) {
	l.Sugar().Panicf(format, args...)
}

// Fatalf Fatal级别格式化日志
func Fatalf(format string, args ...any) {
	l.Sugar().Fatalf(format, args...)
}

// DPanicf Panic级别格式化日志
func DPanicf(format string, args ...any) {
	l.Sugar().DPanicf(format, args...)
}

// Debugln Debug级别格式化日志
func Debugln(args ...any) {
	l.Sugar().Debugln(args...)
}

// Infoln Info级别格式化日志
func Infoln(args ...any) {
	l.Sugar().Infoln(args...)
}

// Errorln Error级别格式化日志
func Errorln(args ...any) {
	l.Sugar().Errorln(args...)
}

// Panicln Panic级别格式化日志
func Panicln(args ...any) {
	l.Sugar().Panicln(args...)
}

// Fatalln Fatal级别格式化日志
func Fatalln(args ...any) {
	l.Sugar().Fatalln(args...)
}

// DPanicln Panic级别格式化日志
func DPanicln(args ...any) {
	l.Sugar().DPanicln(args...)
}
