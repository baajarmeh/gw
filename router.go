package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	apiNameKey = "gw-api-name"
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// dynamicHandler
type dynamicHandler func(relativePath string, r *RouterGroup, caller dynamicCaller)

type dynamicCaller struct {
	argsLength        int
	ctrl              IRestAPI
	handler           reflect.Value
	methodName        string
	argsOrderlyBinder []dynamicArgsBinder
}

type dynamicArgsBinder struct {
	dataType reflect.Type
	bindFunc func(p reflect.Type, c *Context) reflect.Value
}

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
	Decorators        []IDecorator
	permissions       []Permission
	handlerActionName string
}

func createRouterInfo(method, router string, handler Handler, decorators ...IDecorator) RouterInfo {
	val := RouterInfo{
		Method:     method,
		Router:     router,
		Handler:    handler,
		Decorators: decorators,
	}
	pds := filterDecorator(func(d IDecorator) bool {
		return d.Catalog() == permissionDecoratorCatalog
	}, decorators...)
	var perms []Permission
	for _, p := range pds {
		if pd, ok := p.(PermissionDecoratorImpl); ok {
			perms = append(val.permissions, pd.perms...)
		}
	}
	val.permissions = perms
	val.handlerActionName = getHandlerFullName(handler)
	return val
}

func getHandlerFullName(handler Handler) string {
	return ""
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
	routerInfos   []RouterInfo
}

// RouterGroup represents a gw's Group Router info.
type RouterGroup struct {
	*Router
}

// GET register a http Get router of handler.
func (router *Router) GET(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("GET", relativePath, handler, decorators...))
	router.currentRouter.GET(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// POST register a http POST router of handler.
func (router *Router) POST(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("POST", relativePath, handler, decorators...))
	router.currentRouter.POST(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// PUT register a http PUT router of handler.
func (router *Router) PUT(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("PUT", relativePath, handler, decorators...))
	router.currentRouter.PUT(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// HEAD register a http HEAD router of handler.
func (router *Router) HEAD(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("HEAD", relativePath, handler, decorators...))
	router.currentRouter.HEAD(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// DELETE register a http DELETE router of handler.
func (router *Router) DELETE(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("DELETE", relativePath, handler, decorators...))
	router.currentRouter.DELETE(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// OPTIONS register a http OPTIONS router of handler.
func (router *Router) OPTIONS(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("OPTIONS", relativePath, handler, decorators...))
	router.currentRouter.OPTIONS(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// PATCH register a http PATCH router of handler.
func (router *Router) PATCH(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("PATCH", relativePath, handler, decorators...))
	router.currentRouter.PATCH(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
		handle(c, handler, a, b)
	})
}

// Any register a any HTTP method router of handler.
func (router *Router) Any(relativePath string, handler Handler, decorators ...IDecorator) {
	a, b := groupDecorators(decorators...)
	router.routerInfos = append(router.routerInfos, createRouterInfo("Any", relativePath, handler, decorators...))
	router.currentRouter.Any(relativePath, func(c *gin.Context) {
		c.Set(apiNameKey, relativePath)
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
func (router *Router) Group(relativePath string, handler Handler, decorators ...IDecorator) *RouterGroup {
	rg := &RouterGroup{
		router,
	}
	if handler != nil {
		a, b := groupDecorators(decorators...)
		router.routerInfos = append(router.routerInfos, createRouterInfo("Any Group", relativePath, handler, decorators...))
		rg.currentRouter = rg.router.Group(relativePath, func(c *gin.Context) {
			handle(c, handler, a, b)
		})
	} else {
		rg.currentRouter = rg.router.Group(relativePath)
	}
	return rg
}

// RegisterRestAPIs register a collection HTTP routes by gw.IRestAPI.
func (router *RouterGroup) RegisterRestAPIs(ctrls ...IRestAPI) {
	RegisterRestAPI(router, ctrls...)
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

func RegisterRestAPI(router *RouterGroup, restList ...IRestAPI) {
	logger.Info("register router by API RegisterRestAPI(...)")
	for _, ctrl := range restList {
		typ := reflect.TypeOf(ctrl)
		val := reflect.ValueOf(ctrl)
		var relativePath string
		nameCaller, ok := typ.MethodByName("Name")
		if ok {
			relativePath = nameCaller.Func.Call([]reflect.Value{val})[0].String()
		}
		nm := typ.NumMethod()
		for i := 0; i < nm; i++ {
			m := typ.Method(i)
			h, ok := methodParsers[strings.ToLower(m.Name)]
			if ok {
				// FIXME(Ocean): how to check the arguments type is *gw.Context.
				n := 1
				prefix := fmt.Sprintf("invalid operation, method:%s.%s", val.Interface(), m.Name)
				if m.Type.NumOut() != 0 {
					panic(fmt.Sprintf("%s, should be not return any values.", prefix))
				}
				dynBinders := make([]dynamicArgsBinder, n)
				dynBinders[0] = dynamicArgsBinder{
					dataType: reflect.TypeOf(&Context{}),
					bindFunc: ctxBinder,
				}
				dynCaller := dynamicCaller{
					argsLength:        n,
					ctrl:              ctrl,
					methodName:        m.Name,
					handler:           val.MethodByName(m.Name),
					argsOrderlyBinder: dynBinders,
				}
				h(relativePath, router, dynCaller)
			}
		}
	}
}

var methodParsers map[string]dynamicHandler

func init() {
	methodParsers = make(map[string]dynamicHandler)
	methodParsers["get"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.GET(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["query"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		relativePath = strings.TrimRight(relativePath, "/")
		relativePath = fmt.Sprintf("%s/query", relativePath)
		r.GET(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["post"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.POST(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["put"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.PUT(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["delete"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.DELETE(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["options"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.OPTIONS(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["option"] = methodParsers["options"]
	methodParsers["patch"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.PATCH(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["head"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.HEAD(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["any"] = func(relativePath string, r *RouterGroup, caller dynamicCaller) {
		r.Any(relativePath, func(ctx *Context) {
			handleDynamic(ctx, caller)
		})
	}
	methodParsers["all"] = methodParsers["any"]
}

func handleDynamic(ctx *Context, caller dynamicCaller) {
	caller.handler.Call(caller.makeArgs(ctx))
}

func (d dynamicCaller) makeArgs(ctx *Context) []reflect.Value {
	if d.argsLength > 0 {
		var args = make([]reflect.Value, d.argsLength)
		for i := 0; i < d.argsLength; i++ {
			binder := d.argsOrderlyBinder[i]
			args[i] = binder.bindFunc(binder.dataType, ctx)
		}
		return args
	}
	return nil
}

func ctxBinder(typ reflect.Type, ctx *Context) reflect.Value {
	return reflect.ValueOf(ctx)
}

func handle(c *gin.Context, handler Handler, beforeDecorators, afterDecorators []IDecorator) {
	s := hostServer(c)
	user := getUser(c)
	requestID := getRequestID(s, c)
	store := getStore(c, s, *user)
	ctx := makeCtx(c, *user, *s, store, requestID)
	// action before Decorators
	var msg string
	var err error
	for _, d := range beforeDecorators {
		msg, err = d.Call(ctx)
		if err != nil {
			break
		}
	}
	if err != nil {
		if msg == "" {
			msg = "caller decorator fail."
		}
		body := respBody(http.StatusForbidden, requestID, errDefault403Msg, msg)
		c.JSON(http.StatusForbidden, body)
		return
	}

	// process Action handler.
	handler(ctx)

	// action after Decorators
	for _, d := range afterDecorators {
		msg, err = d.Call(ctx)
		if err != nil {
			break
		}
	}
	if err != nil {
		if msg == "" {
			msg = "caller decorator fail."
		}
		body := respBody(http.StatusForbidden, requestID, errDefault403Msg, msg)
		c.JSON(http.StatusForbidden, body)
		return
	}
}

func makeCtx(c *gin.Context, user User, s HostServer, store Store, requestID string) *Context {
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

func groupDecorators(decorators ...IDecorator) (a, b []IDecorator) {
	var before, after []IDecorator
	for _, d := range decorators {
		if d.Point() == DecoratorPointActionBefore {
			before = append(before, d)
		} else if d.Point() == DecoratorPointActionAfter {
			after = append(after, d)
		} else {
			panic(fmt.Sprintf("invalid pointer(%d). decorator: %v", d.Point(), d))
		}
	}
	return before, after
}
