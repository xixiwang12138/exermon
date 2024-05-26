package integration

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	"github.com/xixiwang12138/exermon/xlog"
)

var GormXLoggerIntegration = impl{}

type impl struct{}

func (l *impl) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *impl) Info(ctx context.Context, s string, i ...interface{}) {
	cl := xlog.WithContext(ctx)
	cl.Info(s, i...)
}

func (l *impl) Warn(ctx context.Context, s string, i ...interface{}) {
	cl := xlog.WithContext(ctx)
	cl.Warn(s, i...)
}

func (l *impl) Error(ctx context.Context, s string, i ...interface{}) {
	cl := xlog.WithContext(ctx)
	cl.Error(s, i...)
}

func (l *impl) Debug(ctx context.Context, s string, i ...interface{}) {
	cl := xlog.WithContext(ctx)
	cl.Debug(s, i...)
}

func (l *impl) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.Error(ctx, "[SQL] %s ==> (%dms) | rows: %d | error: %s", sql, elapsed.Milliseconds(), rows, err.Error())
		return
	}
	l.Debug(ctx, fmt.Sprintf("[SQL] %s ==> (%dms) | rows: %d", sql, elapsed.Milliseconds(), rows))
}
