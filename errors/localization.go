package errors

import "os"

// langMap lang -> code -> message
var langMap = map[string]map[uint64]string{
	"zh_CN": {
		100001: "系统内部错误",
		100002: "资源不存在",
		100003: "参数校验失败",
		100004: "禁止访问",
		100005: "基础组件错误",
		100006: "用户鉴权失败",
		100007: "服务器繁忙",
		100008: "请求次数过多",
	},
	"en_US": {
		100001: "Internal Server Error",
		100002: "Resource Not Found",
		100003: "Invalid Parameters",
		100004: "Forbidden",
		100005: "Infra Error",
		100006: "Auth Failed",
		100007: "Server Busy",
		100008: "Too Many Requests",
	},
}

func init() {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "zh_CN"
	}
	for _, e := range []*Error{
		DefaultError, NotFoundError, InvalidParamsError,
		ForbiddenError, InfraError, UnauthError,
		ServerBusyError, TooManyRequestsError,
	} {
		e.Desc = langMap[lang][e.Code]
	}
}
