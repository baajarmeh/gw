package gw

import (
	"fmt"
	"github.com/oceanho/gw/libs/gwreflect"
	"reflect"
	"sync"
)

type ObjectTyper struct {
	Name        string
	DependOn    []TyperDependency
	Typer       reflect.Type
	ActualValue reflect.Value
	newAPI      reflect.Value
	IsPtr       bool
}

const (
	IStoreName               = "github.com/oceanho/gw.IStore"
	IPasswordSignerName      = "github.com/oceanho/gw.IPasswordSigner"
	IPermissionManagerName   = "github.com/oceanho/gw.IPermissionManager"
	IPermissionCheckerName   = "github.com/oceanho/gw.IPermissionChecker"
	ISessionStateManagerName = "github.com/oceanho/gw.ISessionStateManager"
	IdentifierGeneratorName  = "github.com/oceanho/gw.IdentifierGenerator"
	IAuthManagerName         = "github.com/oceanho/gw.IAuthManager"
	IUserManagerName         = "github.com/oceanho/gw.IUserManager"
	ServerStateName          = "github.com/oceanho/gw.ServerState"
	HostServerName           = "github.com/oceanho/gw.HostServer"
	IEventManagerName        = "github.com/oceanho/gw.IEventManager"
)

var nullReflectValue = reflect.ValueOf(nil)

func (ss ServerState) objectTypers() map[string]ObjectTyper {
	var typers = make(map[string]ObjectTyper)
	typers[IStoreName] = newNilApiObjectTyper(IStoreName, ss.Store())
	typers[HostServerName] = newNilApiObjectTyper(HostServerName, ss.s)
	typers[ServerStateName] = newNilApiObjectTyper(ServerStateName, ss)
	typers[IUserManagerName] = newNilApiObjectTyper(IUserManagerName, ss.UserManager())
	typers[IAuthManagerName] = newNilApiObjectTyper(IAuthManagerName, ss.AuthManager())
	typers[IEventManagerName] = newNilApiObjectTyper(IEventManagerName, ss.EventManager())
	typers[IPasswordSignerName] = newNilApiObjectTyper(IPasswordSignerName, ss.PasswordSigner())
	typers[IdentifierGeneratorName] = newNilApiObjectTyper(IdentifierGeneratorName, ss.IDGenerator())
	typers[IPermissionManagerName] = newNilApiObjectTyper(IPermissionManagerName, ss.PermissionManager())
	typers[IPermissionCheckerName] = newNilApiObjectTyper(IPermissionCheckerName, ss.PermissionChecker())
	typers[ISessionStateManagerName] = newNilApiObjectTyper(ISessionStateManagerName, ss.SessionStateManager())
	return typers
}

func newNilApiObjectTyper(name string, value interface{}) ObjectTyper {
	return ObjectTyper{
		Name:        name,
		IsPtr:       false,
		newAPI:      nullReflectValue,
		Typer:       reflect.TypeOf(value),
		ActualValue: reflect.ValueOf(value),
	}
}

type TyperDependency struct {
	Name  string
	IsPtr bool
	Typer reflect.Type
}

type DIConfig struct {
	NewFuncName string
	ResolveFunc func(typers *map[string]ObjectTyper, state interface{}, typerName string) interface{}
}

var defaultDIConfig = DIConfig{
	NewFuncName: "New",
	ResolveFunc: func(typers *map[string]ObjectTyper, state interface{}, typerName string) interface{} {
		var store = state.(IStore)
		var objectTyper, ok = (*typers)[typerName]
		if !ok {
			panic(fmt.Sprintf("object typer(%s) not found", typerName))
		}
		var result reflect.Value
		if objectTyper.newAPI == nullReflectValue {
			result = objectTyper.ActualValue
		} else {
			var values []reflect.Value
			for _, typerDp := range objectTyper.DependOn {
				if typerDp.Name == IStoreName && store != nil {
					values = append(values, reflect.ValueOf(store))
				} else {
					values = append(values, resolver(typers, typerDp, store))
				}
			}
			result = objectTyper.newAPI.Call(values)[0]
		}
		return result.Interface()
	},
}

type IDIProvider interface {
	Register(actual ...interface{}) bool
	RegisterWithTyper(typers ...reflect.Type) bool
	RegisterWithName(typerName string, actual interface{}) bool
	Resolve(typerName string) interface{}
	ResolveByTyper(typer reflect.Type) interface{}
	ResolveWithState(state interface{}, typerName string) interface{}
	ResolveByTyperWithState(state interface{}, typer reflect.Type) interface{}
}

type DefaultDIProviderImpl struct {
	locker       sync.Mutex
	config       DIConfig
	state        *ServerState
	objectTypers map[string]ObjectTyper
}

func DefaultDIProvider(state *ServerState) IDIProvider {
	var di = &DefaultDIProviderImpl{
		state:        state,
		objectTypers: state.objectTypers(),
		config:       defaultDIConfig,
	}
	return di
}

func (d *DefaultDIProviderImpl) Register(actual ...interface{}) bool {
	for _, a := range actual {
		if !d.RegisterWithName("", a) {
			return false
		}
	}
	return true
}

func (d *DefaultDIProviderImpl) RegisterWithTyper(typers ...reflect.Type) bool {
	for _, typer := range typers {
		if !d.RegisterWithName("", reflect.New(typer).Interface()) {
			return false
		}
	}
	return true
}

func (d *DefaultDIProviderImpl) RegisterWithName(name string, actual interface{}) bool {
	var actualValue = reflect.ValueOf(actual)
	var newMethod = actualValue.MethodByName(d.config.NewFuncName)
	if newMethod.Kind() != reflect.Func {
		panic(fmt.Sprintf("typer(%s) has no %s(...) APIs.",
			gwreflect.GetPkgFullName(reflect.TypeOf(actual)), d.config.NewFuncName))
	}
	var newMethodTyper = newMethod.Type()
	var newMethodNumIn = newMethodTyper.NumIn()
	var newMethodNumOut = newMethodTyper.NumOut()
	if newMethodNumOut != 1 {
		panic(fmt.Errorf("actual(%s) named typer New(...) func should be has only one return values", actual))
	}
	var outTyper = newMethodTyper.Out(0)
	var virtualName = gwreflect.GetPkgFullName(outTyper)
	if name != "" {
		virtualName = name
	}
	var newFunParamObjectTypers = make([]TyperDependency, newMethodNumIn)
	for idx := 0; idx < newMethodNumIn; idx++ {
		newFunParamObjectTypers[idx] = typeToObjectTyperName(newMethodTyper.In(idx))
	}
	var objectTyper = ObjectTyper{
		newAPI:      newMethod,
		Name:        virtualName,
		DependOn:    newFunParamObjectTypers,
		ActualValue: actualValue,
		Typer:       reflect.TypeOf(actual),
	}
	d.locker.Lock()
	defer d.locker.Unlock()
	d.objectTypers[virtualName] = objectTyper
	return true
}

func (d *DefaultDIProviderImpl) Resolve(typerName string) interface{} {
	return d.ResolveWithState(d.state.Store(), typerName)
}

func (d *DefaultDIProviderImpl) ResolveByTyper(typer reflect.Type) interface{} {
	return d.ResolveByTyperWithState(d.state.Store(), typer)
}

func (d *DefaultDIProviderImpl) ResolveByTyperWithState(state interface{}, typer reflect.Type) interface{} {
	var typerName = ""
	if method, ok := typer.MethodByName(d.config.NewFuncName); ok {
		typerName = gwreflect.GetPkgFullName(method.Func.Type().Out(0))
	} else {
		typerName = gwreflect.GetPkgFullName(typer)
	}
	return d.ResolveWithState(state, typerName)
}

func (d *DefaultDIProviderImpl) ResolveWithState(state interface{}, typerName string) interface{} {
	return d.config.ResolveFunc(&d.objectTypers, state, typerName)
}

// helpers
func resolver(typers *map[string]ObjectTyper, typerDependency TyperDependency, store IStore) reflect.Value {
	var values []reflect.Value
	var objectTyper, ok = (*typers)[typerDependency.Name]
	if !ok {
		panic(fmt.Sprintf("missing typer(%s)", typerDependency.Name))
	}
	if objectTyper.newAPI == nullReflectValue {
		return objectTyper.ActualValue
	}
	for _, dp := range objectTyper.DependOn {
		dp := dp
		if dp.Name == IStoreName && store != nil {
			values = append(values, reflect.ValueOf(store))
		} else if _, ok := (*typers)[dp.Name]; ok {
			values = append(values, resolver(typers, dp, store))
		} else {
			panic(fmt.Sprintf("object typer(%s) not found", dp.Name))
		}
	}
	var result = objectTyper.newAPI.Call(values)[0]
	switch result.Kind() {
	case reflect.Ptr:
		if typerDependency.IsPtr {
			return result
		} else {
			return reflect.ValueOf(result.Elem().Interface())
		}
	default:
		if !typerDependency.IsPtr {
			return result
		} else {
			var typeValue = reflect.New(reflect.TypeOf(result.Interface()))
			typeValue.Elem().Set(result)
			return typeValue
		}
	}
}

func typeToObjectTyperName(typer reflect.Type) TyperDependency {
	return TyperDependency{
		Typer: typer,
		Name:  gwreflect.GetPkgFullName(typer),
		IsPtr: typer.Kind() == reflect.Ptr,
	}
}
