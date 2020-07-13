package store

import "gorm.io/gorm"

type CacheStore interface {
	Get(key interface{}) interface{}
	Set(key interface{}, value interface{})
}

type Backend interface {
	GetDbStore() *gorm.DB
	GetCacheStore(index int) CacheStore
}
