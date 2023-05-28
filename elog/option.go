package elog

import (
	"context"
	"log"
	"os"
)

func init() {
	log.SetOutput(WithTraceId("default"))
	base = &BaseLogger{
		level:  DEBUG,
		logDir: "./log",
	}
	prefix, _ = os.Getwd()
	prefixLen = len(prefix)
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
