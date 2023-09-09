package logger

import (
	"context"
	"errors"
	"gohub/pkg/helpers"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 操作对象，实现 gormlogger.Interface
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

// NewGormLogger 外部调用。实例化一个 GormLogger 对象，示例：
//     DB, err := gorm.Open(dbConfig, &gorm.Config{
//         Logger: logger.NewGormLogger(),
//     })
func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     Logger,                 // 使用全局的 logger.Logger 对象
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值，单位为千分之一秒
	}
}

// LogMode 实现 gormlogger.Interface 的 LogMode 方法
func (log GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		ZapLogger:     log.ZapLogger,
		SlowThreshold: log.SlowThreshold,
	}
}

// Info 实现 gormlogger.Interface 的 Info 方法
func (log GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	log.logger().Sugar().Debugf(str, args...)
}

// Warn 实现 gormlogger.Interface 的 Warn 方法
func (log GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	log.logger().Sugar().Warnf(str, args...)
}

// Error 实现 gormlogger.Interface 的 Error 方法
func (log GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	log.logger().Sugar().Errorf(str, args...)
}

// Trace 实现 gormlogger.Interface 的 Trace 方法
func (log GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	// 通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", helpers.MicrosecondsStr(elapsed)),
		zap.Int64("rows", rows),
	}

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			// 其他错误使用 error 等级
			logFields = append(logFields, zap.Error(err))
			log.logger().Error("Database Error", logFields...)
		}
	}

	// 慢查询日志
	if log.SlowThreshold != 0 && elapsed > log.SlowThreshold {
		log.logger().Warn("Database Slow Log", logFields...)
	}

	// 记录所有 SQL 请求
	log.logger().Debug("Database Query", logFields...)
}

// logger 内用的辅助方法，确保 Zap 内置信息 Caller 的准确性（如 paginator/paginator.go:148）
func (log GormLogger) logger() *zap.Logger {

	// 跳过 gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// 减去一次封装，以及一次在 logger 初始化里添加 zap.AddCallerSkip(1)
	clone := log.ZapLogger.WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return log.ZapLogger
}
