package gw

import (
	"context"
	"fmt"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/libs/gwjsoner"
	"time"
)

type DefaultSessionStateManagerImpl struct {
	store              IStore
	storeName          string
	storePrefix        string
	expirationDuration time.Duration
	redisTimeout       time.Duration
	cnf                *conf.ApplicationConfig
}

func DefaultSessionStateManager(state *ServerState) *DefaultSessionStateManagerImpl {
	var stateManager = &DefaultSessionStateManagerImpl{}
	stateManager.store = state.Store()
	stateManager.cnf = state.ApplicationConfig()
	stateManager.storeName = stateManager.cnf.Security.Auth.Session.DefaultStore.Name
	stateManager.storePrefix = stateManager.cnf.Security.Auth.Session.DefaultStore.Prefix
	stateManager.expirationDuration = time.Duration(stateManager.cnf.Security.Auth.Cookie.MaxAge) * time.Second
	stateManager.redisTimeout = time.Duration(stateManager.cnf.Settings.TimeoutControl.Redis) * time.Millisecond
	return stateManager
}

func (d *DefaultSessionStateManagerImpl) context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, d.redisTimeout)
}

func (d *DefaultSessionStateManagerImpl) storeKey(sid string) string {
	return fmt.Sprintf("%s.%s", d.storePrefix, sid)
}

func (d *DefaultSessionStateManagerImpl) Remove(sid string) error {
	//FIXME(ocean): deadline executed error ?
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := d.store.GetCacheStoreByName(d.cnf.Security.Auth.Session.DefaultStore.Name)
	return redis.Del(ctx, d.storeKey(sid)).Err()
}

func (d *DefaultSessionStateManagerImpl) Save(sid string, user User) error {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := d.store.GetCacheStoreByName(d.storeName)
	return redis.Set(ctx, d.storeKey(sid), user, d.expirationDuration).Err()
}

func (d *DefaultSessionStateManagerImpl) Query(sid string) (User, error) {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	user := User{}
	redis := d.store.GetCacheStoreByName(d.storeName)
	bytes, err := redis.Get(ctx, d.storeKey(sid)).Bytes()
	if err != nil {
		return EmptyUser, err
	}
	err = gwjsoner.Unmarshal(bytes, &user)
	if err != nil {
		return EmptyUser, err
	}
	return user, nil
}
