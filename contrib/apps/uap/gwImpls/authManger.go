package gwImpls

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type AuthManager struct {
	gw.ServerState
	userAuthCachePrefix string
	cacheStoreName      string
	backendStoreName    string
	expiration          time.Duration
	permPagerExpr       gw.PagerExpr
}

func (a AuthManager) UserAuthCacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", a.userAuthCachePrefix, passport)
}

func (a AuthManager) UserSessionCacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", a.userAuthCachePrefix, passport)
}

func (a AuthManager) Backend() *gorm.DB {
	return a.Store().GetDbStore()
}

func (a AuthManager) Cache() *redis.Client {
	return a.Store().GetCacheStoreByName(a.cacheStoreName)
}

func (a AuthManager) GetAuthUserFromCache(passport string) gw.AuthUser {
	var bytes, err = a.Cache().Get(context.Background(), a.UserAuthCacheKey(passport)).Bytes()
	if err == nil {
		var user gw.AuthUser
		if err := json.Unmarshal(bytes, &user); err == nil {
			return user
		}
	}
	return gw.EmptyUser
}

func (a AuthManager) RemoveAuthUserFromCache(passport string) {
	_ = a.Cache().Del(context.Background(), a.UserAuthCacheKey(passport))
}

func (a AuthManager) Password(secret string) string {
	return a.PasswordSigner().Sign(secret)
}

func (a AuthManager) SaveAuthUserToCache(user gw.AuthUser) {
	if err := a.Cache().Set(context.Background(),
		a.UserAuthCacheKey(user.Passport), user, a.expiration).Err(); err != nil {
		logger.Error("SaveAuthUserToCache fail, err: %v", err)
	}
}

//
// Impl APIs
//
func (a AuthManager) Login(tenantId uint64, passport, secret, verifyCode string, credType gw.CredentialType) (gw.AuthUser, error) {
	var err error = nil
	var user = a.GetAuthUserFromCache(passport)
	var password = a.PasswordSigner().Sign(secret)
	if user.IsEmpty() {
		if credType == gw.UserPassword {
			user, err = a.UserManager().QueryByUser(tenantId, passport, password)
		} else if credType == gw.AccessKeySecret {
			user, err = a.UserManager().QueryByAKS(tenantId, passport, password)
		} else {
			logger.Error("Un-support cred type: %s", credType)
			return gw.EmptyUser, err
		}
	}
	if err != nil {
		return gw.EmptyUser, err
	}
	if user.IsEmpty() {
		return user, gw.ErrorUserNotFound
	}
	_, perms, err := a.PermissionManager().QueryByUser(user.TenantId, user.ID, a.permPagerExpr)
	if err != nil {
		return gw.EmptyUser, err
	}
	user.Permissions = perms
	a.SaveAuthUserToCache(user)
	return user, nil
}

func (a AuthManager) Logout(user gw.AuthUser) bool {
	return true
}

func DefaultAuthManager(state gw.ServerState) AuthManager {
	return AuthManager{
		ServerState:         state,
		cacheStoreName:      "primary",
		backendStoreName:    "primary",
		userAuthCachePrefix: "gw-uap-auth",
		expiration:          time.Hour * 168, // One week.
	}
}
