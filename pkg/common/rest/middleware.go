package rest

import "github.com/gin-gonic/gin"

type Middleware interface {
	Name() string
	Middleware() gin.HandlerFunc
}

var (
	// To assume implementations of Middleware
	_ Middleware = (*MiddlewareCors)(nil)
	_ Middleware = (*MiddlewareAuth)(nil)
)

type MiddlewareCors struct {
	CorsHeader string
}

func (c MiddlewareCors) Name() string {
	return "CORS"
}

func (c MiddlewareCors) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		corsHeaders := c.CorsHeader
		if corsHeaders == "" {
			corsHeaders = "*"
		}

		// log.Debugf("Received a new request from %s, CORS header is enabled and set to: %s", ctx.GetHeader("Origin"), corsHeaders)
		ctx.Header("Access-Control-Allow-Origin", corsHeaders)
		ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Registry-Auth")
		ctx.Header("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")
		ctx.Next()
	}
}

type MiddlewareAuth struct {
	AuthFunc func(ctx *gin.Context)
}

func (a MiddlewareAuth) Name() string {
	return "Auth"
}

func (a MiddlewareAuth) Middleware() gin.HandlerFunc {
	return a.AuthFunc
}
