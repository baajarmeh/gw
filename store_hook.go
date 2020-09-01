package gw

import (
	"gorm.io/gorm"
	"reflect"
)

type DbOpHandler func(db *gorm.DB, ctx *Context, model interface{}) error
type DbOpTyperHandlers map[reflect.Type][]DbOpHandler
type DbOpActionMap map[string]*DbOpTyperHandlers
type DbOpProcessor struct {
	fns *DbOpActionMap
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

func (processor *DbOpProcessor) Register(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}

func (processor *DbOpProcessor) OnSaveAfter(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}

func (processor *DbOpProcessor) OnQueryBefore(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}

func (processor *DbOpProcessor) OnQueryAfter(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}

func (processor *DbOpProcessor) OnDeleteBefore(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}

func (processor *DbOpProcessor) OnDeleteAfter(handler DbOpHandler, models ...interface{}) *DbOpProcessor {
	return processor
}
