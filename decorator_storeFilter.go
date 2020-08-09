package gw

import (
	"sync"
)

const storeDbFilterDecoratorCatalog = "gw_framework_store_filter"

type StoreDbFilterDecorator struct {
	locker sync.Mutex
	items  []Decorator
}

func NewStoreDbFilterDecorator(setupDbFilterHandler StoreDbSetupHandler) []Decorator {
	var decorators []Decorator
	d := Decorator{
		Catalog:  storeDbFilterDecoratorCatalog,
		MetaData: nil,
		Before: func(ctx *Context) (friendlyMsg string, err error) {
			store, ok := ctx.Store.(*backendWrapper)
			if ok {
				store.storeDbSetupHandler = setupDbFilterHandler
				ctx.Store = store
			}
			return "", nil
		},
	}
	decorators = append(decorators, d)
	return decorators
}
