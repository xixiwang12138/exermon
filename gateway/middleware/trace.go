package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xixiwang12138/exermon/elog"
	"time"
)

func TracingLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(elog.TraceIdHeader)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		start := time.Now()
		ctx.Set(elog.TraceIdHeader, traceId)
		ctx.Next()
		ctx.Header(elog.TraceIdHeader, traceId)
		end := time.Now()

		cl := elog.WithTraceId(traceId)
		cl.Info("%-12s %-30s  ===>  %dms\n", ctx.Request.Method, ctx.Request.RequestURI, end.UnixMilli()-start.UnixMilli())
	}
}
