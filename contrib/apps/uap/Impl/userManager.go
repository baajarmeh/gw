package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"github.com/oceanho/gw/libs/gwjsoner"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type UserManager struct {
	store            gw.IStore
	cachePrefix      string
	cacheStoreName   string
	backendStoreName string
	cacheExpiration  time.Duration
	permPagerExpr    gw.PagerExpr
	serverState      *gw.ServerState
}

// DI
func (u UserManager) New(store gw.IStore) gw.IUserManager {
	u.store = store
	return u
}

func (u UserManager) Store() gw.IStore {
	if u.store == nil {
		return u.serverState.Store()
	}
	return u.store
}

func (u UserManager) mapUserType(user Db.User) gw.UserType {
	if user.IsTenancy {
		return gw.TenancyUser
	}
	if user.IsAdmin {
		return gw.AdministrationUser
	}
	return gw.NormalUser
}

func (u UserManager) Backend() *gorm.DB {
	return u.Store().GetDbStoreByName(u.backendStoreName)
}

func (u UserManager) Cache() *redis.Client {
	return u.Store().GetCacheStoreByName(u.cacheStoreName)
}

func (u UserManager) GetFromCache(passport string) gw.User {
	var bytes, err = u.Cache().Get(context.Background(), u.CacheKey(passport)).Bytes()
	if err != nil {
		logger.Error("operation cache fail, err: %v", err)
		return gw.EmptyUser
	}
	var user gw.User
	if gwjsoner.Unmarshal(bytes, &user) != nil {
		return gw.EmptyUser
	}
	return user
}

func (u UserManager) CacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", u.cachePrefix, passport)
}

func (u UserManager) SaveToCache(user gw.User) {
}

//
// APIs
//
func (u UserManager) Create(user *gw.User) error {
	var model Db.User
	// default
	model.IsUser = false
	model.IsAdmin = false
	model.IsTenancy = false

	// passport
	model.TenantID = user.TenantID
	model.Passport = user.Passport
	model.Secret = user.Secret

	// userType
	switch user.UserType {
	case gw.AdministrationUser:
		model.IsAdmin = true
	case gw.TenancyUser:
		model.IsTenancy = true
	case gw.NormalUser:
		model.IsUser = true
	}
	tx := u.Store().GetDbStore()
	// references
	// https://gorm.io/zh_CN/docs/advanced_query.html
	err := tx.Where(Db.User{Passport: model.Passport}).FirstOrCreate(&model).Error
	user.ID = model.ID
	return err
}

func (u UserManager) Modify(user gw.User) error {
	db := u.Store().GetDbStore()
	return db.Model(&Db.User{}).Where("passport = ?", user.Passport).Update("secret", user.Password).Error
}

func (u UserManager) Delete(userId uint64) error {
	db := u.Store().GetDbStore()
	return db.Where("id = ?", userId).Delete(&Db.User{}).Error
}

func (u UserManager) QueryByUser(passport, password string) (gw.User, error) {
	var user gw.User
	var model Db.User
	var err = u.Backend().Take(&model, "passport = ?", passport).Error
	if err != nil || model.Secret != password {
		return gw.EmptyUser, err
	}

	user.ID = model.ID
	user.TenantID = model.TenantID
	user.Passport = model.Passport
	user.Password = model.Secret

	if user.IsEmpty() {
		return gw.EmptyUser, gw.ErrorUserNotFound
	}

	user.UserType = u.mapUserType(model)
	_, perms, err := u.serverState.PermissionManager().QueryByUser(model.TenantID, model.ID, gw.DefaultPageExpr)
	if err != nil {
		return gw.EmptyUser, err
	}
	user.Permissions = make([]*gw.Permission, len(perms))
	copy(user.Permissions, perms)
	return user, nil
}

func (u UserManager) QueryByAKS(accessKey, accessSecret string) (gw.User, error) {
	panic("implement me")
}

func (u UserManager) Query(userId uint64) (gw.User, error) {
	panic("implement me")
}

func (u UserManager) QueryList(tenantId uint64, expr gw.PagerExpr, total int64, out []gw.User) error {
	panic("implement me")
}

func DefaultUserManager(state *gw.ServerState) UserManager {
	var cnf = Config.GetUAP(state.ApplicationConfig())
	return UserManager{
		serverState:      state,
		cacheStoreName:   cnf.User.Cache.Name,
		backendStoreName: cnf.Backend.Name,
		cachePrefix:      cnf.User.Cache.Prefix,
		cacheExpiration:  time.Hour * time.Duration(cnf.User.Cache.ExpirationHours), // One week.
	}
}
