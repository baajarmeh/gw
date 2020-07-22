package db

import "github.com/oceanho/gw/contrib/app"

func Migrate(backend app.Backend) {
	dbStore := backend.GetDbStore()
	dbStore.AutoMigrate(&Tenant{})
}
