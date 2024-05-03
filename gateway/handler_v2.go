package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/exermon/errors"
)

type HandlerFunc[P any, R any] func(c context.Context, req P) R

func Handle[P any, R any](rawFunc HandlerFunc[P, R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		errors.InvalidParamError.RaiseIf(ctx.ShouldBind(param), "invalid param")
		responseHandler(ctx, nil, rawFunc(ctx, *param))
	}
}

type GinHandlerFunc[P any, R any] func(c *gin.Context, req P) R

func GinHandle[P any, R any](rawFunc GinHandlerFunc[P, R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := new(P)
		errors.InvalidParamError.RaiseIf(ctx.ShouldBind(param), "invalid param")
		responseHandler(ctx, nil, rawFunc(ctx, *param))
	}
}

type NoParamHandlerFunc[R any] func(c context.Context) R

func NoParamHandle[R any](rawFunc NoParamHandlerFunc[R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseHandler(ctx, nil, rawFunc(ctx))
	}
}

type GinNoParamHandlerFunc[R any] func(c *gin.Context) R

func GinNoParamHandle[R any](rawFunc GinNoParamHandlerFunc[R]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		responseHandler(ctx, nil, rawFunc(ctx))
	}
}
