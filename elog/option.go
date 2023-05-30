package elog

import (
	"context"
	"log"
)

func init() {
	base = &BaseLogger{
		level:  DEBUG,
		logDir: "./log",
	}
	log.SetOutput(WithTraceId("default"))
}

var (
	base      *BaseLogger
	prefix    string
	prefixLen int
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
