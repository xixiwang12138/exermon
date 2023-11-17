package safe

import (
	"context"
	"github.com/xixiwang12138/exermon/elog"
	"runtime"
)

func Go(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if obj := recover(); obj != nil {
				cl := elog.WithContext(ctx)
				cl.Error("[panic] object: %+v", obj)
				// 获取panic的错误信息和堆栈信息
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				errMsg := string(stack[:length])
				// 打印红色错误信息
				cl.Error(errMsg)
			}
		}()
		f()
	}()
}

func GoWithAfter(ctx context.Context, f, afterPanic func()) {
	go func() {
		defer func() {
			if obj := recover(); obj != nil {
				cl := elog.WithContext(ctx)
				cl.Error("[panic] object: %+v", obj)
				// 获取panic的错误信息和堆栈信息
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				errMsg := string(stack[:length])
				// 打印红色错误信息
				cl.Error(errMsg)

				afterPanic()
			}
		}()
		f()
	}()
}
