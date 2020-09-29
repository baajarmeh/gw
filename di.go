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
	UserTyperName            = "github.com/oceanho/gw.User"
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

func (ss *ServerState) objectTypers() map[string]ObjectTyper {
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
	User                User
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
	EventManager        IEventManager
}

func (bc BuiltinComponent) New(
	User User,
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
	EventManager IEventManager,
) BuiltinComponent {
	bc.User = User
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
	bc.EventManager = EventManager
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
		di, ok := diImpl.(*DefaultDIProviderImpl)
		if !ok {
			panic("di not are gw framework DefaultDIProviderImpl instance")
		}
		var user User
		var store IStore
		ctx, ok := state.(*Context)
		if ok {
			user = ctx.User()
			store = ctx.Store()
		} else {
			ss, ok := state.(*ServerState)
			if !ok {
				panic("state are not a typer of *gw.Context or *gw.ServerState")
			}
			user = EmptyUser
			store = ss.Store()
		}
		return resolverTyperInstance(1, di, user, store, typerName).Interface()
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
	return d.ResolveByTyperWithState(d.state, typer)
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
func resolverTyperInstance(depth int, di *DefaultDIProviderImpl, user User, store IStore, typerName string) reflect.Value {
	if depth > 16 {
		panic(fmt.Sprintf(typerName))
	}
	switch typerName {
	case IStoreName:
		if store != nil {
			return reflect.ValueOf(store)
		}
	case UserTyperName:
		return reflect.ValueOf(user)
	}
	var objectTyper, ok = di.objectTypers[typerName]
	if !ok {
		panic(fmt.Sprintf("missing typer(%s)", typerName))
	}
	if objectTyper.newAPI == NullReflectValue {
		return objectTyper.ActualValue
	}
	var result = objectTyper.newAPI.Call(resolverTyperDependOn(depth, di, objectTyper.DependOn, user, store))[0]
	switch result.Kind() {
	case reflect.Ptr:
		if objectTyper.IsPtr {
			return result
		} else {
			return reflect.ValueOf(result.Elem().Interface())
		}
	default:
		if !objectTyper.IsPtr {
			return result
		} else {
			var typeValue = reflect.New(reflect.TypeOf(result.Interface()))
			typeValue.Elem().Set(result)
			return typeValue
		}
	}
}

func resolverTyperDependOn(depth int, di *DefaultDIProviderImpl, dependencies []TyperDependency, user User, store IStore) []reflect.Value {
	var values = make([]reflect.Value, 0, 8)
	for _, dp := range dependencies {
		values = append(values, resolverTyperInstance(depth+1, di, user, store, dp.Name))
	}
	return values
}

func typeToObjectTyperName(typer reflect.Type) TyperDependency {
	return TyperDependency{
		Typer: typer,
		Name:  gwreflect.GetPkgFullName(typer),
		IsPtr: typer.Kind() == reflect.Ptr,
	}
}
