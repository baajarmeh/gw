package store

import (
	"github.com/oceanho/gw/contrib/app/auth"
	"gorm.io/gorm"
)

type CacheStore interface {
	Get(key interface{}) interface{}
	Set(key interface{}, value interface{})
}

type Backend interface {
	GetDbStore() *gorm.DB
	GetCacheStore(index int) CacheStore
}

func (d DefaultBackendImpl) GetDbStore() *gorm.DB {
	panic("implement me")
}

func (d DefaultBackendImpl) GetCacheStore(index int) CacheStore {
	panic("implement me")
}

type DefaultBackendImpl struct {
}

func GetBackend(user auth.User) Backend {
	return nil
}
