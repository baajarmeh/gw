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

type TyperDependency struct {
	Name  string
	IsPtr bool
}

type IDIProvider interface {
	Register(actual ...interface{}) bool
	RegisterWithName(typerName string, actual interface{}) bool
	Resolve(typerName string) interface{}
	ResolveByTyper(typer reflect.Type) interface{}
	ResolveWithState(state ContextState, typerName string) interface{}
	ResolveByTyperWithState(state ContextState, typer reflect.Type) interface{}
}

type DefaultDIProviderImpl struct {
	locker       sync.Mutex
	objectTypers []ObjectTyper
	typedMaps    map[string]ObjectTyper
	state        ServerState
}

func DefaultDIProvider(state ServerState) IDIProvider {
	return &DefaultDIProviderImpl{
		state:        state,
		objectTypers: make([]ObjectTyper, 0),
		typedMaps:    make(map[string]ObjectTyper),
	}
}

func (d *DefaultDIProviderImpl) Register(actual ...interface{}) bool {
	for _, a := range actual {
		if !d.RegisterWithName("", a) {
			return false
		}
	}
	return true
}

func (d *DefaultDIProviderImpl) RegisterWithName(name string, actual interface{}) bool {
	var actualValue = reflect.ValueOf(actual)
	var newMethod = actualValue.MethodByName("New")
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
	var typeHasRegistered = false
	if _, ok := d.typedMaps[virtualName]; ok {
		for i := 0; i < len(d.objectTypers); i++ {
			if d.objectTypers[i].Name == virtualName {
				typeHasRegistered = true
				d.objectTypers[i] = objectTyper
				break
			}
		}
	}
	d.typedMaps[virtualName] = objectTyper
	if !typeHasRegistered {
		d.objectTypers = append(d.objectTypers, objectTyper)
	}
	return true
}

func (d *DefaultDIProviderImpl) Resolve(typerName string) interface{} {
	ctxState := ContextState{
		state: d.state,
		store: d.state.Store(),
	}
	return d.ResolveWithState(ctxState, typerName)
}

func (d *DefaultDIProviderImpl) ResolveByTyper(typer reflect.Type) interface{} {
	ctxState := ContextState{
		state: d.state,
		store: d.state.Store(),
	}
	return d.ResolveByTyperWithState(ctxState, typer)
}

func (d *DefaultDIProviderImpl) ResolveByTyperWithState(state ContextState, typer reflect.Type) interface{} {
	return d.ResolveWithState(state, gwreflect.GetPkgFullName(typer))
}

func (d *DefaultDIProviderImpl) ResolveWithState(state ContextState, typerName string) interface{} {
	var objectTyper, ok = d.typedMaps[typerName]
	if !ok {
		panic(fmt.Sprintf("object typer(%s) not found", typerName))
	}
	var result reflect.Value
	if objectTyper.newAPI == nullReflectValue {
		result = objectTyper.ActualValue
	} else {
		var values []reflect.Value
		var preparedTypers = make(map[string]ObjectTyper)
		for _, objTyper := range d.typedMaps {
			_objTyper := objTyper
			preparedTypers[objTyper.Name] = _objTyper
		}
		prepareBuiltinTypers(preparedTypers, state)
		for _, typerDp := range objectTyper.DependOn {
			values = append(values, d.resolver(typerDp, preparedTypers))
		}
		result = objectTyper.newAPI.Call(values)[0]
	}
	return result.Interface()
}

// helpers
func (d *DefaultDIProviderImpl) resolver(typerDependency TyperDependency, preparedTypers map[string]ObjectTyper) reflect.Value {
	var values []reflect.Value
	var objectTyper, ok = preparedTypers[typerDependency.Name]
	if !ok {
		panic(fmt.Sprintf("missing typer(%s)", typerDependency.Name))
	}
	if objectTyper.newAPI == nullReflectValue {
		return objectTyper.ActualValue
	}
	for _, dp := range objectTyper.DependOn {
		if _, ok := preparedTypers[dp.Name]; ok {
			values = append(values, d.resolver(dp, preparedTypers))
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
func prepareBuiltinTypers(typers map[string]ObjectTyper, state ContextState) {
	for k, v := range state.objectTypers() {
		typers[k] = v
	}
}

func typeToObjectTyperName(typ reflect.Type) TyperDependency {
	return TyperDependency{
		Name:  gwreflect.GetPkgFullName(typ),
		IsPtr: typ.Kind() == reflect.Ptr,
	}
}

func typeToObjectTyper(typ reflect.Type) ObjectTyper {
	var typer ObjectTyper
	typer.Typer = typ
	typer.Name = gwreflect.GetPkgFullName(typ)
	return typer
}
