package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/oceanho/gw/contrib/app/auth"
	"github.com/oceanho/gw/contrib/app/resp"
	"github.com/oceanho/gw/contrib/app/store"
	"time"
)

type Handler func(ctx *ApiContext) resp.IApiResp

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
	router.router.GET(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) POST(relativePath string, handler Handler) {
	router.router.POST(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) PUT(relativePath string, handler Handler) {
	router.router.PUT(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) HEAD(relativePath string, handler Handler) {
	router.router.HEAD(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) DELETE(relativePath string, handler Handler) {
	router.router.DELETE(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) OPTIONS(relativePath string, handler Handler) {
	router.router.OPTIONS(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) PATCH(relativePath string, handler Handler) {
	router.router.PATCH(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
	})
}

func (router *ApiRouter) Any(relativePath string, handler Handler) {
	router.router.Any(relativePath, func(ctx *gin.Context) {
		handle(ctx, handler)
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
		apiRg.server.Group(relativePath, func(ctx *gin.Context) {
			handle(ctx, handler)
		})
	} else {
		apiRg.server.Group(relativePath)
	}
	return apiRg
}

func handle(ctx *gin.Context, handler Handler) {
	respObj, ok := handler(makeApiCtx(ctx)).(resp.IApiResp)
	if ok {
		renderObj := render.Reader{
			Headers:       respObj.GetHeaders(),
			ContentType:   respObj.GetContentType(),
			Reader:        respObj.GetBodyReader(),
			ContentLength: respObj.GetContentLength(),
		}
		ctx.Render(respObj.GetCode(), renderObj)
	}
}

func makeApiCtx(ctx *gin.Context) *ApiContext {
	_ctx := &ApiContext{
		User:      nil,
		Store:     nil,
		RequestId: genRequestID(),
		Context:   ctx,
	}
	return _ctx
}

func genRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
