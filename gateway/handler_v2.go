package gateway

import (
	"context"
	"github.com/xixiwang12138/exermon/auth"
	"github.com/xixiwang12138/exermon/elog"

	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/exermon/errors"
)

type HandlerFunc[P any, R any] func(c context.Context, req *P) R

func wrapGinContext(ctx *gin.Context) context.Context {
	c := context.WithValue(context.Background(), elog.TraceIdHeader, getTraceId(ctx))
	v, ok := ctx.Get(gin.AuthUserKey)
	if ok {
		c = auth.WrapCtxWithAuthUser(c, v)
	}
	return c
}

func Handle[P any, R any](rawFunc HandlerFunc[P, R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		errors.InvalidParamsError.RaiseIf(ctx.ShouldBind(param))
		responseHandler(ctx, nil, rawFunc(wrapGinContext(ctx), param))
	}
}

// GinHandlerFunc used for custom param binding
type GinHandlerFunc[P any, R any] func(c *gin.Context, ctx context.Context, req *P) R

func GinHandle[P any, R any](rawFunc GinHandlerFunc[P, R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		errors.InvalidParamsError.RaiseIf(ctx.ShouldBind(param))
		responseHandler(ctx, nil, rawFunc(ctx, wrapGinContext(ctx), param))
	}
}

type NoParamHandlerFunc[R any] func(c context.Context) R

func NoParamHandle[R any](rawFunc NoParamHandlerFunc[R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseHandler(ctx, nil, rawFunc(wrapGinContext(ctx)))
	}
}

type GinNoParamHandlerFunc[R any] func(c *gin.Context, ctx context.Context) R

func GinNoParamHandle[R any](rawFunc GinNoParamHandlerFunc[R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseHandler(ctx, nil, rawFunc(ctx, wrapGinContext(ctx)))
	}
}
