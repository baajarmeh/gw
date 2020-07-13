package db

import "github.com/oceanho/gw/contrib/app/store"

func Migrate(backend store.Backend) {
	dbStore := backend.GetDbStore()
	dbStore.AutoMigrate(&Tenant{})
}
