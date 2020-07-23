package db

import (
	"github.com/oceanho/gw/conf"
)

func Migrate(backend conf.Backend) {
	dbStore := backend.GetDbStore()
	dbStore.AutoMigrate(&Tenant{})
}
