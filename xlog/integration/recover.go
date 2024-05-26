package integration

import (
	"github.com/gin-gonic/gin"

	"github.com/xixiwang12138/exermon/errors"
	"github.com/xixiwang12138/exermon/xlog"
)

var ErrCatchXLoggerIntegration = TracingLogger()

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				xl := xlog.WithContext(ctx)
				xl.LogStack(xlog.WARN)
				if e, ok := err.(errors.Error); ok {
					body := gin.H{
						"code": e.Code,
						"msg":  e.Msg,
						"desc": e.Desc,
					}
					if e.Ori != nil {
						body["ori"] = e.Ori.Error()
					}
					ctx.AbortWithStatusJSON(e.HTTPCode, body)
					xl.Error("[Recover], err: %s", err)
					return
				}
				ctx.AbortWithStatusJSON(500, gin.H{"code": 100001, "msg": "系统内部错误", "detail": err})
			}
		}()
		ctx.Next()
	}
}
