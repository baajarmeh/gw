package gw

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Permission struct {
	Name       string
	Descriptor string
}

type IAuthManager interface {
	Login(store Store, passport, secret string) (*User, error)
	Logout(store Store, user *User) bool
}

type ISessionStateStore interface {
	Save(store Store, sid string, user *User) error
	Query(store Store, sid string) (*User, error)
	Remove(store Store, sid string) error
}

type IPermissionManager interface {
	HasPermission(user User, perms ...Permission) bool
	CreatePermissions(group string, perms ...Permission) bool
}

type DefaultAuthManagerImpl struct {
	users map[string]*User
}

type DefaultSessionStateStoreImpl struct {
	storeName          string
	storePrefix        string
	expirationDuration time.Duration
	cnf                conf.Config
	redisTimeout       time.Duration
}

func DefaultSessionStateStore(cnf conf.Config) *DefaultSessionStateStoreImpl {
	defaultSs.cnf = cnf
	timeout := cnf.Service.Security.Timeout.Redis
	defaultSs.storeName = cnf.Service.Security.Auth.Session.DefaultStore.Name
	defaultSs.storePrefix = cnf.Service.Security.Auth.Session.DefaultStore.Prefix
	defaultSs.expirationDuration = time.Duration(cnf.Service.Security.Auth.Cookie.MaxAge) * time.Second
	defaultSs.redisTimeout = time.Duration(timeout) * time.Millisecond
	return defaultSs
}

func (d *DefaultSessionStateStoreImpl) context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, d.redisTimeout)
}

func (d *DefaultSessionStateStoreImpl) storeKey(sid string) string {
	return fmt.Sprintf("%s.%s", d.storePrefix, sid)
}

func (d *DefaultSessionStateStoreImpl) Remove(store Store, sid string) error {
	//FIXME(ocean): deadline executed error ?
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := store.GetCacheStoreByName(d.storeName)
	return redis.Del(ctx, d.storeKey(sid)).Err()
}

func (d *DefaultSessionStateStoreImpl) Save(store Store, sid string, user *User) error {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := store.GetCacheStoreByName(d.storeName)
	return redis.Set(ctx, d.storeKey(sid), user, d.expirationDuration).Err()
}

func (d *DefaultSessionStateStoreImpl) Query(store Store, sid string) (*User, error) {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	user := &User{}
	redis := store.GetCacheStoreByName(d.storeName)
	bytes, err := redis.Get(ctx, d.storeKey(sid)).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes,user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DefaultAuthManagerImpl) Login(store Store, passport, secret string) (*User, error) {
	user, ok := d.users[passport]
	if ok && user.secret == secret {
		return user, nil
	}
	return nil, fmt.Errorf("user:%s not found or serect not match", passport)
}

func (d *DefaultAuthManagerImpl) Logout(store Store, user *User) bool {
	// Nothing to do.
	return true
}

var (
	defaultAm *DefaultAuthManagerImpl
	defaultSs *DefaultSessionStateStoreImpl
	defaultPm *DefaultPermissionManagerImpl
)

func init() {
	defaultAm = &DefaultAuthManagerImpl{
		users: map[string]*User{
			"admin": {
				Id:       10000,
				Passport: "admin",
				TenantId: 0,
				RoleId:   1,
				secret:   "123@456",
			},
			"gw": {
				Id:       10001,
				Passport: "gw",
				TenantId: 0,
				RoleId:   2,
				secret:   "123@456",
			},
			"gw-user1": {
				Id:       100000,
				Passport: "gw-user1",
				TenantId: 10001,
				RoleId:   3,
				secret:   "123@456",
			},
		},
	}
	defaultSs = &DefaultSessionStateStoreImpl{}
	defaultPm = &DefaultPermissionManagerImpl{
		perms: make(map[string]map[string]Permission),
	}
}

type DefaultPermissionManagerImpl struct {
	locker sync.Locker
	perms  map[string]map[string]Permission
}

func (p *DefaultPermissionManagerImpl) CreatePermissions(group string, perms ...Permission) bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	g := p.perms[group]
	for _, pm := range perms {
		if _, ok := g[pm.Name]; !ok {
			g[pm.Name] = pm
		}
	}
	return true
}

func (p *DefaultPermissionManagerImpl) HasPermission(user User, perms ...Permission) bool {
	if user.IsAuth() {
		for _, g := range p.perms {
			for _, pm := range perms {
				if _, ok := g[pm.Name]; ok {
					return true
				}
			}
		}
	}
	return false
}

const gwUserKey = "gw-user"

func gwLogin(c *gin.Context) {
	s := hostServer(c)
	store := s.store
	reqId := getRequestID(c)
	// supports
	// 1. User/Password
	// 2. X-Access-Key/X-Access-Secret
	// 3. Realm auth
	passport, secret, ok := parseCredentials(s, c)
	if !ok {
		c.JSON(http.StatusBadRequest, resp(-1, reqId, "No Credentials.", nil))
		c.Abort()
		return
	}
	sid, secureSid, ok := newSid(s)
	if !ok {
		c.JSON(http.StatusInternalServerError, resp(-1, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	user, err := s.authManager.Login(store, passport, secret)
	if err != nil {
		c.JSON(http.StatusOK, resp(-1, reqId, err.Error(), nil))
		c.Abort()
		return
	}
	if err := s.sessionStore.Save(store, sid, user); err != nil {
		c.JSON(http.StatusInternalServerError, resp(-1, reqId, "Save session fail.", nil))
		c.Abort()
		return
	}
	cks := s.conf.Service.Security.Auth.Cookie
	domain := cks.Domain
	if domain == ":host" || domain == "" {
		domain = c.Request.Host
	}
	c.SetCookie(cks.Key, secureSid, cks.MaxAge, cks.Path, domain, cks.Secure, cks.HttpOnly)
}

func gwLogout(c *gin.Context) {
	s := hostServer(c)
	reqId := getRequestID(c)
	user := getUser(c)
	store := s.store
	cks := s.conf.Service.Security.Auth.Cookie
	ok := s.authManager.Logout(store, user)
	if !ok {
		resp(-1, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		resp(-1, reqId, "session store logout fail", nil)
		return
	}
	_ = s.sessionStore.Remove(store, sid)
	domain := cks.Domain
	if domain == ":host" || domain == "" {
		domain = c.Request.Host
	}
	c.SetCookie(cks.Key, "", -1, cks.Path, domain, cks.Secure, cks.HttpOnly)
}

func parseCredentials(s *HostServer, c *gin.Context) (passport, secret string, result bool) {
	//
	// Auth Param configuration.
	param := s.conf.Service.Security.Auth.ParamKey
	//
	// supports
	// 1. User/Password
	// 2. X-Access-Key/X-Access-Secret
	// 3. Realm auth
	// ===============================
	//

	// 1. User/Password
	passport, _ = c.GetPostForm(param.Passport)
	secret, _ = c.GetPostForm(param.Secret)
	result = passport != "" && secret != ""
	if result {
		return
	}

	// JSON
	if c.ContentType() == "application/json" {
		// json decode
		cred := gin.H{
			param.Passport: "",
			param.Secret:   "",
		}
		err := c.Bind(&cred)
		if err != nil {
			return
		}
		passport = cred[param.Passport].(string)
		secret = cred[param.Secret].(string)
		result = passport != "" && secret != ""
		if result {
			return
		}
	}

	// 2. X-Access-Key/X-Access-Secret
	passport = c.GetHeader("X-Access-Key")
	secret = c.GetHeader("X-Access-Secret")
	result = passport != "" && secret != ""
	if result {
		return
	}

	// 3. Basic auth
	passport, secret, result = c.Request.BasicAuth()
	return
}

func gwAuthChecker(vars map[string]string, urls []conf.AllowUrl) gin.HandlerFunc {
	var allowUrls = make(map[string]bool)
	for _, url := range urls {
		for _, p := range url.Urls {
			s := p
			for k, v := range vars {
				s = strings.Replace(s, k, v, 4)
			}
			allowUrls[s] = true
		}
	}
	return func(c *gin.Context) {
		user := getUser(c)
		path := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		//
		// No auth and request URI not in allowed urls.
		// UnAuthorized
		//
		if (user == nil || !user.IsAuth()) && !allowUrls[path] {
			// Check url are allow dict.
			payload := resp(http.StatusUnauthorized, getRequestID(c), errDefault401Msg, errDefaultPayload)
			c.JSON(http.StatusUnauthorized, payload)
			c.Abort()
			return
		}
		c.Next()
	}
}

//type IUser interface {
//	Id() uint64
//	TenantId() uint64
//	IsAuth() bool
//	IsAdmin() bool
//	IsTenancyAdmin() bool
//	HasPerms(perms ...Permission) bool
//}

type User struct {
	Id       uint64
	TenantId uint64
	Passport string
	RoleId   int // Platform Admin:1, Tenant Admin:2, roleId >= 10000 are custom role.
	secret   string
}

func (usr *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, usr)
}

func (usr *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(usr)
}

func (usr User) IsAuth() bool {
	return &usr != nil && usr.Id > 0
}

func (usr User) IsAdmin() bool {
	return usr.IsAuth() && usr.RoleId == 1
}

func (usr User) IsTenantAdmin() bool {
	return usr.IsAuth() && usr.RoleId == 2
}

func getUser(c *gin.Context) *User {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(*User)
	}
	return nil
}
