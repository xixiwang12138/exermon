package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/xlog"
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
		xl := xlog.FromGin(ctx)
		xl.Error("server process err: ", err.Error())
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
	}
}

// RawHandler 只包含逻辑调用的handler，只关心参数类型，因为任意响应类型最终都会使用json序列化
type RawHandler[P any] func(c *gin.Context, xl *xlog.XLogger, p *P) (any, error)

// RawNoParamHandler 只包含逻辑调用的无参数handler
type RawNoParamHandler func(c *gin.Context, xl *xlog.XLogger) (any, error)

// Handler 带参数的handler包裹器
func Handler[P any](rawFunc RawHandler[P]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		var res any
		var err error
		if err = ctx.ShouldBind(param); err == nil {
			res, err = rawFunc(ctx, xlog.FromGin(ctx), param)
		}
		responseHandler(ctx, err, res)
	}
}

// NoParamHandler 无参handler包裹器
func NoParamHandler(rawFunc RawNoParamHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res any
		var err error
		res, err = rawFunc(ctx, xlog.FromGin(ctx))
		responseHandler(ctx, err, res)
	}
}
