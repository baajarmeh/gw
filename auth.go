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
	Key        string `json:"key"`
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
}

func (perms Permission) String() string {
	b, _ := json.Marshal(perms)
	return string(b)
}

type IAuthManager interface {
	Login(store Store, passport, secret, verifyCode string) (*User, error)
	Logout(store Store, user *User) bool
}

type ISessionStateManager interface {
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

type DefaultSessionStateManagerImpl struct {
	storeName          string
	storePrefix        string
	expirationDuration time.Duration
	cnf                conf.Config
	redisTimeout       time.Duration
}

func DefaultSessionStateManager(cnf conf.Config) *DefaultSessionStateManagerImpl {
	defaultSs.cnf = cnf
	timeout := cnf.Service.Security.Timeout.Redis
	defaultSs.storeName = cnf.Service.Security.Auth.Session.DefaultStore.Name
	defaultSs.storePrefix = cnf.Service.Security.Auth.Session.DefaultStore.Prefix
	defaultSs.expirationDuration = time.Duration(cnf.Service.Security.Auth.Cookie.MaxAge) * time.Second
	defaultSs.redisTimeout = time.Duration(timeout) * time.Millisecond
	return defaultSs
}

func (d *DefaultSessionStateManagerImpl) context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, d.redisTimeout)
}

func (d *DefaultSessionStateManagerImpl) storeKey(sid string) string {
	return fmt.Sprintf("%s.%s", d.storePrefix, sid)
}

func (d *DefaultSessionStateManagerImpl) Remove(store Store, sid string) error {
	//FIXME(ocean): deadline executed error ?
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := store.GetCacheStoreByName(d.storeName)
	return redis.Del(ctx, d.storeKey(sid)).Err()
}

func (d *DefaultSessionStateManagerImpl) Save(store Store, sid string, user *User) error {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := store.GetCacheStoreByName(d.storeName)
	return redis.Set(ctx, d.storeKey(sid), user, d.expirationDuration).Err()
}

func (d *DefaultSessionStateManagerImpl) Query(store Store, sid string) (*User, error) {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	user := &User{}
	redis := store.GetCacheStoreByName(d.storeName)
	bytes, err := redis.Get(ctx, d.storeKey(sid)).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DefaultAuthManagerImpl) Login(store Store, passport, secret, verifyCode string) (*User, error) {
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
	defaultSs *DefaultSessionStateManagerImpl
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
	defaultSs = &DefaultSessionStateManagerImpl{}
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

const (
	gwUserKey = "gw-user"
)

func gwLogin(c *gin.Context) {
	s := hostServer(c)
	reqId := getRequestID(c)
	pKey := s.conf.Service.Security.Auth.ParamKey
	// supports
	// 1. User/Password
	// 2. X-Access-Key/X-Access-Secret
	// 3. Realm auth (Basic auth)
	passport, secret, verifyCode, ok := parseCredentials(s, c)
	if !ok {
		c.JSON(http.StatusBadRequest, respBody(400, reqId, "Invalid Credentials.", nil))
		c.Abort()
		return
	}

	// check params
	if p, ok := s.authParamValidators[pKey.Passport]; ok {
		if !p.MatchString(passport) {
			c.JSON(http.StatusBadRequest, respBody(400, reqId, "Invalid Credentials Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.Secret]; ok {
		if !p.MatchString(secret) {
			c.JSON(http.StatusBadRequest, respBody(400, reqId, "Invalid Credentials Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.VerifyCode]; ok {
		if !p.MatchString(verifyCode) {
			c.JSON(http.StatusBadRequest, respBody(400, reqId, "Invalid Credentials Formatter.", nil))
			c.Abort()
			return
		}
	}

	// Login
	user, err := s.authManager.Login(s.store, passport, secret, verifyCode)
	if err != nil || user == nil {
		c.JSON(http.StatusOK, respBody(400, reqId, err.Error(), nil))
		c.Abort()
		return
	}
	sid, ok := encryptSid(s, passport)
	if !ok {
		c.JSON(http.StatusInternalServerError, respBody(-1, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	if err := s.sessionStateManager.Save(s.store, user.Passport, user); err != nil {
		c.JSON(http.StatusInternalServerError, respBody(-1, reqId, "Save session fail.", err.Error()))
		c.Abort()
		return
	}
	cks := s.conf.Service.Security.Auth.Cookie
	domain := cks.Domain
	// if domain == ":host" || domain == "" {
	// 	domain = ""
	// }
	expiredAt := time.Duration(cks.MaxAge) * time.Second
	token := gin.H{
		"Token":     sid,
		"ExpiredAt": time.Now().Add(expiredAt).Unix(),
	}
	payload := respBody(0, reqId, nil, token)
	// token, header, X-Auth-Token
	c.Header("X-Auth-Token", sid)
	c.SetCookie(cks.Key, sid, cks.MaxAge, cks.Path, domain, cks.Secure, cks.HttpOnly)
	c.JSON(http.StatusOK, payload)
}

func gwLogout(c *gin.Context) {
	s := hostServer(c)
	reqId := getRequestID(c)
	user := getUser(c)
	cks := s.conf.Service.Security.Auth.Cookie
	ok := s.authManager.Logout(s.store, user)
	if !ok {
		respBody(500, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		respBody(500, reqId, "session store logout fail", nil)
		return
	}
	_ = s.sessionStateManager.Remove(s.store, sid)
	domain := cks.Domain
	if domain == ":host" || domain == "" {
		domain = c.Request.Host
	}
	c.SetCookie(cks.Key, "", -1, cks.Path, domain, cks.Secure, cks.HttpOnly)
}

func parseCredentials(s *HostServer, c *gin.Context) (passport, secret string, verifyCode string, result bool) {
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
	verifyCode, _ = c.GetPostForm(param.VerifyCode)
	result = passport != "" && secret != ""
	if result {
		return
	}

	// JSON
	if c.ContentType() == "application/json" {
		// json decode
		cred := gin.H{
			param.Passport:   "",
			param.Secret:     "",
			param.VerifyCode: "",
		}
		err := c.Bind(&cred)
		if err != nil {
			return
		}
		passport = cred[param.Passport].(string)
		secret = cred[param.Secret].(string)
		verifyCode = cred[param.VerifyCode].(string)
		result = passport != "" && secret != ""
		if result {
			return
		}
	}

	// 2. X-Access-Key/X-Access-Secret
	passport = c.GetHeader("X-Access-Key")
	secret = c.GetHeader("X-Access-Secret")
	verifyCode = c.GetHeader("X-Access-VerifyCode")
	result = passport != "" && secret != ""
	if result {
		return
	}

	// 3. Basic auth
	passport, secret, result = c.Request.BasicAuth()
	if verifyCode == "" {
		verifyCode = c.Query(param.VerifyCode)
	}
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
			auth := hostServer(c).conf.Service.Security.Auth
			authServer := strings.TrimRight(auth.AuthServer, "/")
			loginUrl := fmt.Sprintf("%s/%s", authServer, strings.TrimLeft(auth.LoginUrl, "/"))
			logoutUrl := fmt.Sprintf("%s/%s", authServer, strings.TrimLeft(auth.LogoutUrl, "/"))
			// Check url are allow dict.
			payload := gin.H{
				"Auth": gin.H{
					"LogIn": gin.H{
						"Url":    loginUrl,
						"Method": "GET/POST",
						"Types": []string{
							"User/Password",
							"X-Access-Key/Secret",
							"Basic Auth",
						},
					},
					"LogOut": gin.H{
						"Url":    logoutUrl,
						"Method": "GET/POST",
					},
				},
			}
			body := respBody(http.StatusUnauthorized, getRequestID(c), errDefault401Msg, payload)
			c.JSON(http.StatusUnauthorized, body)
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
