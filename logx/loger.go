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
