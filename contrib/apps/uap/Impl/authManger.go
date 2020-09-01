package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/logger"
	"gorm.io/gorm"
	"time"
)

type AuthManager struct {
	*gw.ServerState
	cachePrefix      string
	cacheStoreName   string
	backendStoreName string
	cacheExpiration  time.Duration
	permPagerExpr    gw.PagerExpr
}

func (a AuthManager) UserAuthCacheKey(passport string) string {
	return fmt.Sprintf("%s.%s", a.cachePrefix, passport)
}

func (a AuthManager) Password(secret string) string {
	return a.PasswordSigner().Sign(secret)
}

func (a AuthManager) Backend() *gorm.DB {
	return a.Store().GetDbStore()
}

func (a AuthManager) Cache() *redis.Client {
	return a.Store().GetCacheStoreByName(a.cacheStoreName)
}

func (a AuthManager) GetAuthUserFromCache(passport string) gw.User {
	var bytes, err = a.Cache().Get(context.Background(), a.UserAuthCacheKey(passport)).Bytes()
	if err == nil {
		var user gw.User
		if err := json.Unmarshal(bytes, &user); err == nil {
			return user
		}
	}
	return gw.EmptyUser
}

func (a AuthManager) RemoveAuthUserFromCache(passport string) {
	_ = a.Cache().Del(context.Background(), a.UserAuthCacheKey(passport))
}

func (a AuthManager) SaveAuthUserToCache(user gw.User) {
	if err := a.Cache().Set(context.Background(),
		a.UserAuthCacheKey(user.Passport), user, a.cacheExpiration).Err(); err != nil {
		logger.Error("SaveAuthUserToCache fail, err: %v", err)
	}
}

//
// Impl APIs
//
func (a AuthManager) Login(param gw.AuthParameter) (gw.User, error) {
	var err error = nil
	var user = a.GetAuthUserFromCache(param.Passport)
	var password = a.PasswordSigner().Sign(param.Password)
	if user.IsEmpty() {
		switch param.CredType {
		case gw.UserPasswordAuth, gw.BasicAuth:
			user, err = a.UserManager().QueryByUser(param.TenantId, param.Passport, password)
		case gw.AksAuth:
			user, err = a.UserManager().QueryByAKS(param.TenantId, param.Passport, password)
		default:
			logger.Error("Un-support cred type: %s", param.CredType)
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

func (a AuthManager) Logout(user gw.User) bool {
	return true
}

func DefaultAuthManager(state *gw.ServerState) AuthManager {
	var cnf = Config.GetUAP(state.ApplicationConfig())
	return AuthManager{
		ServerState:      state,
		backendStoreName: cnf.Auth.Backend.Name,
		cacheStoreName:   cnf.Auth.Cache.Name,
		cachePrefix:      cnf.Auth.Cache.Prefix,
		cacheExpiration:  time.Hour * time.Duration(cnf.Auth.Cache.ExpirationHours), // One week.
	}
}
