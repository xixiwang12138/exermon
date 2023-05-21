package _examples

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xixiwang12138/exermon/gateway"
)

func main() {
	group := map[string]func(group *gin.RouterGroup){
		"/api/user": func(group *gin.RouterGroup) {
			group.GET("/info", gateway.Handler(GetUserInfo))
		},
	}
	s := gateway.NewServer(
		gateway.WithServer(group),
		gateway.WithGracefully(),
		gateway.WithPort(":8082"))
	s.Start()
}

func GetUserInfo(c *gin.Context, xl context.Context, p *string) (*any, error) {
	return nil, nil
}
