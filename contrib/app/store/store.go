package store

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw/contrib/app/auth"
	"github.com/oceanho/gw/contrib/app/conf"
	"gorm.io/gorm"
	"sync"
)

type Backend interface {
	GetDbStore() *gorm.DB
	GetDbStoreByName(name string) *gorm.DB
	GetCacheStore() *redis.Client
	GetCacheStoreByName(name string) *redis.Client
}

type DefaultBackendImpl struct {
}

func (d DefaultBackendImpl) GetDbStore() *gorm.DB {
	panic("implement me")
}

func (d DefaultBackendImpl) GetDbStoreByName(name string) *gorm.DB {
	panic("implement me")
}

func (d DefaultBackendImpl) GetCacheStore() *redis.Client {
	panic("implement me")
}

func (d DefaultBackendImpl) GetCacheStoreByName(name string) *redis.Client {
	panic("implement me")
}

var backend Backend
var once sync.Once

func Initial(conf conf.Config, initial func(conf conf.Config) Backend) {
	once.Do(func() {
		backend = initial(conf)
	})
}

func Default(cnf conf.Config) Backend {
	return &DefaultBackendImpl{}
}

// Check represents a API that for check the store package requires has been initialed
// Throw a exception if your has been not call the Initial(...) API.
func Check() {
	if backend == nil {
		panic("store.backend is null. call store.Initial(...) please.")
	}
}

func GetBackend(c *gin.Context, user auth.User) Backend {
	return backend
}
