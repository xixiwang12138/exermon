package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xixiwang12138/exermon/elog"
	"runtime/debug"
)

var (
	DefaultError         = &Error{Code: 100001, HTTPCode: 200}
	NotFoundError        = &Error{Code: 100002, HTTPCode: 404}
	InvalidParamsError   = &Error{Code: 100003, HTTPCode: 400}
	ForbiddenError       = &Error{Code: 100004, HTTPCode: 403}
	InfraError           = &Error{Code: 100005, HTTPCode: 500}
	UnauthError          = &Error{Code: 100006, HTTPCode: 401}
	ServerBusyError      = &Error{Code: 100007, HTTPCode: 503}
	TooManyRequestsError = &Error{Code: 100008, HTTPCode: 429}
)

type Error struct {
	Code     uint64 // 错误码
	Desc     string // 错误简单描述
	HTTPCode int    // HTTP状态码(若有)

	Msg string // 错误信息
	Ori error  // 原始错误
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s, origin: %s", e.Code, e.Msg, e.Ori)
}

func (e *Error) RaiseWithFormatIf(cond any, msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	e.RaiseIf(cond, m)
}

func (e *Error) RaiseIf(cond any, msg ...string) {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	switch cond.(type) {
	case bool:
		if cond.(bool) {
			err := &Error{
				Code:     e.Code,
				Desc:     e.Desc,
				HTTPCode: e.HTTPCode,
				Msg:      m,
			}
			panic(err)
		}
	case error:
		if cond != nil {
			err := &Error{
				Code:     e.Code,
				Desc:     e.Desc,
				HTTPCode: e.HTTPCode,
				Msg:      m,
				Ori:      cond.(error),
			}
			panic(err)
		}
	}
}

func (e *Error) RaiseIfIsError(err error, targetErr error, msg ...string) {
	if errors.Is(err, targetErr) {
		e.RaiseIf(err, msg...)
	}
}

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			el := elog.WithContext(ctx)
			if err := recover(); err != nil {
				if e, ok := err.(*Error); ok {
					body := gin.H{
						"code": e.Code,
						"msg":  e.Msg,
						"desc": e.Desc,
					}
					if e.Ori != nil {
						body["ori"] = e.Ori.Error()
					}
					ctx.AbortWithStatusJSON(e.HTTPCode, body)
					//el.Error(ctx, "recovered", slog.Any("error", err))
					return
				}
				el.Error("unexpected error: %s; stack: %s", err, string(debug.Stack()))
				//el.Error("recovered", slog.Any("unexpected", err), slog.String("stack", string(debug.Stack())))
				ctx.AbortWithStatusJSON(500, gin.H{"code": -1, "desc": "Unknown error", "detail": err})
			}
		}()
		ctx.Next()
	}
}
