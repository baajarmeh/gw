package store

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Backend interface {
	GetDbStore() *gorm.DB
	GetDbStoreByName(name string) *gorm.DB
	GetCacheStore() *redis.Client
	GetCacheStoreByName(name string) *redis.Client
}
