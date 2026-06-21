package logx

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger 将 zap logger 适配为 gorm logger.Interface
type GormLogger struct {
	*Loger
	SlowThreshold             time.Duration
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建 GORM 日志适配器
func NewGormLogger(l *Loger, cfg logger.Config) *GormLogger {
	return &GormLogger{
		Loger:                     l,
		SlowThreshold:             cfg.SlowThreshold,
		LogLevel:                  cfg.LogLevel,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
	}
}

// LogMode 设置日志级别
func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return &newLogger
}

// Info 打印 Info 级别日志
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Info {
		g.Loger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 打印 Warn 级别日志
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		g.Loger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 打印 Error 级别日志
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.Loger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 打印 SQL 执行追踪日志
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && g.LogLevel >= logger.Error && (!g.IgnoreRecordNotFoundError || err.Error() != "record not found"):
		g.Loger.Error("SQL执行错误",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		g.Loger.Warn("SQL慢查询",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case g.LogLevel >= logger.Info:
		g.Loger.Info("SQL执行",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
