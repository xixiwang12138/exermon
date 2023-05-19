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
	v := ctx.Value(TraceIdHeader)
	s := "empty"
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
