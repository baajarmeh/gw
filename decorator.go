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

type DecoratorHandler func(ctx *Context) (status int, err error, payload interface{})

type DecoratorHandlerResult struct {
	status  int
	error   error
	payload interface{}
}

func NewDecoratorHandlerResult(status int, err error, payload interface{}) DecoratorHandlerResult {
	return DecoratorHandlerResult{
		status:  status,
		error:   err,
		payload: payload,
	}
}

func (r DecoratorHandlerResult) Result() (status int, err error, payload interface{}) {
	return r.status, r.error, r.payload
}

var (
	DecoratorHandlerOK  = NewDecoratorHandlerResult(0, nil, nil)
	DecoratorHandler400 = NewDecoratorHandlerResult(400, ErrBadRequest, errDefault400Msg)
	DecoratorHandler401 = NewDecoratorHandlerResult(401, ErrUnauthorized, errDefault401Msg)
	DecoratorHandler403 = NewDecoratorHandlerResult(403, ErrPermissionDenied, errDefault403Msg)
	DecoratorHandler404 = NewDecoratorHandlerResult(404, ErrNotFoundRequest, errDefault404Msg)
	DecoratorHandler500 = NewDecoratorHandlerResult(500, ErrInternalServerError, errDefault500Msg)
)

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
