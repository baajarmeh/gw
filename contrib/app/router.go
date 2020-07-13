package app

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app/auth"
	"github.com/oceanho/gw/contrib/app/req"
	"github.com/oceanho/gw/contrib/app/store"
)

type Handler func(ctx *ApiContext)

type ApiContext struct {
	*gin.Context
	RequestId string
	User      auth.User
	Store     store.Backend
}

type ApiRouter struct {
	server *gin.Engine
	router *gin.RouterGroup
	currentRouter *gin.RouterGroup
}

type ApiRouteGroup struct {
	*ApiRouter
}

func (router *ApiRouter) GET(relativePath string, handler Handler) {
	router.currentRouter.GET(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) POST(relativePath string, handler Handler) {
	router.currentRouter.POST(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) PUT(relativePath string, handler Handler) {
	router.currentRouter.PUT(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}
func (router *ApiRouter) HEAD(relativePath string, handler Handler) {
	router.currentRouter.HEAD(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) DELETE(relativePath string, handler Handler) {
	router.currentRouter.DELETE(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) OPTIONS(relativePath string, handler Handler) {
	router.currentRouter.OPTIONS(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) PATCH(relativePath string, handler Handler) {
	router.currentRouter.PATCH(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) Any(relativePath string, handler Handler) {
	router.currentRouter.Any(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

func (router *ApiRouter) Handlers() gin.HandlersChain {
	return router.currentRouter.Handlers
}

func (router *ApiRouter) Group(relativePath string, handler Handler) *ApiRouteGroup {
	rg := &ApiRouteGroup{
		router,
	}
	if handler != nil {
		rg.currentRouter = rg.router.Group(relativePath, func(c *gin.Context) {
			handle(c, handler)
		})
	} else {
		rg.currentRouter = rg.router.Group(relativePath)
	}
	return rg
}

func handle(c *gin.Context, handler Handler) {
	handler(makeApiCtx(c))
}

func makeApiCtx(c *gin.Context) *ApiContext {
	ctx := &ApiContext{
		User:      auth.GetUser(c),
		Store:     store.GetBackend(c),
		RequestId: req.GetRequestId(c),
		Context:   c,
	}
	return ctx
}
