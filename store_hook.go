package gw

import (
	"gorm.io/gorm"
	"reflect"
	"sync"
)

type DbOpHandler func(db *gorm.DB, ctx *Context, model interface{}) error

type DbOpTyperHandlers struct {
	locker   sync.Mutex
	handlers map[reflect.Type][]DbOpHandler
}

type DbOpProcessor struct {
	fns *map[string]*DbOpTyperHandlers
}

func NewDbOpProcessor() *DbOpProcessor {
	var maps = make(map[string]*DbOpTyperHandlers)
	maps["gw:on_save_before"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	maps["gw:on_save_after"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	maps["gw:on_query_before"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	maps["gw:on_query_after"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	maps["gw:on_delete_before"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	maps["gw:on_delete_after"] = &DbOpTyperHandlers{
		handlers: make(map[reflect.Type][]DbOpHandler),
	}
	return &DbOpProcessor{
		fns: &maps,
	}
}

func (processor *DbOpProcessor) SaveBefore() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_save_before"]
}

func (processor *DbOpProcessor) SaveAfter() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_save_after"]
}

func (processor *DbOpProcessor) QueryBefore() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_query_before"]
}

func (processor *DbOpProcessor) QueryAfter() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_query_after"]
}

func (processor *DbOpProcessor) DeleteBefore() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_delete_before"]
}

func (processor *DbOpProcessor) DeleteAfter() *DbOpTyperHandlers {
	return (*(processor.fns))["gw:on_delete_after"]
}

func (handlers *DbOpTyperHandlers) Register(handler DbOpHandler, models ...interface{}) *DbOpTyperHandlers {
	handlers.locker.Lock()
	defer handlers.locker.Unlock()
	for _, m := range models {
		_model := m
		typer := reflect.TypeOf(_model)
		if handlers.handlers[typer] == nil {
			if len(handlers.handlers[typer]) == 0 {
				handlers.handlers[typer] = make([]DbOpHandler, 0, 8)
			}
			handlers.handlers[typer] = append(handlers.handlers[typer], handler)
		}
	}
	return handlers
}
