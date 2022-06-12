package rest

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
)

type Server struct {
	middlewares []Middleware

	Eng    *gin.Engine
	Router Router
	Name   string
	Addr   string
}

var buildInServer *Server

func (s *Server) init() {
	s.check()

	for _, m := range s.middlewares {
		s.Eng.Use(m.Middleware())
	}
	for _, route := range s.Router.SingleRoutes() {
		s.Eng.Handle(route.Method(), route.Path(), route.Handlers()...)
	}
	for _, group := range s.Router.GroupRoutes() {
		g := s.Eng.Group(group.GroupPath(), group.Middlewares()...)
		for _, route := range group.Routes() {
			g.Handle(route.Method(), route.Path(), route.Handlers()...)
		}
	}
}

func (s *Server) check() {
	if s.middlewares == nil {
		s.middlewares = make([]Middleware, 0)
	}
	if s.Router == nil {
		log.Fatalln("server's Router cannot be nil")
	}
}

func (s *Server) AddMiddleware(m ...Middleware) {
	if s.middlewares == nil {
		s.middlewares = make([]Middleware, 0)
	}
	s.middlewares = append(s.middlewares, m...)
}

func (s *Server) Middlewares() []Middleware {
	return slices.Clone(s.middlewares)
}

func (s *Server) Serve() error {
	s.init()
	return s.Eng.Run(s.Addr)
}

func (s *Server) ServeTLS(cert, key string) error {
	s.init()
	return s.Eng.RunTLS(s.Addr, cert, key)
}

func init() {

	// gin.LoggerWithWriter(nil)

	buildInServer = &Server{
		Eng:         gin.Default(), // TODO 加一个 Logger
		middlewares: make([]Middleware, 0),
	}
}

// SetAddr 设置监听地址
func SetAddr(addr string) {
	buildInServer.Addr = addr
}

// AddMiddleware 添加全局中间件
func AddMiddleware(m ...Middleware) {
	buildInServer.AddMiddleware(m...)
}

// SetRouter 将 router 和 backend 设置进内建的 rest server 内
func SetRouter(router Router) {
	buildInServer.Router = router
}

func SetName(name string) {
	buildInServer.Name = name
}

// Serve 将会阻塞 goroutine，当服务发生错误关闭时会返回 err
func Serve() error {
	return buildInServer.Serve()
}

func ServerTLS(cert, key string) error {
	return buildInServer.ServeTLS(cert, key)
}
