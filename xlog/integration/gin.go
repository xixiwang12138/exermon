package integration

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xixiwang12138/exermon/elog"
	"github.com/xixiwang12138/exermon/xlog"
)

var GinXLoggerIntegration = TracingLogger()

func TracingLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(elog.TraceIdHeader)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		logger := xlog.NewLogger(traceId)
		ctx.Set(xlog.XLOG, logger)
		start := time.Now()
		ctx.Next()
		ctx.Header(xlog.XLOG, traceId)
		end := time.Now()
		xlog.WithContext(ctx).Info("%-12s %-30s  ===>  %dms\n", ctx.Request.Method, ctx.Request.RequestURI, end.UnixMilli()-start.UnixMilli())
	}
}
