package gwImpls

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"github.com/oceanho/gw/libs/gwjsoner"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type UserManager struct {
	gw.ServerState
	userCachePrefix  string
	cacheStoreName   string
	backendStoreName string
	expiration       time.Duration
	permPagerExpr    gw.PagerExpr
}

func (u UserManager) Backend() *gorm.DB {
	return u.Store().GetDbStoreByName(u.backendStoreName)
}

func (u UserManager) Cache() *redis.Client {
	return u.Store().GetCacheStoreByName(u.cacheStoreName)
}

func (u UserManager) User() *gorm.DB {
	return u.Backend().Model(dbModel.User{})
}

func (u UserManager) GetFromCache(passport string) gw.AuthUser {
	var bytes, err = u.Cache().Get(context.Background(), u.CacheKey(passport)).Bytes()
	if err != nil {
		logger.Error("operation cache fail, err: %v", err)
		return gw.EmptyUser
	}
	var user gw.AuthUser
	if gwjsoner.Unmarshal(bytes, &user) != nil {
		return gw.EmptyUser
	}
	return user
}

func (u UserManager) CacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", u.userCachePrefix, passport)
}

func (u UserManager) SaveToCache(user gw.AuthUser) {
}

//
// APIs
//
func (u UserManager) Create(user *gw.AuthUser) error {
	store := u.Store()
	db := store.GetDbStore()
	var model dbModel.User
	err := db.First(&model, "tenant_id=? and passport=?", user.TenantId, user.Passport).Error
	if err != nil && err.Error() != "record not found" {
		return err
	}
	if model.ID > 0 {
		return gw.ErrorUserHasExists
	}

	// default
	model.IsUser = false
	model.IsAdmin = false
	model.IsTenancy = false

	// passport
	model.Passport = user.Passport
	model.TenantId = user.TenantId
	model.Secret = user.Secret

	switch user.UserType {
	case gw.Administrator:
		model.IsAdmin = true
		break
	case gw.TenantAdministrator:
		model.IsTenancy = true
		break
	case gw.User:
		model.IsUser = true
		break
	default:
		model.IsUser = false
		model.IsAdmin = false
		model.IsTenancy = false
	}
	tx := store.GetDbStore().Begin()
	err = tx.Create(&model).Error
	if err != nil {
		return err
	}
	err = tx.Commit().Error
	user.ID = model.ID
	return err
}

func (u UserManager) Modify(user gw.AuthUser) error {
	panic("implement me")
}

func (u UserManager) Delete(tenantId, userId uint64) error {
	panic("implement me")
}

func (u UserManager) QueryByUser(tenantId uint64, passport, password string) (gw.AuthUser, error) {
	var user gw.AuthUser
	var err = u.User().Find(&user, "tenant_id=? and passport=? and secret=?", tenantId, passport, password).Error
	return user, err
}

func (u UserManager) QueryByAKS(tenantId uint64, accessKey, accessSecret string) (gw.AuthUser, error) {
	panic("implement me")
}

func (u UserManager) Query(tenantId, userId uint64) (gw.AuthUser, error) {
	panic("implement me")
}

func (u UserManager) QueryList(tenantId uint64, expr gw.PagerExpr, total int64, out []gw.AuthUser) error {
	panic("implement me")
}

func DefaultUserManager(state gw.ServerState) UserManager {
	return UserManager{
		ServerState:      state,
		cacheStoreName:   "primary",
		backendStoreName: "primary",
		userCachePrefix:  "gw-uap-user",
		expiration:       time.Hour * 168, // One week.
	}
}
