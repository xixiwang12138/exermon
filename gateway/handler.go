package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/exermon/elog"
	"net/http"
)

var (
	responseHandler = processResponse
)

func processResponse(ctx *gin.Context, err error, res any) {
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": res,
		})
	} else {
		cl := elog.WithTraceId(getTraceId(ctx))
		cl.Error("Server process err: %+v", err)
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
	}
}

// RawHandler 只包含逻辑调用的handler，只关心参数类型，因为任意响应类型最终都会使用json序列化
type RawHandler[P any, R any] func(c *gin.Context, ctx context.Context, p *P) (*R, error)

// RawNoParamHandler 只包含逻辑调用的无参数handler
type RawNoParamHandler[R any] func(c *gin.Context, ctx context.Context) (*R, error)

// Handler 带参数的handler包裹器
func Handler[P any, R any](rawFunc RawHandler[P, R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		var res any
		var err error
		if err = ctx.ShouldBind(param); err == nil {
			res, err = rawFunc(ctx, context.WithValue(ctx, elog.TraceIdHeader, getTraceId(ctx)), param)
		}
		responseHandler(ctx, err, res)
	}
}

// NoParamHandler 无参handler包裹器
func NoParamHandler[R any](rawFunc RawNoParamHandler[R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res any
		var err error

		res, err = rawFunc(ctx, context.WithValue(ctx, elog.TraceIdHeader, getTraceId(ctx)))
		responseHandler(ctx, err, res)
	}
}

func getTraceId(ctx *gin.Context) string {
	traceId, ok := ctx.Get(elog.TraceIdHeader)
	if !ok {
		traceId = elog.DefaultTrace
	}
	return traceId.(string)
}
