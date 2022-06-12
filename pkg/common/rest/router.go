package rest

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Router 定义了一个路由组，服务需要自己实现一个 Router 来进行组册
type Router interface {
	SingleRoutes() []Route
	GroupRoutes() []Group
}

type Group interface {
	GroupPath() string
	Middlewares() []gin.HandlerFunc
	Routes() []Route
}

type Route interface {
	Handlers() []gin.HandlerFunc
	Method() string
	Path() string
}

type buildInRoute struct {
	handlers []gin.HandlerFunc
	method   string
	path     string
}

// Handlers 返回路由的 HandleFunc 列表的浅拷贝，不要尝试修改
func (d *buildInRoute) Handlers() []gin.HandlerFunc {
	return slices.Clone(d.handlers)
}

func (d *buildInRoute) Method() string {
	return d.method
}

func (d *buildInRoute) Path() string {
	return d.path
}

type defaultGroup struct {
	groupPath   string
	middlewares []gin.HandlerFunc
	routes      []Route
}

func (d *defaultGroup) GroupPath() string {
	return d.groupPath
}

// Middlewares 同 Handlers
func (d *defaultGroup) Middlewares() []gin.HandlerFunc {
	if d.middlewares == nil {
		return []gin.HandlerFunc{}
	}
	return slices.Clone(d.middlewares)
}

func (d *defaultGroup) Routes() []Route {
	if d.routes == nil {
		return []Route{}
	}
	return slices.Clone(d.routes)
}

func NewRoute(path string, method string, handler ...gin.HandlerFunc) Route {
	return &buildInRoute{
		handlers: handler,
		method:   method,
		path:     path,
	}
}

func NewGroup(groupPath string, middleware []gin.HandlerFunc, route ...Route) Group {
	return &defaultGroup{
		groupPath:   groupPath,
		middlewares: middleware,
		routes:      route,
	}
}
