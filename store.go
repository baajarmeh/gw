package gw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/go-redis/redis/v8"
	mysqlDb "github.com/go-sql-driver/mysql"
	"github.com/oceanho/gw/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// IStore represents a Store engine of gw framework.
type IStore interface {
	GetDbStore() *gorm.DB
	GetDbStoreByName(name string) *gorm.DB
	GetCacheStore() *redis.Client
	GetCacheStoreByName(name string) *redis.Client
}

// StoreDbSetupHandler represents a database ORM object handler that can be replace A *gorm.DB instances features.
type StoreDbSetupHandler func(ctx *Context, db *gorm.DB) *gorm.DB

// StoreCacheSetupHandler represents a redis Client object handler that can be replace A *redis.Client instances features.
type StoreCacheSetupHandler func(ctx *Context, client *redis.Client, user User) *redis.Client

// SessionStateHandler represents a Session state manager handler.
type SessionStateHandler func(state *ServerState) ISessionStateManager

// AuthManagerHandler represents a Permission state manager handler.
type AuthManagerHandler func(state *ServerState) IAuthManager

// PermissionManagerHandler represents a Permission state manager handler.
type PermissionManagerHandler func(state *ServerState) IPermissionManager

// RespBodyBuildFunc represents a response body build handler.
type RespBodyBuildFunc func(ctx *gin.Context, status int, requestID string, err interface{}, msgBody interface{}) interface{}

type backendWrapper struct {
	user                    User
	store                   IStore
	ctx                     *Context
	storeDbSetupHandlers    []StoreDbSetupHandler
	storeCacheSetupHandlers []StoreCacheSetupHandler
}

func (b *backendWrapper) GetDbStore() *gorm.DB {
	db := b.store.GetDbStore()
	if db == nil {
		panic("got db store fail, ret is nil.")
	}
	return b.globalDbStep(db)
}

func (b *backendWrapper) GetDbStoreByName(name string) *gorm.DB {
	db := b.store.GetDbStoreByName(name)
	if db == nil {
		panic("got db store by name fail, ret is nil.")
	}
	return b.globalDbStep(db)
}

func (b *backendWrapper) globalDbStep(db *gorm.DB) *gorm.DB {
	for _, h := range b.storeDbSetupHandlers {
		db = h(b.ctx, db)
	}
	return db
}

func (b *backendWrapper) globalCacheSetup(db *redis.Client) *redis.Client {
	for _, h := range b.storeCacheSetupHandlers {
		db = h(b.ctx, db, b.user)
	}
	return db
}

func (b *backendWrapper) GetCacheStore() *redis.Client {
	db := b.store.GetCacheStore()
	if db == nil {
		panic("got cache store fail, ret is nil.")
	}
	return b.globalCacheSetup(db)
}

func (b *backendWrapper) GetCacheStoreByName(name string) *redis.Client {
	db := b.store.GetCacheStoreByName(name)
	if db == nil {
		panic("got cache store by name fail, ret is nil.")
	}
	return b.globalCacheSetup(db)
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
		//logger.Warn("got db: %s fail. not found.", name)
		panic(fmt.Sprintf("got db: %s fail. not found.", name))
	}
	return db
}

func (d DefaultBackendImpl) GetCacheStore() *redis.Client {
	return d.GetCacheStoreByName("primary")
}

func (d DefaultBackendImpl) GetCacheStoreByName(name string) *redis.Client {
	db, ok := d.caches[name]
	if !ok {
		//logger.Warn("got cache: %s fail. not found.", name)
		panic(fmt.Sprintf("got cache: %s fail. not found.", name))
	}
	return db
}

func GetContextFromDB(db *gorm.DB) (*Context, bool) {
	var obj, ok = db.Get(gwDbContextKey)
	if !ok {
		return nil, false
	}
	ctx, ok := obj.(*Context)
	return ctx, ok
}

func DefaultBackend(cnf *conf.ApplicationConfig) IStore {
	storeBackend := DefaultBackendImpl{
		dbs:    make(map[string]*gorm.DB),
		caches: make(map[string]*redis.Client),
	}
	dbs := cnf.Backend.Db
	caches := cnf.Backend.Cache
	for _, v := range dbs {
		storeBackend.dbs[v.Name] = createDb(v)
	}
	for _, v := range caches {
		db := createCache(v)
		storeBackend.caches[v.Name] = db
	}
	return storeBackend
}

func createDb(db conf.Db) *gorm.DB {
	if db.Driver != "mysql" {
		panic("not supports db.Driver: %s, only supports mysql.")
	}

	params := make(map[string]string)
	for k, v := range db.Args {
		params[k] = fmt.Sprintf("%v", v)
	}

	dbConf := mysqlDb.NewConfig()
	dbConf.DBName = db.Database
	dbConf.User = db.User
	dbConf.Passwd = db.Password
	parseTime, ok := params["parseTime"]
	if ok {
		dbConf.ParseTime, _ = strconv.ParseBool(parseTime)
	}
	dbConf.Addr = fmt.Sprintf("%s:%d", db.Addr, db.Port)
	dbConf.Params = params

	gDialect := mysql.Open(dbConf.FormatDSN())
	//
	// https://gorm.io/docs/performance.html
	gDbConf := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}
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
	setupDb(gDb)
	return gDb
}

func setupDb(db *gorm.DB) {
	err := db.Callback().Create().Before("gorm:create").Register("gw:create_before", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			fns, ok := ctx.server.DbOpProcessor.CreateBefore().handlers[db.Statement.Schema.ModelType]
			if ok {
				for _, f := range fns {
					err := f(db, ctx)
					if err != nil {
						_ = db.AddError(err)
					}
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Create().After("gorm:create").Register("gw:create_after", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			fns, ok := ctx.server.DbOpProcessor.CreateAfter().handlers[db.Statement.Schema.ModelType]
			if ok {
				for _, f := range fns {
					err = f(db, ctx)
					if err != nil {
						_ = db.AddError(err)
					}
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Update().Before("gorm:update").Register("gw:update_before", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			fns, ok := ctx.server.DbOpProcessor.UpdateBefore().handlers[db.Statement.Schema.ModelType]
			if ok {
				for _, f := range fns {
					err = f(db, ctx)
					if err != nil {
						_ = db.AddError(err)
					}
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Update().After("gorm:update").Register("gw:update_after", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			handlers := ctx.server.DbOpProcessor.UpdateAfter().Handlers(db.Statement.Schema.ModelType)
			for _, f := range handlers {
				err := f(db, ctx)
				if err != nil {
					_ = db.AddError(err)
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Delete().Before("gorm:delete").Register("gw:update_before", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			if ok {
				handlers := ctx.server.DbOpProcessor.DeleteBefore().Handlers(db.Statement.Schema.ModelType)
				for _, f := range handlers {
					err := f(db, ctx)
					if err != nil {
						_ = db.AddError(err)
					}
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Delete().After("gorm:delete").Register("gw:delete_after", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			handlers := ctx.server.DbOpProcessor.DeleteAfter().Handlers(db.Statement.Schema.ModelType)
			for _, f := range handlers {
				err := f(db, ctx)
				if err != nil {
					_ = db.AddError(err)
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Query().Before("gorm:query").Register("gw:query_before", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			handlers := ctx.server.DbOpProcessor.QueryBefore().Handlers(db.Statement.Schema.ModelType)
			for _, f := range handlers {
				err := f(db, ctx)
				if err != nil {
					_ = db.AddError(err)
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
	err = db.Callback().Query().After("gorm:query").Register("gw:query_after", func(db *gorm.DB) {
		var obj, ok = db.Get(gwDbContextKey)
		if !ok {
			return
		}
		if ctx, ok := obj.(*Context); ok {
			fns, ok := ctx.server.DbOpProcessor.QueryAfter().handlers[db.Statement.Schema.ModelType]
			if ok {
				for _, f := range fns {
					err = f(db, ctx)
					if err != nil {
						_ = db.AddError(err)
					}
				}
			}
		}
	})
	if err != nil {
		panic(fmt.Sprintf("setup db hooks fail, err: %v", err))
	}
}

func createCache(cache conf.Cache) *redis.Client {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cache.Addr, cache.Port),
		Password: cache.Password,
		DB:       cache.DB,
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("redis not pong, name:%s. addr: %s, port: %d", cache.Name, cache.Addr, cache.Port))
	}
	return client
}
