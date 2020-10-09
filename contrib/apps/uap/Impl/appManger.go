package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type AppManager struct {
	*gw.ServerState
	store            gw.IStore
	cachePrefix      string
	cacheStoreName   string
	backendStoreName string
	cacheExpiration  time.Duration
	permPagerExpr    gw.PagerExpr
}

// DI

func (a AppManager) New(store gw.IStore) gw.IAppManager {
	a.store = store
	return a
}

func (a AppManager) Store() gw.IStore {
	if a.store == nil {
		return a.ServerState.Store()
	}
	return a.store
}

func (a AppManager) cacheKey(name string) string {
	return fmt.Sprintf("%s.%s", a.cachePrefix, name)
}

func (a AppManager) Password(secret string) string {
	return a.PasswordSigner().Sign(secret)
}

func (a AppManager) Backend() *gorm.DB {
	return a.Store().GetDbStoreByName(a.backendStoreName)
}

func (a AppManager) cache() *redis.Client {
	return a.Store().GetCacheStoreByName(a.cacheStoreName)
}

func (a AppManager) GetFromCache(appName string) *gw.AppInfo {
	var bytes, err = a.cache().Get(context.Background(), a.cacheKey(appName)).Bytes()
	if err == nil {
		var appInfo gw.AppInfo
		if err := json.Unmarshal(bytes, &appInfo); err == nil {
			return &appInfo
		}
	}
	return nil
}

func (a AppManager) RemoveFromCache(name string) {
	_ = a.cache().Del(context.Background(), a.cacheKey(name))
}

func (a AppManager) SaveToCache(app gw.AppInfo) {
	if err := a.cache().Set(context.Background(),
		a.cacheKey(app.Name), app, a.cacheExpiration).Err(); err != nil {
		logger.Error("AppManager.SaveToCache fail, err: %v", err)
	}
}

//
// APIs impl
//
func (a AppManager) Create(app gw.AppInfo) error {
	// https://gorm.io/zh_CN/docs/advanced_query.html
	err := a.Backend().Where(Db.App{Key: app.Key}).FirstOrCreate(&Db.App{
		Key:        app.Key,
		Name:       app.Name,
		Router:     app.Router,
		Descriptor: app.Descriptor,
	}).Error
	if err != nil {
		return err
	}
	a.SaveToCache(app)
	return nil
}

func (a AppManager) QueryByName(name string) *gw.AppInfo {
	var app = a.GetFromCache(name)
	if app != nil {
		return app
	}
	var model Db.App
	if err := a.Backend().First(&model, "name = ?", name).Error; err != nil {
		return nil
	}
	app = &gw.AppInfo{}
	app.ID = model.ID
	app.Name = model.Name
	app.Router = model.Router
	app.Descriptor = model.Descriptor
	return app
}

func DefaultAppManager(state *gw.ServerState) gw.IAppManager {
	var cnf = Config.GetUAP(state.ApplicationConfig())
	return AppManager{
		ServerState:      state,
		backendStoreName: cnf.Backend.Name,
		cacheStoreName:   cnf.AppManager.Cache.Name,
		cachePrefix:      cnf.AppManager.Cache.Prefix,
		cacheExpiration:  time.Hour * time.Duration(cnf.AppManager.Cache.ExpirationHours), // One week.
	}
}
