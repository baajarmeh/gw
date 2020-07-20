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
	queries   map[string][]string
	params    map[string]interface{}
}

func (c *ApiContext) Query(key string) string {
	return c.queries[key][0]
}

func (c *ApiContext) QueryArray(key string) []string {
	return c.queries[key]
}

type ApiRouter struct {
	server        *gin.Engine
	router        *gin.RouterGroup
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

func (router *ApiRouter) Use(middleware ...gin.HandlerFunc) {
	router.currentRouter.Use(middleware...)
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

func reflectRouter(relativePath string, handler Handler) {
	// prefix  pattern   suffix.
}

func handle(c *gin.Context, handler Handler) {
	handler(makeApiCtx(c))
}

//
// TODO(Ocean): impl it.
//  sp, reg, str, uint64
// ===========================
//
//func handleBySubPath(c *gin.Context, handler Handler) {
//	handler(makeApiCtx(c))
//}
//
//func handleByUint64(c *gin.Context, handler Handler) {
//	handler(makeApiCtx(c))
//}
//
//func handleByStr(c *gin.Context, handler Handler) {
//	handler(makeApiCtx(c))
//}
//
//func handleByRegex(c *gin.Context, handler Handler) {
//	handler(makeApiCtx(c))
//}

func makeApiCtx(c *gin.Context) *ApiContext {
	user := auth.GetUser(c)
	requestId := req.GetRequestId(c)
	backendStore := store.GetBackend(c, user)
	ctx := &ApiContext{
		User:      user,
		RequestId: requestId,
		Store:     backendStore,
		Context:   c,
		queries:   make(map[string][]string),
		params:    make(map[string]interface{}),
	}
	return ctx
}

func parseParams() map[string]interface{} {
	panic("impl please.")
}
