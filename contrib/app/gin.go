package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app/auth"
	"github.com/oceanho/gw/contrib/app/store"
	"time"
)

type Handler func(ctx *ApiContext) interface{}

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

func (router *ApiRouter) GET(relativePath string, handlers ...Handler) {
	router.router.GET(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) POST(relativePath string, handlers ...Handler) {
	router.router.POST(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) PUT(relativePath string, handlers ...Handler) {
	router.router.PUT(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) HEAD(relativePath string, handlers ...Handler) {
	router.router.HEAD(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) DELETE(relativePath string, handlers ...Handler) {
	router.router.DELETE(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) OPTIONS(relativePath string, handlers ...Handler) {
	router.router.OPTIONS(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) PATCH(relativePath string, handlers ...Handler) {
	router.router.PATCH(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) Any(relativePath string, handlers ...Handler) {
	router.router.Any(relativePath, func(ctx *gin.Context) {
		handle(ctx, handlers...)
	})
}

func (router *ApiRouter) Handlers() gin.HandlersChain {
	return router.router.Handlers
}

func (router *ApiRouter) Group(relativePath string, handlers ...Handler) *ApiRouteGroup {
	apiRg := &ApiRouteGroup{
		router,
	}
	if len(handlers) > 0 {
		apiRg.server.Group(relativePath, func(ctx *gin.Context) {
			handle(ctx, handlers...)
		})
	} else {
		apiRg.server.Group(relativePath)
	}
	return apiRg
}

func handle(ctx *gin.Context, handlers ...Handler) {
	// FIXME(Ocean): Should be removal the following code.
	// fmt.Printf("handlers: %v\n", handlers)
	results := make([]interface{}, 0)
	for _, handler := range handlers {
		results = append(results, handler(makeApiCtx(ctx)))
	}
	if len(results) == 1 {
		ctx.JSON(200, results[0])
		return
	}
	ctx.JSON(200, results)
}

func makeApiCtx(ctx *gin.Context) *ApiContext {
	_ctx := &ApiContext{
		User:    nil,
		Store:   nil,
		RequestId: genRequestID(),
		Context: ctx,
	}
	return _ctx
}

func genRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
