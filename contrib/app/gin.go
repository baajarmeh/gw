package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app/auth"
	"github.com/oceanho/gw/contrib/app/store"
	"time"
)

type Handler func(ctx *ApiContext)
type HandlerRet func(ctx *ApiContext) interface{}

type ApiContext struct {
	*gin.Context
	RequestId string
	User      *auth.User
	Store     *store.Backend
}

type ApiRouter struct {
	server *gin.Engine
	router *gin.RouterGroup
}

type ApiRouteGroup struct {
	*ApiRouter
}

func (router *ApiRouter) GET(relativePath string, handler Handler) {
	router.router.GET(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) POST(relativePath string, handler Handler) {
	router.router.POST(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) PUT(relativePath string, handler Handler) {
	router.router.PUT(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}
func (router *ApiRouter) HEAD(relativePath string, handler Handler) {
	router.router.HEAD(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) DELETE(relativePath string, handler Handler) {
	router.router.DELETE(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) OPTIONS(relativePath string, handler Handler) {
	router.router.OPTIONS(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) PATCH(relativePath string, handler Handler) {
	router.router.PATCH(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) Any(relativePath string, handler Handler) {
	router.router.Any(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) Handlers() gin.HandlersChain {
	return router.router.Handlers
}

func (router *ApiRouter) Group(relativePath string, handler Handler) *ApiRouteGroup {
	apiRg := &ApiRouteGroup{
		router,
	}
	if handler != nil {
		apiRg.server.Group(relativePath, func(c *gin.Context) {
			handle(c, handler)
		})
	} else {
		apiRg.server.Group(relativePath)
	}
	return apiRg
}

func handle(c *gin.Context, handler Handler) {
	handler(makeApiCtx(c))
}

func makeApiCtx(c *gin.Context) *ApiContext {
	ctx := &ApiContext{
		User:      nil,
		Store:     nil,
		RequestId: genRequestID(),
		Context:   c,
	}
	return ctx
}

func genRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
