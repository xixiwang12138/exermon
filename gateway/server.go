package gateway

import (
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"

	"github.com/xixiwang12138/exermon/elog"
	"github.com/xixiwang12138/exermon/errors"
	"github.com/xixiwang12138/exermon/gateway/middleware"
	"github.com/xixiwang12138/exermon/log"

	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type ResponseProcessor func(ctx *gin.Context, err error, res any)

var (
	DefaultMiddlewares = []gin.HandlerFunc{middleware.Cors(), log.TracingLogger(), errors.Recover()}
)

type Option func(*Server)

var (
	WithPort = func(port string) Option {
		return func(s *Server) {
			s.port = port
		}
	}
	WithSSL = func(keyPath, certPath string) Option {
		return func(s *Server) {
			s.cert = &struct {
				key  string
				cert string
			}{key: keyPath, cert: certPath}
		}
	}

	WithServer = func(group map[string]func(g *gin.RouterGroup)) Option {
		return func(s *Server) {
			s.group = group
		}
	}

	WithGlobalMiddlewares = func(fs ...gin.HandlerFunc) Option {
		return func(s *Server) {
			s.globalMiddlewares = append(s.globalMiddlewares, fs...)
		}
	}

	WithGracefully = func() Option {
		return func(s *Server) {
			s.globalMiddlewares = append(s.globalMiddlewares, Grateful)
			s.grateful = true
		}
	}

	WithAutoTLS = func(domain ...string) Option {
		return func(s *Server) {
			s.autoTls = true
			s.domains = domain
		}
	}

	WithResponseProcessor = func(p ResponseProcessor) Option {
		return func(s *Server) {
			//s.responseHandler = p
			responseHandler = p
		}
	}
)

type Server struct {
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
	autoTls  bool
	domains  []string
}

func NewServer(options ...Option) *Server {
	s := &Server{}
	s.g = gin.New()

	for _, option := range options {
		option(s)
	}
	return s
}

func (server *Server) init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		el := elog.WithTraceId("start")
		el.Info("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	server.g = gin.New()
	server.responseHandler = responseHandler
	server.globalMiddlewares = append(DefaultMiddlewares, server.globalMiddlewares...)
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

func (server *Server) start() {
	ssl := server.cert
	cl := elog.WithTraceId("start")
	//HTTPS启动
	if server.autoTls {
		if err := autotls.Run(server.g, server.domains...); err != nil {
			panic(err)
		}
	}
	if ssl != nil && exist(ssl.key) && exist(ssl.cert) {
		server.g.Use(middleware.TlsHandler(server.port))
		cl.Info("[HTTPS] ===> HTTPS： %s", server.port)
		if err := server.g.RunTLS(server.port, ssl.cert, ssl.key); err != nil {
			panic(err)
		}
	} else {
		cl.Info("[HTTP] CERT FILE NOT EXIST ===> HTTP: %s", server.port)
		if err := server.g.Run(server.port); err != nil {
			panic(err)
		}
	}
}

func (server *Server) Start() {
	server.init()
	server.g.Use(server.globalMiddlewares...)
	server.registerGroups()
	go server.start()
	if server.grateful {
		//grateful stop Server
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		select {
		case sig := <-c:
			{
				cl := elog.WithTraceId("terminate")
				cl.Info("Got %s signal. Waiting all requests can be handled...\n", sig)
				time.Sleep(4 * time.Second) //ensure all request can be processed
				server.Close()
				os.Exit(1)
			}
		}
	}
	select {}
}

func (server *Server) registerGroups() {
	if server.group == nil {
		return
	}
	for s, f := range server.group {
		server.registerGroup(s, f)
	}
}

func (server *Server) registerGroup(groupPrefix string, registerFunc func(g *gin.RouterGroup)) {
	group := server.g.Group(groupPrefix)
	registerFunc(group)
}

var (
	stop int64 = 0
)

func Grateful(c *gin.Context) {
	if stopped := atomic.LoadInt64(&stop); stopped == 1 {
		c.JSON(404, gin.H{"msg": "服务维护中，请稍后再试"})
		c.Abort()
		return
	}
	c.Next()
	return
}

func (server *Server) Close() {
	atomic.AddInt64(&stop, 1)
	el := elog.WithTraceId("close")
	el.Info("gateway Server stop handle request....")
}
