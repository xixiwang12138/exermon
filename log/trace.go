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
		ctx.Header(TraceHeader, traceId)
		start := time.Now()
		ctx.Set(TraceHeader, traceId)
		ctx.Next()
		end := time.Now()
		slog.InfoContext(ctx, "Access Log",
			slog.String("method", ctx.Request.Method),
			slog.String("url", ctx.Request.RequestURI),
			slog.Int64("cost", end.UnixMilli()-start.UnixMilli()))
	}
}
