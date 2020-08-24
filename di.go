package gw

import (
	"fmt"
	"github.com/oceanho/gw/libs/gwreflect"
	"github.com/oceanho/gw/logger"
	"reflect"
	"sync"
)

type ObjectTyper struct {
	Name        string
	DependOn    []string
	Typer       reflect.Type
	ActualValue reflect.Value
	newAPI      reflect.Value
}

type BuildFunc func(state ContextState, typers map[string]ObjectTyper) (interface{}, error)

type IDIProvider interface {
	Register(actual interface{}) bool
	RegisterWithName(name string, actual interface{}) bool
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

type parameter struct {
	Name  string
	Value reflect.Value
}

func DefaultDIProvider(state ServerState) IDIProvider {
	return &DefaultDIProviderImpl{
		state:        state,
		objectTypers: make([]ObjectTyper, 0),
		typedMaps:    make(map[string]ObjectTyper),
	}
}

func (d *DefaultDIProviderImpl) Register(actual interface{}) bool {
	return d.RegisterByName("", actual)
}

func (d *DefaultDIProviderImpl) RegisterByName(name string, actual interface{}) bool {
	return d.RegisterWithName(name, actual)
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
	var newFunParamObjectTypers = make([]string, newMethodNumIn)
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
		logger.Error("object typer(%s) not found", typerName)
		return nil
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
		for _, dp := range objectTyper.DependOn {
			values = append(values, d.resolver(dp, preparedTypers))
		}
		result = objectTyper.newAPI.Call(values)[0]
	}
	if result.Kind() == reflect.Ptr {
		//FIXME(OceanHo): how to returns a pointer typer?
		// looking for: #OC-20200824.001
		return result.Elem().Interface()
	}
	return result.Interface()
}

// helpers
func (d *DefaultDIProviderImpl) resolver(typerName string, preparedTypers map[string]ObjectTyper) reflect.Value {
	var values []reflect.Value
	var objectTyper, ok = preparedTypers[typerName]
	if !ok {
		panic(fmt.Sprintf("missing typer(%s)", typerName))
	}
	if objectTyper.newAPI == nullReflectValue {
		return objectTyper.ActualValue
	}
	for _, dp := range objectTyper.DependOn {
		if objTyper, ok := preparedTypers[dp]; ok {
			values = append(values, d.resolver(objTyper.Name, preparedTypers))
		}
	}
	return objectTyper.newAPI.Call(values)[0]
}
func prepareBuiltinTypers(typers map[string]ObjectTyper, state ContextState) {
	for k, v := range state.objectTypers() {
		typers[k] = v
	}
}

func typeToObjectTyperName(typ reflect.Type) string {
	return gwreflect.GetPkgFullName(typ)
}

func typeToObjectTyper(typ reflect.Type) ObjectTyper {
	var typer ObjectTyper
	typer.Typer = typ
	typer.Name = gwreflect.GetPkgFullName(typ)
	return typer
}
