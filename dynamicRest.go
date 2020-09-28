package gw

import (
	"fmt"
	"github.com/oceanho/gw/logger"
	"net/http"
	"reflect"
	"strings"
)

var (
	restApiRegister map[string]restHandler
)

func init() {
	restApiRegister = make(map[string]restHandler)
	restApiRegister["get"] = restHandler{
		"Get",
		"Get",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter("GET", relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["detail"] = restHandler{
		"Get",
		"Detail",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/detail/:id", relativePath)
			r.createRouter("GET", relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["query"] = restHandler{
		"Get",
		"Query",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/query", relativePath)
			r.createRouter(http.MethodGet, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["querylist"] = restHandler{
		"Get",
		"QueryList",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			relativePath = strings.TrimRight(relativePath, "/")
			relativePath = fmt.Sprintf("%s/queryList", relativePath)
			r.createRouter(http.MethodGet, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["post"] = restHandler{
		"Post",
		"Post",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodPost, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["put"] = restHandler{
		"Put",
		"Put",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodPut, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["delete"] = restHandler{
		"Delete",
		"Delete",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodDelete, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["options"] = restHandler{
		"Options",
		"Options",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodOptions, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["patch"] = restHandler{
		"Patch",
		"Patch",
		func(relativePath, actionPkgName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			var handlerActionName = fmt.Sprintf("%s.Patch", actionPkgName)
			r.createRouter(http.MethodPatch, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["head"] = restHandler{
		"Head",
		"Head",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodHead, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["connect"] = restHandler{
		"Connect",
		"Connect",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodConnect, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["trace"] = restHandler{
		"Trace",
		"Trace",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter(http.MethodTrace, relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["any"] = restHandler{
		"Any",
		"Any",
		func(relativePath, handlerActionName string, r *RouterGroup, dynamicCaller DynamicCaller) {
			r.createRouter("any", relativePath, func(ctx *Context) {
				handleDynamicApi(ctx, dynamicCaller)
			}, handlerActionName, dynamicCaller.decorators...)
		},
	}
	restApiRegister["all"] = restApiRegister["any"]
}

// restHandler
type restHandler struct {
	httpMethod     string
	actionFuncName string
	register       func(relativePath, actionPkgFunName string, r *RouterGroup, dynamicCaller DynamicCaller)
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

// RegisterRestAPIs register a collection HTTP routes by gw.IDynamicRestAPI.
func (router *RouterGroup) RegisterRestAPIs(restAPIs ...IDynamicRestAPI) {
	registerRestAPIImpl(router, restAPIs...)
}

func registerRestAPIImpl(router *RouterGroup, restAPIs ...IDynamicRestAPI) {
	logger.Info("register router by API RegisterRestAPI(...)")
	for i := 0; i < len(restAPIs); i++ {
		rest := restAPIs[i]
		var restPkgId string
		typ := reflect.TypeOf(rest)
		val := reflect.ValueOf(rest)
		ctrlCallArgs := []reflect.Value{
			reflect.ValueOf(rest),
		}
		if typ.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("%s should be are pointer.", rest.Name()))
		}
		el := typ.Elem()
		restPkgId = fmt.Sprintf("%s.(*%s)", el.PkgPath(), el.Name())
		relativePath := strings.ToLower(val.MethodByName("Name").Call(nil)[0].String())
		var name = "SetupDecorator"
		var globalDecorators = make([]Decorator, 0, 0)
		if _, ok := typ.MethodByName(name); ok {
			d := val.MethodByName(name).Call(nil)[0].Interface()
			if r, ok := d.([]Decorator); ok {
				globalDecorators = append(globalDecorators, r...)
			} else if r, ok := d.(Decorator); ok {
				globalDecorators = append(globalDecorators, r)
			} else if r, ok := d.(*Decorator); ok {
				globalDecorators = append(globalDecorators, *r)
			} else if r, ok := d.([]*Decorator); ok {
				for _, _r := range r {
					_r1 := *_r
					globalDecorators = append(globalDecorators, _r1)
				}
			}
		}
		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			dyApiRegister, ok := restApiRegister[strings.ToLower(m.Name)]
			if ok {
				var apiSpecifyDecorators = make([]Decorator, 0, 0)
				name = "SetupOn" + m.Name + "Decorator"
				_, ok := typ.MethodByName(name)
				if ok {
					d := val.MethodByName(name).Call(nil)[0].Interface()
					if r, ok := d.([]Decorator); ok {
						apiSpecifyDecorators = append(apiSpecifyDecorators, r...)
					} else if r, ok := d.(Decorator); ok {
						apiSpecifyDecorators = append(apiSpecifyDecorators, r)
					} else if r, ok := d.(*Decorator); ok {
						apiSpecifyDecorators = append(apiSpecifyDecorators, *r)
					} else if r, ok := d.([]*Decorator); ok {
						for _, _r := range r {
							_r1 := *_r
							apiSpecifyDecorators = append(apiSpecifyDecorators, _r1)
						}
					}
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
				// OnXBefore
				name = fmt.Sprintf("On%sBefore", m.Name)
				onBefore, ok := typ.MethodByName(name)
				if ok {
					onBeforeHandler := onBefore.Func.Call(ctrlCallArgs)[0].Interface()
					if handler, ok := onBeforeHandler.(DecoratorHandler); ok {
						decorators = append(decorators, Decorator{
							Before: handler,
						})
					} else if handler, ok := onBeforeHandler.([]DecoratorHandler); ok {
						for _, h := range handler {
							h := h
							decorators = append(decorators, Decorator{
								Before: h,
							})
						}
					}
				}

				decorators = append(decorators, apiSpecifyDecorators...)
				decorators = append(decorators, globalDecorators...)

				// OnXAfter
				name = fmt.Sprintf("On%sAfter", m.Name)
				onAfter, ok := typ.MethodByName(name)
				if ok {
					onAfterHandler := onAfter.Func.Call(ctrlCallArgs)[0].Interface()
					if handler, ok := onAfterHandler.(DecoratorHandler); ok {
						decorators = append(decorators, Decorator{
							After: handler,
						})
					} else if handler, ok := onAfterHandler.([]DecoratorHandler); ok {
						for _, h := range handler {
							h := h
							decorators = append(decorators, Decorator{
								After: h,
							})
						}
					}
				}
				bindingFuncPkgName := fmt.Sprintf("%s.%s", restPkgId, m.Name)
				dynCaller := DynamicCaller{
					argInNumber:        n,
					decorators:         decorators,
					bindingFuncPkgName: bindingFuncPkgName,
					handler:            val.MethodByName(m.Name),
					argsOrderlyBinder:  dynBinders,
				}
				dyApiRegister.register(relativePath, bindingFuncPkgName, router, dynCaller)
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
