package store

import (
	"github.com/gin-gonic/gin"
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

func GetBackend(c *gin.Context) Backend {
	return nil
}
