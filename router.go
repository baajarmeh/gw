package gw

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// Context represents a gw Context object, it's extension from gin.Context.
type Context struct {
	*gin.Context
	RequestID string
	User      User
	Store     Store
	queries   map[string][]string
	params    map[string]interface{}
}

// Query returns a string from queries.
func (c *Context) Query(key string) string {
	val := c.Context.Query(key)
	if val == "" {
		queries := c.QueryArray(key)
		if len(queries) > 0 {
			val = queries[0]
		}
	}
	return val
}

// QueryArray returns a string array from queries.
func (c *Context) QueryArray(key string) []string {
	return c.queries[key]
}

// Router represents a gw's Router info.
type Router struct {
	server        *gin.Engine
	router        *gin.RouterGroup
	currentRouter *gin.RouterGroup
}

// RouteGroup represents a gw's Group Router info.
type RouteGroup struct {
	*Router
}

// GET register a http Get router of handler.
func (router *Router) GET(relativePath string, handler Handler) {
	router.currentRouter.GET(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// POST register a http POST router of handler.
func (router *Router) POST(relativePath string, handler Handler) {
	router.currentRouter.POST(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// PUT register a http PUT router of handler.
func (router *Router) PUT(relativePath string, handler Handler) {
	router.currentRouter.PUT(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// HEAD register a http HEAD router of handler.
func (router *Router) HEAD(relativePath string, handler Handler) {
	router.currentRouter.HEAD(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// DELETE register a http DELETE router of handler.
func (router *Router) DELETE(relativePath string, handler Handler) {
	router.currentRouter.DELETE(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// OPTIONS register a http OPTIONS router of handler.
func (router *Router) OPTIONS(relativePath string, handler Handler) {
	router.currentRouter.OPTIONS(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// PATCH register a http PATCH router of handler.
func (router *Router) PATCH(relativePath string, handler Handler) {
	router.currentRouter.PATCH(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// Any register a any HTTP method router of handler.
func (router *Router) Any(relativePath string, handler Handler) {
	router.currentRouter.Any(relativePath, func(c *gin.Context) {
		handle(c, handler)
	})
}

// Use register a gin Middleware of handler.
func (router *Router) Use(middleware ...gin.HandlerFunc) {
	router.currentRouter.Use(middleware...)
}

// Handlers returns the current router of gin.HandlersChain
func (router *Router) Handlers() gin.HandlersChain {
	return router.currentRouter.Handlers
}

// Group returns a new route group.
func (router *Router) Group(relativePath string, handler Handler) *RouteGroup {
	rg := &RouteGroup{
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

// OK response a JSON formatter to client with http status = 200.
func (c *Context) OK(payload interface{}) {
	c.JSON(http.StatusOK, 0, payload)
}

// Fault403 response a JSON formatter to client with http status = 403.
func (c *Context) Fault403(status int, errMsg, payload interface{}) {
	result := gin.H{
		"status":  status,
		"msg":     errMsg,
		"payload": payload,
	}
	c.Context.JSON(403, result)
}

// Fault404 response a JSON formatter to client with http status = 404.
func (c *Context) Fault404(status int, outs ...interface{}) {
	c.Context.JSON(http.StatusNotFound, resp(status, outs...))
}

// Fault500 response a JSON formatter to client with http status = 500.
func (c *Context) Fault500(status int, outs ...interface{}) {
	c.Context.JSON(http.StatusInternalServerError, resp(status, outs...))
}

// JSON response a JSON formatter to client.
func (c *Context) JSON(code int, outs ...interface{}) {
	c.JSONStatus(code, 0, outs)
}

// JSON response a JSON formatter to client.
func (c *Context) JSONStatus(code int, status int, outs ...interface{}) {
	c.Context.JSON(code, resp(status, outs...))
}

func resp(status int, outs ...interface{}) interface{} {
	var errMsg interface{}
	var payload interface{}
	if len(outs) == 2 {
		errMsg = outs[0]
		payload = outs[1]
	} else if len(outs) == 1 {
		errMsg = nil
		payload = outs[0]
	}
	return gin.H{
		"status":  status,
		"err":     errMsg,
		"payload": payload,
	}
}

func reflectRouter(relativePath string, handler Handler) {
	// prefix  pattern   suffix.
}

func handle(c *gin.Context, handler Handler) {
	handler(makeCtx(c))
}

//
// TODO(Ocean): impl it.
//  sp, reg, str, uint64
// ===========================
//
// func handleBySubPath(c *gin.Context, handler Handler) {
// 	handler(makeCtx(c))
// }
//
// func handleByUint64(c *gin.Context, handler Handler) {
// 	handler(makeCtx(c))
// }
//
// func handleByStr(c *gin.Context, handler Handler) {
// 	handler(makeCtx(c))
// }
//
// func handleByRegex(c *gin.Context, handler Handler) {
// 	handler(makeCtx(c))
// }

func makeCtx(c *gin.Context) *Context {
	user := getUser(c)
	requestID := getRequestID(c)
	backendStore := getStore(c, user)
	ctx := &Context{
		User:      user,
		RequestID: requestID,
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
