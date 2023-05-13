package _examples

import (
	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/exermon/gateway"
	"github.com/xixiwang12138/exermon/gateway/middleware"
	"github.com/xixiwang12138/xlog"
)

func main() {
	group := map[string]func(group *gin.RouterGroup){
		"/api/user": func(group *gin.RouterGroup) {
			group.GET("/info", gateway.Handler(GetUserInfo))
		},
	}
	s := gateway.NewServer(
		gateway.WithServer(group),
		gateway.WithGrafully(),
		gateway.WithPort(":8082"),
		gateway.WithGloabMiddlewares(middleware.Cors()))
	s.Start()
}

func GetUserInfo(c *gin.Context, xl *xlog.XLogger, p *string) (any, error) {
	return nil, nil
}
