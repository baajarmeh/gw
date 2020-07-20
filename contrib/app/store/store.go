package store

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw/contrib/app/auth"
	"gorm.io/gorm"
)

type Backend interface {
	GetDbStore() *gorm.DB
	GetDbStoreByName(name string) *gorm.DB
	GetCacheStore() *redis.Client
	GetCacheStoreByName(name string) *redis.Client
}

func GetBackend(c *gin.Context, user auth.User) Backend {
	return nil
}
