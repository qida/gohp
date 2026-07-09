package logx

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// L 全局日志实例
var L *Loger

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
)

// Config 日志配置
type Config struct {
	Mode       string `mapstructure:"mode" yaml:"mode"` // 日志模式，可选值：console, file
	Level      string `mapstructure:"level" yaml:"level"`
	Filename   string `mapstructure:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress"`
}

// InitLoger 初始化日志
// mode: "console" - 仅输出到控制台, "file" - 仅输出到文件, 其他值 - 同时输出到控制台和文件
func InitLoger(cfg *Config) (*Loger, error) {
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
	case "mail":
	case "kafka":
	case "rocketmq":

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

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	L = &Loger{Logger: logger}
	return L, nil
}

// getEncoder 获取编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Debug Debug级别格式化日志
func (l *Loger) Debug(args ...any) {
	l.Sugar().Debug(args...)
}

// Info Info级别格式化日志
func (l *Loger) Info(args ...any) {
	l.Sugar().Info(args...)
}

// Error Error级别格式化日志
func (l *Loger) Error(args ...any) {
	l.Sugar().Error(args...)
}

// Panic Panic级别格式化日志
func (l *Loger) Panic(args ...any) {
	l.Sugar().Panic(args...)
}

// Fatal Fatal级别格式化日志
func (l *Loger) Fatal(args ...any) {
	l.Sugar().Fatal(args...)
}

// DPanic Panic级别格式化日志
func (l *Loger) DPanic(args ...any) {
	l.Sugar().DPanic(args...)
}

// Debugf Debug级别格式化日志
func (l *Loger) Debugf(format string, args ...any) {
	l.Sugar().Debugf(format, args...)
}

// Infof Info级别格式化日志
func (l *Loger) Infof(format string, args ...any) {
	l.Sugar().Infof(format, args...)
}

// Errorf Error级别格式化日志
func (l *Loger) Errorf(format string, args ...any) {
	l.Sugar().Errorf(format, args...)
}

// Panicf Panic级别格式化日志
func (l *Loger) Panicf(format string, args ...any) {
	l.Sugar().Panicf(format, args...)
}

// Fatalf Fatal级别格式化日志
func (l *Loger) Fatalf(format string, args ...any) {
	l.Sugar().Fatalf(format, args...)
}

// DPanicf Panic级别格式化日志
func (l *Loger) DPanicf(format string, args ...any) {
	l.Sugar().DPanicf(format, args...)
}

// Debugln Debug级别格式化日志
func (l *Loger) Debugln(args ...any) {
	l.Sugar().Debugln(args...)
}

// Infoln Info级别格式化日志
func (l *Loger) Infoln(args ...any) {
	l.Sugar().Infoln(args...)
}

// Errorln Error级别格式化日志
func (l *Loger) Errorln(args ...any) {
	l.Sugar().Errorln(args...)
}

// Panicln Panic级别格式化日志
func (l *Loger) Panicln(args ...any) {
	l.Sugar().Panicln(args...)
}

// Fatalln Fatal级别格式化日志
func (l *Loger) Fatalln(args ...any) {
	l.Sugar().Fatalln(args...)
}

// DPanicln Panic级别格式化日志
func (l *Loger) DPanicln(args ...any) {
	l.Sugar().DPanicln(args...)
}
