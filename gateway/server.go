package gateway

import (
	"exermon/gateway/middleware"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/xixiwang12138/xlog"
)

type Option func(*server)

var (
	WithPort = func(port string) Option {
		return func(s *server) {
			s.port = port
		}
	}
	WithSSL = func(keyPath, certPath string) Option {
		return func(s *server) {
			s.cert = &struct {
				key  string
				cert string
			}{key: keyPath, cert: certPath}
		}
	}

	WithServer = func(group map[string]func(g *gin.RouterGroup)) Option {
		return func(s *server) {
			s.group = group
		}
	}

	WithGloabMiddlewares = func(fs ...gin.HandlerFunc) Option {
		return func(s *server) {
			s.globalMiddlewares = append(s.globalMiddlewares, fs...)
		}
	}

	WithGrafully = func() Option {
		return func(s *server) {
			s.grateful = true
		}
	}

	WithResponseProcessor = func(p ResponseProcessor) Option {
		return func(s *server) {
			//s.responseHandler = p
			responseHandler = p
		}
	}
)

type ResponseProcessor func(ctx *gin.Context, err error, res any)

var (
	DefaultMiddlewares = []gin.HandlerFunc{gin.Recovery(), Grateful()}
)

type server struct {
	g *gin.Engine

	group             map[string]func(*gin.RouterGroup)
	globalMiddlewares []gin.HandlerFunc
	responseHandler   ResponseProcessor

	port string
	cert *struct {
		key  string
		cert string
	}

	grateful bool
}

func (server *server) init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}
	server.g = gin.New()
	server.responseHandler = responseHandler
}

func (server *server) RegisterGlobalMiddleWare(middles ...gin.HandlerFunc) {
	server.g.Use(middles...)
}

func (server *server) RegisterMiddleWare(middle gin.HandlerFunc) {
	server.g.Use(middle)
}

func exist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (server *server) start() {
	ssl := server.cert
	//HTTPS启动
	if ssl != nil && exist(ssl.key) && exist(ssl.cert) {
		server.g.Use(middleware.TlsHandler(server.port))
		xlog.Info("[HTTPS] ===> HTTPS", server.port)
		if err := server.g.RunTLS(server.port, ssl.cert, ssl.key); err != nil {
			panic(err)
		}
	} else {
		xlog.Info("[HTTP] CERT FILE NOT EXIST ===> HTTP: ", &server.port, "\n\n")
		if err := server.g.Run(server.port); err != nil {
			panic(err)
		}
	}
}

func (server *server) Setup() {
	server.init()
	server.RegisterGlobalMiddleWare()
	server.registerGroups()
	go server.start()
	if server.grateful {
		//grateful stop server
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		select {
		case sig := <-c:
			{
				xlog.Infof("Got %s signal. Waiting all requests can be handled...\n", sig)
				time.Sleep(4 * time.Second) //ensure all request can be processed
				server.Close()
				os.Exit(1)
			}
		}
	}
}

func (server *server) registerGroups() {
	if server.group == nil {
		return
	}
	for s, f := range server.group {
		server.registerGroup(s, f)
	}
}

func (server *server) registerGroup(groupPrefix string, registerFunc func(g *gin.RouterGroup)) {
	group := server.g.Group(groupPrefix)
	registerFunc(group)
}

var (
	stop int64 = 0
)

func Grateful() gin.HandlerFunc {
	return func(c *gin.Context) {
		if stopped := atomic.LoadInt64(&stop); stopped == 1 {
			c.JSON(404, gin.H{"msg": "服务维护中，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

func (server *server) Close() {
	atomic.AddInt64(&stop, 1)
	xlog.Info("gateway server stop handle request....")
}
