package elog

import (
	"context"
	"log"
)

var (
	base *BaseLogger
)

func SetConfig(level LogLevel, logDir string) {
	l, err := newBaseLogger(level, logDir)
	if err != nil {
		log.Fatal("init log context error", err.Error())
	}
	base = l
}

func WithContext(ctx context.Context) *Logger {
	return &Logger{
		BaseLogger: base,
		traceId:    ctx.Value(TraceIdHeader).(string),
	}
}
