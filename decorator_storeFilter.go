package gw

import (
	"sync"
)

const storeDbFilterDecoratorCatalog = "gw_framework_store_filter"

type StoreDbFilterDecorator struct {
	locker sync.Mutex
	items  []Decorator
}

func NewStoreDbFilterDecorator(setupDbFilterHandler StoreDbSetupHandler) Decorator {
	var d = Decorator{
		Catalog:  storeDbFilterDecoratorCatalog,
		MetaData: nil,
		Before: func(ctx *Context) (status int, err error, payload interface{}) {
			store, ok := ctx.Store.(*backendWrapper)
			if ok {
				store.storeDbSetupHandler = setupDbFilterHandler
				ctx.Store = store
			}
			return 0, nil, nil
		},
	}
	return d
}
