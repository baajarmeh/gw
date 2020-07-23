package gw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	dbMySQL "github.com/go-sql-driver/mysql"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

// Store represents a Store engine of gw framework.
type Store interface {
	GetDbStore() *gorm.DB
	GetDbStoreByName(name string) *gorm.DB
	GetCacheStore() *redis.Client
	GetCacheStoreByName(name string) *redis.Client
}

// StoreDbHandler defines a database ORM object handler.
type StoreDbHandler func(ctx *gin.Context, user User) *gorm.DB

type internalBackendWrapper struct {
	ctx   *gin.Context
	user  User
	store Store
}

func (b internalBackendWrapper) GetDbStore() *gorm.DB {
	db := b.store.GetDbStore().WithContext(context.Background())
	if db == nil {
		panic("got db store fail, ret is nil.")
	}
	return b.addGlobalDbFilter(db)
}

func (b internalBackendWrapper) GetDbStoreByName(name string) *gorm.DB {
	db := b.store.GetDbStoreByName(name).WithContext(context.Background())
	if db == nil {
		panic("got db store by name fail, ret is nil.")
	}
	return b.addGlobalDbFilter(db)
}

func (b internalBackendWrapper) addGlobalDbFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (b internalBackendWrapper) addGlobalCacheSettings(db *redis.Client) *redis.Client {
	return db
}

func (b internalBackendWrapper) GetCacheStore() *redis.Client {
	db := b.store.GetCacheStore()
	if db == nil {
		panic("got cache store fail, ret is nil.")
	}
	return b.addGlobalCacheSettings(db)
}

func (b internalBackendWrapper) GetCacheStoreByName(name string) *redis.Client {
	db := b.store.GetCacheStoreByName(name)
	if db == nil {
		panic("got cache store by name fail, ret is nil.")
	}
	return b.addGlobalCacheSettings(db)
}

type DefaultBackendImpl struct {
	dbs    map[string]*gorm.DB
	caches map[string]*redis.Client
}

func (d DefaultBackendImpl) GetDbStore() *gorm.DB {
	return d.GetDbStoreByName("primary")
}

func (d DefaultBackendImpl) GetDbStoreByName(name string) *gorm.DB {
	db, ok := d.dbs[name]
	if !ok {
		logger.Warn("got db: %s fail. not found.", name)
	}
	return db
}

func (d DefaultBackendImpl) GetCacheStore() *redis.Client {
	return d.GetCacheStoreByName("primary")
}

func (d DefaultBackendImpl) GetCacheStoreByName(name string) *redis.Client {
	db, ok := d.caches[name]
	if !ok {
		logger.Warn("got cache: %s fail. not found.", name)
	}
	return db
}

var backend Store
var once sync.Once

func Initial(conf conf.Config, initial func(conf conf.Config) Store) {
	once.Do(func() {
		backend = initial(conf)
	})
}

func DefaultBackend(cnf conf.Config) Store {
	storeBackend := DefaultBackendImpl{
		dbs:    make(map[string]*gorm.DB),
		caches: make(map[string]*redis.Client),
	}
	dbs := cnf.Common.Backend.Db
	caches := cnf.Common.Backend.Cache
	for _, v := range dbs {
		storeBackend.dbs[v.Name] = createDb(v)
	}
	for _, v := range caches {
		storeBackend.caches[v.Name] = createCache(v)
	}
	return storeBackend
}

func createDb(db conf.Db) *gorm.DB {
	if db.Driver != "mysql" {
		panic("not supports db.Driver: %s, only supports mysql.")
	}

	dbConf := dbMySQL.NewConfig()
	dbConf.Addr = fmt.Sprintf("%s:%d", db.Addr, db.Port)
	dbConf.DBName = db.Database
	dbConf.User = db.User
	dbConf.Passwd = db.Password
	pt, _ := strconv.ParseBool(db.Args["parseTime"])
	dbConf.ParseTime = pt

	gDialect := mysql.Open(dbConf.FormatDSN())
	gDbConf := &gorm.Config{}
	gDb, err := gorm.Open(gDialect, gDbConf)
	if err != nil {
		panic("create backend store db fail.")
	}
	sqlDb, err := gDb.DB()
	if err != nil {
		panic("get sql.DB fail.")
	}
	if err := sqlDb.Ping(); err != nil {
		panic(fmt.Sprintf("db not pong, db name:%s. addr: %s, port: %d", db.Name, db.Addr, db.Port))
	}
	return gDb
}

func createCache(db conf.Cache) *redis.Client {
	return nil
}

func GetBackend(ctx *gin.Context, user User) Store {
	bk := &internalBackendWrapper{
		ctx:   ctx,
		user:  user,
		store: backend,
	}
	return bk
}
