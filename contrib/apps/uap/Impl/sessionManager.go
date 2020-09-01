package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Config"
	"github.com/oceanho/gw/logger"
	"time"
)

type SessionManager struct {
	*gw.ServerState
	cachePrefix     string
	cacheStoreName  string
	cacheExpiration time.Duration
	permPagerExpr   gw.PagerExpr
}

func (sm SessionManager) SessionCacheKey(sid string) string {
	return fmt.Sprintf("%s.%s", sm.cachePrefix, sid)
}

func (sm SessionManager) Cache() *redis.Client {
	return sm.Store().GetCacheStoreByName(sm.cacheStoreName)
}

func (sm SessionManager) GetSessionFromCache(sid string) (gw.User, error) {
	var bytes, err = sm.Cache().Get(context.Background(), sm.SessionCacheKey(sid)).Bytes()
	if err == nil {
		var user gw.User
		err = json.Unmarshal(bytes, &user)
		return user, err
	}
	return gw.EmptyUser, gw.ErrorSessionNotFound
}

func (sm SessionManager) RemoveSessionUserFromCache(sid string) error {
	return sm.Cache().Del(context.Background(), sm.SessionCacheKey(sid)).Err()
}

func (sm SessionManager) SaveSessionUserToCache(sid string, user gw.User) error {
	if err := sm.Cache().Set(context.Background(),
		sm.SessionCacheKey(sid), user, sm.cacheExpiration).Err(); err != nil {
		logger.Error("SaveSessionUserToCache fail, err: %v", err)
		return err
	}
	return nil
}

//
// APIs
//
func (sm SessionManager) Save(sid string, user gw.User) error {
	return sm.SaveSessionUserToCache(sid, user)
}

func (sm SessionManager) Query(sid string) (gw.User, error) {
	return sm.GetSessionFromCache(sid)
}

func (sm SessionManager) Remove(sid string) error {
	return sm.RemoveSessionUserFromCache(sid)
}

func DefaultSessionManager(state *gw.ServerState) SessionManager {
	var cnf = Config.GetUAP(state.ApplicationConfig())
	var sm = SessionManager{
		ServerState:     state,
		cachePrefix:     cnf.Session.Cache.Prefix,
		cacheStoreName:  cnf.Session.Cache.Name,
		cacheExpiration: time.Hour * time.Duration(cnf.Session.Cache.ExpirationHours), // One week.
	}
	return sm
}
