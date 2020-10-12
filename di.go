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
	NullReflectValue             = reflect.ValueOf(nil)
	BuiltinComponentTyper        = reflect.TypeOf(BuiltinComponent{})
	ErrorObjectTyperIsNonPointer = fmt.Errorf("object should be a pointer typer")
)

func (ss *ServerState) objectTypers() map[string]ObjectTyper {
	var typers = make(map[string]ObjectTyper)
	typers[UserTyperName] = newNilApiObjectTyper(UserTyperName, EmptyUser)
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
	AppManager          IAppManager
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
	AppManager IAppManager,
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
	bc.AppManager = AppManager
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
	ResolveFunc func(di interface{}, state interface{}, typerName string) (error, interface{})
}

var defaultDIConfig = DIConfig{
	NewFuncName: "New",
	ResolveFunc: func(diImpl interface{}, state interface{}, typerName string) (error, interface{}) {
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
		e, v := resolverTyperInstance(di, user, store, typerName)
		return e, v.Interface()
	},
}

type IDIProvider interface {
	Register(actual ...interface{}) bool
	RegisterWithTyper(typers ...reflect.Type) bool
	RegisterWithName(typerName string, actual interface{}) bool
	Resolve(typerName string) (error, interface{})
	ResolveByTyper(typer reflect.Type) (error, interface{})
	ResolveByObjectTyper(object interface{}) error
	ResolveWithState(state interface{}, typerName string) (error, interface{})
	ResolveByTyperWithState(state interface{}, typer reflect.Type) (error, interface{})
	Check() bool
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
		panic(fmt.Sprintf("typer(%s) has no name of %s(...) method",
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

func (d *DefaultDIProviderImpl) Resolve(typerName string) (error, interface{}) {
	return d.ResolveWithState(d.state.Store(), typerName)
}

func (d *DefaultDIProviderImpl) ResolveByTyper(typer reflect.Type) (error, interface{}) {
	return d.ResolveByTyperWithState(d.state, typer)
}

func (d *DefaultDIProviderImpl) ResolveByObjectTyper(object interface{}) error {
	typer := reflect.TypeOf(object)
	if typer.Kind() != reflect.Ptr {
		return ErrorObjectTyperIsNonPointer
	}
	e, v := d.ResolveByTyper(typer)
	if e != nil {
		return e
	}
	reflect.ValueOf(object).Elem().Set(reflect.ValueOf(v))
	return nil
}

func (d *DefaultDIProviderImpl) ResolveByTyperWithState(state interface{}, typer reflect.Type) (error, interface{}) {
	var typerName, ok = d.typerMappers[typer]
	if !ok {
		typerName = gwreflect.GetPkgFullName(typer)
	}
	return d.ResolveWithState(state, typerName)
}

func (d *DefaultDIProviderImpl) ResolveWithState(state interface{}, typerName string) (error, interface{}) {
	return d.config.ResolveFunc(d, state, typerName)
}

// Check all of register Object has cycle references.
func (d *DefaultDIProviderImpl) Check() bool {
	return false
}

// helpers
func resolverTyperInstance(di *DefaultDIProviderImpl, user User, store IStore, typerName string) (error, reflect.Value) {
	switch typerName {
	case IStoreName:
		if store != nil {
			return nil, reflect.ValueOf(store)
		}
	case UserTyperName:
		return nil, reflect.ValueOf(user)
	}
	var objectTyper, ok = di.objectTypers[typerName]
	if !ok {
		return fmt.Errorf("missing typer(%s)", typerName), NullReflectValue
	}
	if objectTyper.newAPI == NullReflectValue {
		return nil, objectTyper.ActualValue
	}
	e, values := resolveTyperDependOn(di, &objectTyper, user, store)
	if e != nil {
		return e, NullReflectValue
	}
	var result = objectTyper.newAPI.Call(values)[0]
	switch result.Kind() {
	case reflect.Ptr:
		if objectTyper.IsPtr {
			return nil, result
		} else {
			return nil, reflect.ValueOf(result.Elem().Interface())
		}
	default:
		if !objectTyper.IsPtr {
			return nil, result
		} else {
			var typeValue = reflect.New(reflect.TypeOf(result.Interface()))
			typeValue.Elem().Set(result)
			return nil, typeValue
		}
	}
}

func resolveTyperDependOn(di *DefaultDIProviderImpl, objTyper *ObjectTyper, user User, store IStore) (error, []reflect.Value) {
	var values = make([]reflect.Value, 0, 8)
	for _, dp := range objTyper.DependOn {
		e, value := resolverTyperInstance(di, user, store, dp.Name)
		if e != nil {
			return e, nil
		}
		if dp.IsPtr && value.Kind() != reflect.Ptr {
			var typeValue = reflect.New(reflect.TypeOf(value.Interface()))
			typeValue.Elem().Set(value)
			value = typeValue
		}
		// FIXME(Ocean): has some problem when register object are not receiver but references by pointer typer
		//else if !dp.IsPtr && value.Kind() != reflect.Interface && value.Kind() == reflect.Ptr {
		//	value = reflect.ValueOf(value.Elem().Interface())
		//}
		values = append(values, value)
	}
	return nil, values
}

func typeToObjectTyperName(typer reflect.Type) TyperDependency {
	return TyperDependency{
		Typer: typer,
		Name:  gwreflect.GetPkgFullName(typer),
		IsPtr: typer.Kind() == reflect.Ptr,
	}
}
