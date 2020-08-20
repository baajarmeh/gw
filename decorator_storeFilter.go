package gw

import (
	"sync"
)

const storeDbFilterDecoratorCatalog = "gw_framework_store_filter"

type StoreDbFilterDecorator struct {
	locker sync.Mutex
	items  []Decorator
}

func NewStoreDbSetupDecorator(setupDbFilterHandler StoreDbSetupHandler) Decorator {
	var d = Decorator{
		Catalog:  storeDbFilterDecoratorCatalog,
		MetaData: nil,
		Before: func(ctx *Context) (status int, err error, payload interface{}) {
			store, ok := ctx.State.Store().(*backendWrapper)
			if ok {
				store.storeDbSetupHandler = setupDbFilterHandler
				ctx.State.s.Store = store
			}
			return 0, nil, nil
		},
	}
	return d
}
