package elog

import (
	"context"
	"github.com/xixiwang12138/exermon/conf"
	"log"
)

var (
	base         *BaseLogger
	silentLogger *BaseLogger
)

// Setup level可配置: DEBUG, INFO, WARN, ERROR, FATAL
func Setup(cf *conf.LogConfig) {
	var err error
	base, err = newBaseLogger(logNameToLevel[cf.Level], cf.DirPath)
	if err != nil {
		log.Fatal("init log context error", err.Error())
	}

	silentLogger, err = newBaseLogger(SILENT, cf.DirPath)
	if err != nil {
		log.Fatal("init log context error", err.Error())
	}
}

func WithContext(ctx context.Context) *Logger {
	v := ctx.Value(TraceIdHeader)
	s := "default"
	if _, ok := v.(string); ok {
		s = v.(string)
	}

	v2 := ctx.Value(TraceLogEnableHeader)
	logEnable := enable
	if _, ok := v2.(int); ok {
		logEnable = v2.(int)
	}

	if logEnable == silence {
		return &Logger{
			BaseLogger: silentLogger,
			traceId:    s,
		}
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

const enable = 0
const silence = 1

func SilenceLogCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceLogEnableHeader, silence)
}

func EnableLogCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceLogEnableHeader, enable)
}

func TempSilence(ctx context.Context, f func()) context.Context {
	ctx = SilenceLogCtx(ctx)
	f()
	ctx = EnableLogCtx(ctx)
	return ctx
}
