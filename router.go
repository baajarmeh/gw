package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"sync"
)

const (
	gwRouterInfoKey = "gw-router"
)

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

// RouterInfo represents a gw Router Info.
type RouterInfo struct {
	Method            string
	UrlPath           string
	Handler           Handler
	Decorators        []Decorator
	Permissions       []Permission
	handlerActionName string
	beforeDecorators  []Decorator
	afterDecorators   []Decorator
}

func (r *RouterInfo) String() string {
	// https://books.studygolang.com/gobyexample/string-formatting
	// Padding left.
	before, after := splitDecorators(r.Decorators...)
	return fmt.Sprintf("%-8s%s -> %s "+
		"[decorators(%d before, %d after)]", r.Method, r.UrlPath, r.handlerActionName, len(before), len(after))
}

func (router *Router) createRouter(method, relativePath string, handler Handler, handlerActionName string, decorators ...Decorator) {
	urlPath := path.Join(router.currentRouter.BasePath(), relativePath)
	//var _decorators []Decorator
	//_decorators = append(_decorators, NewStoreDbSetupDecorator(func(ctx Context, db *gorm.DB) *gorm.DB {
	//	return db.Set(gwDbUserInfoKey, ctx.User())
	//}))
	//_decorators = append(_decorators, decorators...)
	routerInfo := RouterInfo{
		Method:     method,
		UrlPath:    urlPath,
		Handler:    handler,
		Decorators: decorators,
	}
	beforeDecorators, afterDecorators := splitDecorators(decorators...)
	pds := FilterDecorator(func(d Decorator) bool {
		return d.Catalog == permissionDecoratorCatalog
	}, decorators...)
	var perms []Permission
	for _, p := range pds {
		if pd, ok := p.MetaData.([]Permission); ok {
			perms = append(routerInfo.Permissions, pd...)
		}
	}
	if handlerActionName == "" {
		handlerActionName = getHandlerFullName(handler)
	}
	routerInfo.Permissions = perms
	routerInfo.afterDecorators = afterDecorators
	routerInfo.beforeDecorators = beforeDecorators
	routerInfo.handlerActionName = handlerActionName
	if method == "any" {
		router.currentRouter.Any(relativePath, func(c *gin.Context) {
			c.Set(gwRouterInfoKey, routerInfo)
			handle(c)
		})
		return
	}
	if method == "group" {
		router.currentRouter.Group(relativePath, func(c *gin.Context) {
			c.Set(gwRouterInfoKey, routerInfo)
			handle(c)
		})
		return
	}
	// gin router
	router.currentRouter.Handle(method, relativePath, func(c *gin.Context) {
		c.Set(gwRouterInfoKey, routerInfo)
		handle(c)
	})
	router.locker.Lock()
	defer router.locker.Unlock()
	router.routerInfos = append(router.routerInfos, routerInfo)
}

func getHandlerFullName(handler Handler) string {
	var val = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	return fmt.Sprintf("%s(ctx *Context)", val)
}

// GET register a http Get router of handler.
func (router *Router) GET(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodGet, relativePath, handler, "", decorators...)
}

// POST register a http POST router of handler.
func (router *Router) POST(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPost, relativePath, handler, "", decorators...)
}

// PUT register a http PUT router of handler.
func (router *Router) PUT(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPut, relativePath, handler, "", decorators...)
}

// HEAD register a http HEAD router of handler.
func (router *Router) HEAD(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodHead, relativePath, handler, "", decorators...)
}

// DELETE register a http DELETE router of handler.
func (router *Router) DELETE(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodDelete, relativePath, handler, "", decorators...)
}

// OPTIONS register a http OPTIONS router of handler.
func (router *Router) OPTIONS(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodOptions, relativePath, handler, "", decorators...)
}

// PATCH register a http PATCH router of handler.
func (router *Router) PATCH(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodPatch, relativePath, handler, "", decorators...)
}

// Any register a all HTTP methods router of handler.
func (router *Router) Any(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter("any", relativePath, handler, "", decorators...)
}

// Any register a HTTP Connect router of handler.
func (router *Router) Connect(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodConnect, relativePath, handler, "", decorators...)
}

// Trace register a HTTP Trace router of handler.
func (router *Router) Trace(relativePath string, handler Handler, decorators ...Decorator) {
	router.createRouter(http.MethodTrace, relativePath, handler, "", decorators...)
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
	if handler == nil {
		rg.currentRouter = rg.router.Group(relativePath)
	} else {
		router.createRouter("group", relativePath, handler, "", decorators...)
	}
	rg.prefix = path.Join(router.prefix, relativePath)
	return rg
}
