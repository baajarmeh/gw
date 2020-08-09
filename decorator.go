package gw

import "sync"

type Decorator struct {
	Catalog  string
	MetaData interface{}
	Before   DecoratorHandler
	After    DecoratorHandler
}

type DecoratorList struct {
	ds     []Decorator
	locker sync.Mutex
}

func NewDecorators(decorators ...Decorator) *DecoratorList {
	dList := &DecoratorList{}
	dList.Append(decorators...)
	return dList
}

func (dc *DecoratorList) All() []Decorator {
	dc.locker.Lock()
	dc.locker.Unlock()
	var dcs = make([]Decorator, len(dc.ds))
	copy(dcs, dc.ds)
	return dcs
}

func (dc *DecoratorList) Append(decorators ...Decorator) *DecoratorList {
	dc.locker.Lock()
	dc.locker.Unlock()
	if len(decorators) < 1 {
		return dc
	}
	dc.ds = append(dc.ds, decorators...)
	return dc
}

type DecoratorHandler func(ctx *Context) (friendlyMsg string, err error)

// helpers
func FilterDecorator(filter func(d Decorator) bool, decorators ...Decorator) []Decorator {
	var result []Decorator
	for _, dc := range decorators {
		if filter(dc) {
			result = append(result, dc)
		}
	}
	return result
}
