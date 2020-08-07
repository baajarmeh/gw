package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	apiNameKey = "gw-api-name"
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// Context represents a gw Context object, it's extension from gin.Context.
type Context struct {
	*gin.Context
	RequestID      string
	User           User
	Store          Store
	startTime      time.Time
	logger         Logger
	queries        map[string][]string
	params         map[string]interface{}
	globalDbFilter []interface{}
}

// RouterInfo represents a gw Router Info.
type RouterInfo struct {
	Method            string
	Router            string
	Handler           Handler
	Decorators        []Decorator
	Permissions       []Permission
	handlerActionName string
}

func (r RouterInfo) String() string {
	// https://books.studygolang.com/gobyexample/string-formatting
	// Padding left.
	return fmt.Sprintf("%-8s%s -> %s", r.Method, r.Router, r.handlerActionName)
}

func createRouterInfo(method, router string, handlerActionName string, handler Handler, decorators ...Decorator) RouterInfo {
	routerInfo := RouterInfo{
		Method:     method,
		Router:     router,
		Handler:    handler,
		Decorators: decorators,
	}
	pds := FilterDecorator(func(d Decorator) bool {
		return d.Catalog == permissionDecoratorCatalog
	}, decorators...)
	var perms []Permission
	for _, p := range pds {
		if pd, ok := p.MetaData.([]Permission); ok {
			perms = append(routerInfo.Permissions, pd...)
		}
	}
	routerInfo.Permissions = perms
	if handlerActionName == "" {
		handlerActionName = getHandlerFullName(handler)
	}
	routerInfo.handlerActionName = handlerActionName
	return routerInfo
}

func getHandlerFullName(handler Handler) string {
	var val = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	return fmt.Sprintf("%s(ctx *Context)", val)
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
	locker        sync.Mutex
	prefix        string
	server        *gin.Engine
	router        *gin.RouterGroup
	currentRouter *gin.RouterGroup
	routerInfos   []RouterInfo
}

// RouterGroup represents a gw's Group Router info.
type RouterGroup struct {
	*Router
}

func (router *Router) storeRouterStateWithHandlerName(method string, relativePath,
	handlerName string, handler Handler, decorators ...Decorator) {
	router.locker.Lock()
	defer router.locker.Unlock()
	str := path.Join(router.prefix, relativePath)
	router.routerInfos = append(router.routerInfos, createRouterInfo(method, str, handlerName, handler, decorators...))
}

func (router *Router) storeRouterState(method string, relativePath string, handler Handler, decorators ...Decorator) {
	handlerName := getHandlerFullName(handler)
	// Rest API, not register router info.
	// It's routers save by RegisterRestApis(...)
	if strings.Count(handlerName, "github.com/oceanho/gw.init") > 0 {
		return
	}
	router.storeRouterStateWithHandlerName(method, relativePath, handlerName, handler, decorators...)
}

// GET register a http Get router of handler.
func (router *Router) GET(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("GET", relativePath, handler, decorators...)
	router.currentRouter.GET(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// POST register a http POST router of handler.
func (router *Router) POST(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("POST", relativePath, handler, decorators...)
	router.currentRouter.POST(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// PUT register a http PUT router of handler.
func (router *Router) PUT(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("PUT", relativePath, handler, decorators...)
	router.currentRouter.PUT(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// HEAD register a http HEAD router of handler.
func (router *Router) HEAD(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("HEAD", relativePath, handler, decorators...)
	router.currentRouter.HEAD(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// DELETE register a http DELETE router of handler.
func (router *Router) DELETE(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("DELETE", relativePath, handler, decorators...)
	router.currentRouter.DELETE(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// OPTIONS register a http OPTIONS router of handler.
func (router *Router) OPTIONS(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("OPTIONS", relativePath, handler, decorators...)
	router.currentRouter.OPTIONS(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// PATCH register a http PATCH router of handler.
func (router *Router) PATCH(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("PATCH", relativePath, handler, decorators...)
	router.currentRouter.PATCH(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
	})
}

// Any register a any HTTP method router of handler.
func (router *Router) Any(relativePath string, handler Handler, decorators ...Decorator) {
	a, b := splitDecorators(decorators...)
	router.storeRouterState("Any", relativePath, handler, decorators...)
	router.currentRouter.Any(relativePath, func(c *gin.Context) {
		handle(c, handler, a, b)
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
func (router *Router) Group(relativePath string, handler Handler, decorators ...Decorator) *RouterGroup {
	rg := &RouterGroup{
		router,
	}
	if handler != nil {
		a, b := splitDecorators(decorators...)
		router.storeRouterState("Group", relativePath, handler, decorators...)
		rg.currentRouter = rg.router.Group(relativePath, func(c *gin.Context) {
			handle(c, handler, a, b)
		})
	} else {
		rg.currentRouter = rg.router.Group(relativePath)
		rg.prefix = path.Join(router.prefix, relativePath)
	}
	return rg
}

// Config returns a snapshot of the current Context's conf.Config object.
func (c *Context) Config() conf.Config {
	return config(c.Context)
}

// Bind define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.Bind(...).
func (c *Context) Bind(out interface{}) error {
	if err := c.Context.Bind(out); err != nil {
		c.JSON400Msg(400, fmt.Sprintf("invalid request parameters, details: \n%v", err))
		return err
	}
	return nil
}

// BindQuery define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.BindQuery(...).
func (c *Context) BindQuery(out interface{}) error {
	if err := c.Context.BindQuery(out); err != nil {
		c.JSON400Msg(400, fmt.Sprintf("invalid request parameters, details: \n%v", err))
		return err
	}
	return nil
}

func handle(c *gin.Context, handler Handler, beforeDecorators, afterDecorators []Decorator) {
	s := hostServer(c)
	user := getUser(c)
	requestID := getRequestID(s, c)
	store := getStore(c, s, *user)
	ctx := makeCtx(c, *user, store, requestID)
	// action before Decorators
	var msg string
	var err error
	for _, d := range beforeDecorators {
		msg, err = d.Before(ctx)
		if err != nil {
			break
		}
	}
	if err != nil {
		if msg == "" {
			msg = "caller decorator fail."
		}
		body := respBody(http.StatusBadRequest, requestID, errDefault403Msg, msg)
		c.JSON(http.StatusBadRequest, body)
		return
	}

	// process Action handler.
	handler(ctx)

	// action after Decorators
	l := len(afterDecorators)
	if l < 1 {
		return
	}
	for i := l - 1; i >= 0; i-- {
		msg, err = afterDecorators[i].After(ctx)
		if err != nil {
			break
		}
	}
	if err != nil {
		if msg == "" {
			msg = "caller decorator fail."
		}
		body := respBody(http.StatusBadRequest, requestID, errDefault403Msg, msg)
		c.JSON(http.StatusBadRequest, body)
		return
	}
}

func makeCtx(c *gin.Context, user User, store Store, requestID string) *Context {
	ctx := &Context{
		User:      user,
		RequestID: requestID,
		Store:     store,
		Context:   c,
		startTime: time.Now(),
		logger:    getLogger(c),
		queries:   make(map[string][]string),
		params:    make(map[string]interface{}),
	}
	return ctx
}

func splitDecorators(decorators ...Decorator) (before, after []Decorator) {
	var a, b []Decorator
	for _, d := range decorators {
		if d.After != nil {
			b = append(b, d)
		}
		if d.Before != nil {
			a = append(a, d)
		}
	}
	return b, a
}
