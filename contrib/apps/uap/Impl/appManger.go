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
	cachePrefix      string
	cacheStoreName   string
	backendStoreName string
	cacheExpiration  time.Duration
	permPagerExpr    gw.PagerExpr
}

func (a AppManager) CacheKey(name string) string {
	return fmt.Sprintf("%s.%s", a.cachePrefix, name)
}

func (a AppManager) Password(secret string) string {
	return a.PasswordSigner().Sign(secret)
}

func (a AppManager) Backend() *gorm.DB {
	return a.Store().GetDbStoreByName(a.backendStoreName)
}

func (a AppManager) Cache() *redis.Client {
	return a.Store().GetCacheStoreByName(a.cacheStoreName)
}

func (a AppManager) GetFromCache(passport string) gw.User {
	var bytes, err = a.Cache().Get(context.Background(), a.CacheKey(passport)).Bytes()
	if err == nil {
		var user gw.User
		if err := json.Unmarshal(bytes, &user); err == nil {
			return user
		}
	}
	return gw.EmptyUser
}

func (a AppManager) RemoveFromCache(passport string) {
	_ = a.Cache().Del(context.Background(), a.CacheKey(passport))
}

func (a AppManager) SaveToCache(user gw.User) {
	if err := a.Cache().Set(context.Background(),
		a.CacheKey(user.Passport), user, a.cacheExpiration).Err(); err != nil {
		logger.Error("SaveAuthUserToCache fail, err: %v", err)
	}
}

//
// APIs impl
//
func (a AppManager) Create(app gw.AppInfo) error {
	return a.Backend().FirstOrCreate(&Db.App{
		Name:       app.Name,
		Router:     app.Router,
		Descriptor: app.Descriptor,
	}, "name = ?", app.Name).Error
}

func (a AppManager) QueryByName(name string) *gw.AppInfo {
	var model Db.App
	if err := a.Backend().First(&model, "name = ?", name).Error; err != nil {
		return nil
	}
	var app gw.AppInfo
	app.ID = model.ID
	app.Name = model.Name
	app.Router = model.Router
	app.Descriptor = model.Descriptor
	return &app
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
