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
	ICryptoHashName          = "github.com/oceanho/gw.ICryptoHash"
	ICryptoProtectName       = "github.com/oceanho/gw.ICryptoProtect"
	IdentifierGeneratorName  = "github.com/oceanho/gw.IdentifierGenerator"
	IAuthManagerName         = "github.com/oceanho/gw.IAuthManager"
	IUserManagerName         = "github.com/oceanho/gw.IUserManager"
	ServerStateName          = "github.com/oceanho/gw.ServerState"
	HostServerName           = "github.com/oceanho/gw.HostServer"
	IEventManagerName        = "github.com/oceanho/gw.IEventManager"
)

var (
	NullReflectValue      = reflect.ValueOf(nil)
	BuiltinComponentTyper = reflect.TypeOf(BuiltinComponent{})
)

func (ss ServerState) objectTypers() map[string]ObjectTyper {
	var typers = make(map[string]ObjectTyper)
	typers[IStoreName] = newNilApiObjectTyper(IStoreName, ss.Store())
	typers[HostServerName] = newNilApiObjectTyper(HostServerName, ss.s)
	typers[ServerStateName] = newNilApiObjectTyper(ServerStateName, ss)
	typers[IUserManagerName] = newNilApiObjectTyper(IUserManagerName, ss.UserManager())
	typers[IAuthManagerName] = newNilApiObjectTyper(IAuthManagerName, ss.AuthManager())
	typers[ICryptoHashName] = newNilApiObjectTyper(ICryptoHashName, ss.CryptoHash())
	typers[ICryptoProtectName] = newNilApiObjectTyper(ICryptoProtectName, ss.CryptoProtect())
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
		newAPI:      NullReflectValue,
		Typer:       reflect.TypeOf(value),
		ActualValue: reflect.ValueOf(value),
	}
}

type BuiltinComponent struct {
	Store               IStore
	UserManager         IUserManager
	AuthManager         IAuthManager
	SessionStateManager ISessionStateManager
	PermissionManager   IPermissionManager
	PermissionChecker   IPermissionChecker
	CryptoHash          ICryptoHash
	CryptoProtect       ICryptoProtect
	PasswordSigner      IPasswordSigner
	IDGenerator         IdentifierGenerator
}

func (bc BuiltinComponent) New(
	Store IStore,
	UserManager IUserManager,
	AuthManager IAuthManager,
	SessionStateManager ISessionStateManager,
	PermissionManager IPermissionManager,
	PermissionChecker IPermissionChecker,
	CryptoHash ICryptoHash,
	CryptoProtect ICryptoProtect,
	PasswordSigner IPasswordSigner,
	IDGenerator IdentifierGenerator,
) BuiltinComponent {
	bc.Store = Store
	bc.UserManager = UserManager
	bc.AuthManager = AuthManager
	bc.SessionStateManager = SessionStateManager
	bc.PermissionManager = PermissionManager
	bc.PermissionChecker = PermissionChecker
	bc.CryptoHash = CryptoHash
	bc.CryptoProtect = CryptoProtect
	bc.PasswordSigner = PasswordSigner
	bc.IDGenerator = IDGenerator
	return bc
}

type TyperDependency struct {
	Name  string
	IsPtr bool
	Typer reflect.Type
}

type DIConfig struct {
	NewFuncName string
	ResolveFunc func(di interface{}, state interface{}, typerName string) interface{}
}

var defaultDIConfig = DIConfig{
	NewFuncName: "New",
	ResolveFunc: func(diImpl interface{}, state interface{}, typerName string) interface{} {
		var di = diImpl.(*DefaultDIProviderImpl)
		if di == nil {
			panic("diImpl not are DefaultDIProviderImpl")
		}
		var store interface{}
		if state != nil {
			store = state.(IStore)
		}
		var objectTyper, ok = di.objectTypers[typerName]
		if !ok {
			panic(fmt.Sprintf("object typer(%s) not found", typerName))
		}
		var result reflect.Value
		if objectTyper.newAPI == NullReflectValue {
			result = objectTyper.ActualValue
		} else {
			var values []reflect.Value
			for _, typerDp := range objectTyper.DependOn {
				if typerDp.Name == IStoreName && store != nil {
					values = append(values, reflect.ValueOf(store))
				} else {
					values = append(values, resolver(di, typerDp, store))
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
	typerMappers map[reflect.Type]string
}

func DefaultDIProvider(state *ServerState) IDIProvider {
	stateTypers := state.objectTypers()
	var registeredTypers = make(map[reflect.Type]string)
	for _, d := range stateTypers {
		registeredTypers[d.Typer] = d.Name
	}
	var di = &DefaultDIProviderImpl{
		state:        state,
		objectTypers: stateTypers,
		typerMappers: registeredTypers,
		config:       defaultDIConfig,
	}
	di.Register(BuiltinComponent{})
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
	d.typerMappers[outTyper] = virtualName
	return true
}

func (d *DefaultDIProviderImpl) Resolve(typerName string) interface{} {
	return d.ResolveWithState(d.state.Store(), typerName)
}

func (d *DefaultDIProviderImpl) ResolveByTyper(typer reflect.Type) interface{} {
	return d.ResolveByTyperWithState(d.state.Store(), typer)
}

func (d *DefaultDIProviderImpl) ResolveByTyperWithState(state interface{}, typer reflect.Type) interface{} {
	var typerName, ok = d.typerMappers[typer]
	if !ok {
		typerName = gwreflect.GetPkgFullName(typer)
	}
	return d.ResolveWithState(state, typerName)
}

func (d *DefaultDIProviderImpl) ResolveWithState(state interface{}, typerName string) interface{} {
	return d.config.ResolveFunc(d, state, typerName)
}

// helpers
func resolver(defaultDiImpl *DefaultDIProviderImpl, typerDependency TyperDependency, state interface{}) reflect.Value {
	var store interface{}
	if state != nil {
		store = state.(IStore)
	}
	var values []reflect.Value
	var objectTyper, ok = defaultDiImpl.objectTypers[typerDependency.Name]
	if !ok {
		panic(fmt.Sprintf("missing typer(%s)", typerDependency.Name))
	}
	if objectTyper.newAPI == NullReflectValue {
		return objectTyper.ActualValue
	}
	for _, dp := range objectTyper.DependOn {
		dp := dp
		if dp.Name == IStoreName && store != nil {
			values = append(values, reflect.ValueOf(store))
		} else if _, ok := defaultDiImpl.objectTypers[dp.Name]; ok {
			values = append(values, resolver(defaultDiImpl, dp, store))
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
