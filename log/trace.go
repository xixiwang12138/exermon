package log

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TracingLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader(TraceHeader)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		start := time.Now()
		ctx.Set(TraceHeader, traceId)
		ctx.Next()
		ctx.Header(TraceHeader, traceId)
		end := time.Now()
		slog.InfoContext(ctx, "%-12s %-30s  ===>  %dms", ctx.Request.Method, ctx.Request.RequestURI, end.UnixMilli()-start.UnixMilli())
	}
}
