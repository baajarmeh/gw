package gwImpls

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"time"
)

type AuthManager struct {
	userCachePrefix  string
	cacheStoreName   string
	backendStoreName string
	expiration       time.Duration
	permPagerExpr    gw.PagerExpr
	state              gw.ServerState
}

func (a AuthManager) getUserCacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", a.userCachePrefix, passport)
}

func (a AuthManager) Login(passport, secret, credType, verifyCode string) (gw.User, error) {
	var gwUser gw.User
	var store = a.state.State.Store()()
	var password = a.state.PasswordSigner().Sign(secret)
	var cache = store.GetCacheStoreByName(a.cacheStoreName)
	var userCacheKey = a.getUserCacheKey(passport)
	var bytes, err = cache.Get(context.Background(), userCacheKey).Bytes()
	if err == nil && len(bytes) > 0 {
		err := json.Unmarshal(bytes, &gwUser)
		return gwUser, err
	}
	var user dbModel.User
	db := store.GetDbStoreByName(a.backendStoreName)
	err = db.Where("passport = ? and secret = ?", passport, password).First(&user).Error
	if err != nil {
		return gw.EmptyUser, err
	}
	// secret/password checker
	_, perms, err := a.state.PermissionManager().QueryByUser(user.TenantId, user.ID, a.permPagerExpr)
	if err != nil {
		return gw.EmptyUser, err
	}
	gwUser.Id = user.ID
	gwUser.TenantId = user.TenantId
	gwUser.Passport = user.Passport
	gwUser.Permissions = perms
	_ = cache.Set(context.Background(), userCacheKey, gwUser, a.expiration).Err()
	return gwUser, nil
}

func (a AuthManager) Logout(user gw.User) bool {
	var userCacheKey = a.getUserCacheKey(user.Passport)
	cache := a.state.State.Store()().GetCacheStoreByName(a.cacheStoreName)
	err := cache.Del(context.Background(), userCacheKey).Err()
	return err != nil
}

func DefaultAuthManager(state gw.ServerState) AuthManager {
	return AuthManager{
		state:              state,
		cacheStoreName:   "primary",
		backendStoreName: "primary",
		expiration:       time.Hour * 168, // One week.
	}
}
