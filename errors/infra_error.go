package errors

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

var (
	DefaultError      = &Error{Code: 100001, Desc: "系统内部错误", HTTPCode: 200}
	NotFoundError     = &Error{Code: 100002, Desc: "资源不存在", HTTPCode: 404}
	InvalidParamError = &Error{Code: 100003, Desc: "参数错误", HTTPCode: 400}
	ForbiddenError    = &Error{Code: 100004, Desc: "无权限", HTTPCode: 403}
	InfraError        = &Error{Code: 100005, Desc: "基础组件错误", HTTPCode: 500}
	UnauthError       = &Error{Code: 100005, Desc: "用户验证失败", HTTPCode: 401}
)

type Error struct {
	Code     uint64 // 错误码
	Desc     string // 错误简单描述
	HTTPCode int    // HTTP状态码(若有)

	Msg string // 错误信息
	Ori error  // 原始错误
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s, origin: %s", e.Code, e.Msg, e.Ori)
}

func (e Error) RaiseIf(cond any, msg string, args ...any) {
	switch cond.(type) {
	case bool:
		if cond.(bool) {
			err := Error{
				Code:     e.Code,
				Desc:     e.Desc,
				HTTPCode: e.HTTPCode,
				Msg:      fmt.Sprintf(msg, args...),
			}
			panic(err)
		}
	case error:
		if cond != nil {
			err := Error{
				Code:     e.Code,
				Desc:     e.Desc,
				HTTPCode: e.HTTPCode,
				Msg:      fmt.Sprintf(msg, args...),
				Ori:      cond.(error),
			}
			panic(err)
		}
	}
}

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(Error); ok {
					body := gin.H{
						"code": e.Code,
						"msg":  e.Msg,
						"desc": e.Desc,
					}
					if e.Ori != nil {
						body["ori"] = e.Ori.Error()
					}
					ctx.AbortWithStatusJSON(e.HTTPCode, body)
					slog.ErrorContext(ctx, "recovered", slog.Any("error", err))
					return
				}
				slog.ErrorContext(ctx, "recovered", slog.Any("unexpected", err), slog.String("stack", string(debug.Stack())))
				ctx.AbortWithStatusJSON(500, gin.H{"code": 100001, "msg": "系统内部错误", "detail": err})
			}
		}()
		ctx.Next()
	}
}
