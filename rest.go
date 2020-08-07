package gw

import (
	"fmt"
	"github.com/oceanho/gw/logger"
	"net/http"
	"path"
	"reflect"
	"strings"
)

var (
	restApiRegister                 map[string]restHandler
	errorDynamicCallerBeforeHandler = fmt.Errorf("dnamic Caller Before Handler fail")
	errorDynamicCallerAfterHandler  = fmt.Errorf("dnamic Caller After Handler fail")
)

func init() {
	restApiRegister = make(map[string]restHandler)
	restApiRegister["get"] = restHandler{
		"Get",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}

	restApiRegister["query"] = restHandler{
		"Get",
		"query",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/query", relativePath)
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["queryList"] = restHandler{
		"Get",
		"queryList",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/queryList", relativePath)
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["post"] = restHandler{
		"Post",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.POST(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["put"] = restHandler{
		"Put",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.PUT(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["delete"] = restHandler{
		"Delete",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.DELETE(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["options"] = restHandler{
		"Options",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.OPTIONS(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["option"] = restApiRegister["options"]
	restApiRegister["patch"] = restHandler{
		"Patch",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.PATCH(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["head"] = restHandler{
		"Head",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.HEAD(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["any"] = restHandler{
		"Any",
		"",
		func(relativePath string, r *RouterGroup, caller restCaller) {
			r.Any(relativePath, func(ctx *Context) {
				handleDynamicRestApi(ctx, caller)
			})
		},
	}
	restApiRegister["all"] = restApiRegister["any"]
}

// restHandler
type restHandler struct {
	httpMethod  string
	extraRouter string
	register    func(relativePath string, r *RouterGroup, caller restCaller)
}

// restCaller
type restCaller struct {
	rest                   IRestAPI
	argsNumber             int
	handler                reflect.Value
	hasActionBeforeHandler bool
	hasActionAfterHandler  bool
	beforeHandler          reflect.Value
	afterHandler           reflect.Value
	beforeDecorators       []Decorator
	afterDecorators        []Decorator
	afterDecoratorsMaxIdx  int
	handlerActionName      string
	argsOrderlyBinder      []restArgsBinder
}

type restArgsBinder struct {
	dataType reflect.Type
	bindFunc func(p reflect.Type, c *Context) reflect.Value
}

// RegisterRestAPIs register a collection HTTP routes by gw.IRestAPI.
func (router *RouterGroup) RegisterRestAPIs(ctrls ...IRestAPI) {
	RegisterRestAPI(router, ctrls...)
}

func RegisterRestAPI(router *RouterGroup, restAPIs ...IRestAPI) {
	logger.Info("register router by API RegisterRestAPI(...)")
	for _, rest := range restAPIs {
		var relativePath, restName string
		typ := reflect.TypeOf(rest)
		val := reflect.ValueOf(rest)
		if typ.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("%s should be are pointer.", rest.Name()))
		}
		restName = fmt.Sprintf("%s.%s", typ.Elem().PkgPath(), typ.Elem().Name())
		nameCaller, ok := typ.MethodByName("Name")
		if ok {
			relativePath = nameCaller.Func.Call([]reflect.Value{val})[0].String()
		}
		var name = "SetupDecorator"
		var restDecorators []Decorator
		_, ok = typ.MethodByName(name)
		if ok {
			restDecorators = val.MethodByName(name).Call(nil)[0].Interface().([]Decorator)
		}
		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			dyApiRegister, ok := restApiRegister[strings.ToLower(m.Name)]
			if ok {
				var actionDecorators []Decorator
				name = "SetupOn" + m.Name + "Decorator"
				dm, ok := typ.MethodByName(name)
				if ok {
					actionDecorators = dm.Func.Call(nil)[0].Interface().([]Decorator)
				}
				// FIXME(Ocean): how to check the arguments type is *gw.Context.
				n := 1
				prefix := fmt.Sprintf("invalid operation, method: %s.%s", restName, m.Name)
				if m.Type.NumOut() != 0 {
					panic(fmt.Sprintf("%s, should be not return any values.", prefix))
				}
				dynBinders := make([]restArgsBinder, n)
				dynBinders[0] = restArgsBinder{
					dataType: reflect.TypeOf(&Context{}),
					bindFunc: ctxBinder,
				}
				var decorators []Decorator
				decorators = append(decorators, restDecorators...)
				decorators = append(decorators, actionDecorators...)
				before, after := splitDecorators(decorators...)
				handlerActionName := fmt.Sprintf("%s.%s(*gw.Context)", restName, m.Name)
				dynCaller := restCaller{
					argsNumber:             n,
					beforeDecorators:       before,
					afterDecorators:        after,
					afterDecoratorsMaxIdx:  len(after) - 1,
					rest:                   rest,
					handlerActionName:      handlerActionName,
					handler:                val.MethodByName(m.Name),
					argsOrderlyBinder:      dynBinders,
					hasActionBeforeHandler: false,
					hasActionAfterHandler:  false,
					beforeHandler:          reflect.ValueOf(nil),
					afterHandler:           reflect.ValueOf(nil),
				}
				name = "On" + m.Name + "Before"
				if _, ok := typ.MethodByName(name); ok {
					dynCaller.hasActionBeforeHandler = true
					dynCaller.beforeHandler = val.MethodByName(name)
				}
				name = "On" + m.Name + "After"
				if _, ok := typ.MethodByName(name); ok {
					dynCaller.hasActionAfterHandler = true
					dynCaller.afterHandler = val.MethodByName(name)
				}
				url := path.Join(relativePath, dyApiRegister.extraRouter)
				router.storeRouterStateWithHandlerName(
					strings.ToUpper(dyApiRegister.httpMethod), url, handlerActionName, nil, decorators...)
				dyApiRegister.register(relativePath, router, dynCaller)
			}
		}
	}
}

func handleDynamicRestApi(c *Context, caller restCaller) {
	args := caller.makeArgs(c)
	// XXX handler AT first call.
	s := hostServer(c.Context)
	requestID := getRequestID(s, c.Context)
	var err error
	var msg string
	var isParserOK = false
	if caller.hasActionBeforeHandler {
		returns := caller.beforeHandler.Call(args)
		if len(returns) == 1 {
			msg, isParserOK = returns[0].Interface().(string)
			if !isParserOK {
				err, isParserOK = returns[0].Interface().(error)
				if !isParserOK {
					logger.Warn("rest caller before handler %s "+
						"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
						caller.handlerActionName)
					err = errorDynamicCallerBeforeHandler
				}
			}
		} else if len(returns) == 2 {
			msg, isParserOK = returns[0].Interface().(string)
			if !isParserOK {
				logger.Warn("rest caller before handler %s "+
					"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
					caller.handlerActionName)
			}
			err, isParserOK = returns[1].Interface().(error)
			if !isParserOK {
				logger.Warn("rest caller before handler %s "+
					"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
					caller.handlerActionName)
				err = errorDynamicCallerBeforeHandler
			}
		}
		if err != nil {
			if msg == "" {
				msg = "rest call before handler fail."
			}
			body := respBody(http.StatusBadRequest, requestID, errDefault400Msg, msg)
			c.JSON(http.StatusBadRequest, body)
			return
		}
	}
	// before decorators
	msg = ""
	err = nil
	for _, d := range caller.beforeDecorators {
		msg, err = d.Before(c)
		if err != nil {
			break
		}
	}
	if err != nil {
		if msg == "" {
			msg = "rest call before decorator fail."
		}
		body := respBody(http.StatusBadRequest, requestID, errDefault400Msg, msg)
		c.JSON(http.StatusBadRequest, body)
		return
	}

	// rest handler
	caller.handler.Call(caller.makeArgs(c))

	// after decorators
	if caller.afterDecoratorsMaxIdx >= 0 {
		msg = ""
		err = nil
		for i := caller.afterDecoratorsMaxIdx; i >= 0; i-- {
			msg, err = caller.afterDecorators[i].After(c)
			if err != nil {
				break
			}
		}
		if err != nil {
			if msg == "" {
				msg = "rest call before decorator fail."
			}
			body := respBody(http.StatusBadRequest, requestID, errDefault400Msg, msg)
			c.JSON(http.StatusBadRequest, body)
			return
		}
	}
	// after caller handler
	msg = ""
	err = nil
	if caller.hasActionAfterHandler {
		returns := caller.afterHandler.Call(args)
		if len(returns) == 1 {
			msg, isParserOK = returns[0].Interface().(string)
			if !isParserOK {
				err, isParserOK = returns[0].Interface().(error)
				if !isParserOK {
					logger.Warn("rest caller after handler %s "+
						"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
						caller.handlerActionName)
					err = errorDynamicCallerAfterHandler
				}
			}
		} else if len(returns) == 2 {
			msg, isParserOK = returns[0].Interface().(string)
			if !isParserOK {
				logger.Warn("rest caller after handler %s "+
					"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
					caller.handlerActionName)
			}
			err, isParserOK = returns[1].Interface().(error)
			if !isParserOK {
				logger.Warn("rest caller before handler %s "+
					"return value are invalid. should be returns as (msg string, err error)/(string)/(error)",
					caller.handlerActionName)
				err = errorDynamicCallerAfterHandler
			}
		}
		if err != nil {
			if msg == "" {
				msg = "rest call after handler fail."
			}
			body := respBody(http.StatusBadRequest, requestID, errDefault400Msg, msg)
			c.JSON(http.StatusBadRequest, body)
			return
		}
	}
}

func (d restCaller) makeArgs(ctx *Context) []reflect.Value {
	if d.argsNumber > 0 {
		var args = make([]reflect.Value, d.argsNumber)
		for i := 0; i < d.argsNumber; i++ {
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
