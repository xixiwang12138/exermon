package elog

import (
	"context"
	"github.com/xixiwang12138/exermon/conf"
	"log"
)

var (
	base *BaseLogger
)

// Setup level可配置: DEBUG, INFO, WARN, ERROR, FATAL
func Setup(cf *conf.LogConfig) {
	l, err := newBaseLogger(logNameToLevel[cf.Level], cf.DirPath)
	if err != nil {
		log.Fatal("init log context error", err.Error())
	}
	base = l
}

func WithContext(ctx context.Context) *Logger {
	v := ctx.Value(TraceIdHeader)
	s := "default"
	if _, ok := v.(string); ok {
		s = v.(string)
	}
	return &Logger{
		BaseLogger: base,
		traceId:    s,
	}
}

func WithTraceId(traceId string) *Logger {
	return &Logger{
		BaseLogger: base,
		traceId:    traceId,
	}
}
