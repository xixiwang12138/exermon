package gateway

// 用于对Gin返回值的封装，包括全局异常处理
import (
	"github.com/gin-gonic/gin"
)

type Resp struct{}

// Success 返回成功
//func Success(c *gin.Context, data interface{}) {
//	var resp struct {
//		Code int32       `json:"code"`
//		Data interface{} `json:"data"`
//	}
//	resp.Data = data
//	resp.Code = 0
//	c.JSON(gateway.StatusOK, resp)
//	return
//}

// Err 返回失败
func Err(c *gin.Context, httpCode int, code int, msg string) {
	c.JSON(httpCode, map[string]interface{}{
		"code":  code,
		"error": msg, //错误提示
	})
	return
}

//// CatchError 异常捕获中间件
//func CatchError() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			if definederr := recover(); definederr != nil {
//				url := c.Request.URL
//				method := c.Request.Method
//				log.Printf("| url [%s] | method | [%s] | error [%s] |", url, method, definederr)
//				log.Printf(cts.ErrorFormat, definederr)
//				//var exception definederr.Exception
//				////判断是否为自定义异常
//				//switch definederr.(type) {
//				//case string: //自定义异常
//				//	definederr := json.Unmarshal([]byte(string(definederr.(string))), &exception)
//				//	if definederr != nil {
//				//		Err(c, gateway.StatusBadRequest, exception.Code, "未知错误，请联系管理员！")
//				//		c.Abort()
//				//		return
//				//	}
//				//	Err(c, exception.HttpCode, exception.Code, exception.Error)
//				//default:
//				Err(c, 500, 1, definederr.(runtime.Error).Error())
//				//}
//				c.Abort()
//			}
//		}()
//		c.Next()
//	}
//}

// CommonResponse 统一返回结构
type CommonResponse struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
