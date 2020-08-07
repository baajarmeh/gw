package gw

import (
	"fmt"
	"github.com/oceanho/gw/logger"
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
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}

	restApiRegister["query"] = restHandler{
		"Get",
		"query",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/query", relativePath)
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["queryList"] = restHandler{
		"Get",
		"queryList",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/queryList", relativePath)
			r.GET(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["post"] = restHandler{
		"Post",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.POST(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["put"] = restHandler{
		"Put",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.PUT(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["delete"] = restHandler{
		"Delete",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.DELETE(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["options"] = restHandler{
		"Options",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.OPTIONS(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["option"] = restApiRegister["options"]
	restApiRegister["patch"] = restHandler{
		"Patch",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.PATCH(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["head"] = restHandler{
		"Head",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.HEAD(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["any"] = restHandler{
		"Any",
		"",
		func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.Any(relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, dynamicCaller.decorators...)
		},
	}
	restApiRegister["all"] = restApiRegister["any"]
}

// restHandler
type restHandler struct {
	httpMethod  string
	extraRouter string
	register    func(relativePath string, r *RouterGroup, dynamicCaller DynamicCaller)
}

// DynamicCaller ...
type DynamicCaller struct {
	argInNumber        int
	retOutNumber       int
	handler            reflect.Value
	decorators         []Decorator
	bindingFuncPkgName string
	argsOrderlyBinder  []restArgsBinder
}

type restArgsBinder struct {
	dataType reflect.Type
	bindFunc func(p reflect.Type, c *Context) reflect.Value
}

// RegisterRestAPIs register a collection HTTP routes by gw.IRestAPI.
func (router *RouterGroup) RegisterRestAPIs(restAPIs ...IRestAPI) {
	RegisterRestAPI(router, restAPIs...)
}

func RegisterRestAPI(router *RouterGroup, restAPIs ...IRestAPI) {
	logger.Info("register router by API RegisterRestAPI(...)")
	for _, rest := range restAPIs {
		var restPkgId string
		typ := reflect.TypeOf(rest)
		val := reflect.ValueOf(rest)
		if typ.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("%s should be are pointer.", rest.Name()))
		}
		el := typ.Elem()
		restPkgId = fmt.Sprintf("%s.*%s{}", el.PkgPath(), el.Name())
		relativePath := strings.ToLower(val.MethodByName("Name").Call(nil)[0].String())
		var name = "SetupDecorator"
		var globalDecorators []Decorator
		if _, ok := typ.MethodByName(name); ok {
			globalDecorators = val.MethodByName(name).Call(nil)[0].Interface().([]Decorator)
		}
		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			dyApiRegister, ok := restApiRegister[strings.ToLower(m.Name)]
			if ok {
				var apiSpecifyDecorators []Decorator
				name = "Setup" + m.Name + "Decorator"
				_, ok := typ.MethodByName(name)
				if ok {
					apiSpecifyDecorators = val.MethodByName(name).Call(nil)[0].Interface().([]Decorator)
				}
				// FIXME(Ocean): how to check the arguments type is *gw.Context.
				n := 1
				prefix := fmt.Sprintf("invalid operation, method: %s.%s", restPkgId, m.Name)
				if m.Type.NumOut() != 0 {
					panic(fmt.Sprintf("%s, should be not return any values.", prefix))
				}
				dynBinders := make([]restArgsBinder, n)
				dynBinders[0] = restArgsBinder{
					dataType: reflect.TypeOf(&Context{}),
					bindFunc: ctxBinder,
				}
				var decorators []Decorator
				decorators = append(decorators, apiSpecifyDecorators...)
				decorators = append(decorators, globalDecorators...)
				bindingFuncPkgName := fmt.Sprintf("%s.%s(*gw.Context)", restPkgId, m.Name)
				dynCaller := DynamicCaller{
					argInNumber:        n,
					decorators:         decorators,
					bindingFuncPkgName: bindingFuncPkgName,
					handler:            val.MethodByName(m.Name),
					argsOrderlyBinder:  dynBinders,
				}
				httpMethod := strings.ToUpper(dyApiRegister.httpMethod)
				httpUrlPath := path.Join(relativePath, dyApiRegister.extraRouter)
				// Save router info
				router.storeRouterStateWithHandlerName(httpMethod, httpUrlPath, bindingFuncPkgName, nil, decorators...)
				dyApiRegister.register(relativePath, router, dynCaller)
			}
		}
	}
}

// handleDynamicApi ...
func handleDynamicApi(c *Context, dynamicCaller DynamicCaller) {
	dynamicCaller.handler.Call(dynamicCaller.makeArgs(c))
}

func (d DynamicCaller) makeArgs(ctx *Context) []reflect.Value {
	if d.argInNumber > 0 {
		var args = make([]reflect.Value, d.argInNumber)
		for i := 0; i < d.argInNumber; i++ {
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
