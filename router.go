package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"time"
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// Context represents a gw Context object, it's extension from gin.Context.
type Context struct {
	*gin.Context
	RequestID string
	User      User
	Store     Store
	startTime time.Time
	queries   map[string][]string
	params    map[string]interface{}
}

// Query returns a string from queries.
func (c Context) Query(key string) string {
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
func (c Context) QueryArray(key string) []string {
	return c.queries[key]
}

// StartTime returns the Context start *time.Time
func (c Context) StartTime() *time.Time {
	return &c.startTime
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

// Config returns a snapshot of the current Context's conf.Config object.
func (c *Context) Config() conf.Config {
	return getConfig(c.Context)
}

// Bind represent a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.Bind(...).
func (c *Context) Bind(out interface{}) error {
	if err := c.Context.Bind(out); err != nil {
		c.Err400Msg(4000, fmt.Sprintf("invalid request parameters, details: \n%v", err))
		return err
	}
	return nil
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
		startTime: time.Now(),
		queries:   make(map[string][]string),
		params:    make(map[string]interface{}),
	}
	return ctx
}

func parseParams() map[string]interface{} {
	panic("impl please.")
}