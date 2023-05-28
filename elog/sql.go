package elog

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"time"
)

var GormLogger = impl{}

type impl struct{}

func (l *impl) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *impl) Info(ctx context.Context, s string, i ...interface{}) {
	cl := WithContext(ctx)
	cl.Info(s, i...)
}

func (l *impl) Warn(ctx context.Context, s string, i ...interface{}) {
	cl := WithContext(ctx)
	cl.Warning(s, i...)
}

func (l *impl) Error(ctx context.Context, s string, i ...interface{}) {
	cl := WithContext(ctx)
	cl.Error(s, i...)
}

func (l *impl) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.Info(ctx, fmt.Sprintf("[SQL] %s ==> (%dms) | rows: %d", sql, elapsed.Milliseconds(), rows))
}
