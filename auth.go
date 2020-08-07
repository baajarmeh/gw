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

func NewPerm(key, name, descriptor string) Permission {
	return Permission{
		Key:        key,
		Name:       name,
		Descriptor: descriptor,
	}
}

func NewPermSameKeyName(kn, descriptor string) Permission {
	return NewPerm(kn, kn, descriptor)
}

func NewPermSameKeyNameDesc(knd string) Permission {
	return NewPermSameKeyName(knd, knd)
}

type IAuthManager interface {
	Login(store Store, passport, secret, credType, verifyCode string) (*User, error)
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
	users map[string]*defaultUser
}

type defaultUser struct {
	User
	secret string
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
	timeout := cnf.Service.Settings.TimeoutControl.Redis
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

func (d *DefaultAuthManagerImpl) Login(store Store, passport, secret, credType, verifyCode string) (*User, error) {
	user, ok := d.users[passport]
	if ok && user.secret == secret {
		return &user.User, nil
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
		users: map[string]*defaultUser{
			"admin": {
				User: User{
					Id:         10000,
					Passport:   "admin",
					TenantId:   0,
					MainRoleId: 1,
				},
				secret: "123@456",
			},
			"gw": {
				User: User{
					Id:         10001,
					Passport:   "gw",
					TenantId:   0,
					MainRoleId: 2,
				},
				secret: "123@456",
			},
			"gw-user1": {
				User: User{
					Id:         100000,
					Passport:   "gw-user1",
					TenantId:   10001,
					MainRoleId: 100000,
				},
				secret: "123@456",
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
		if user.IsAdmin() {
			return true
		}
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

// GW framework login API.
func gwLogin(c *gin.Context) {
	s := hostServer(c)
	reqId := getRequestID(s, c)
	pKey := s.conf.Service.Security.Auth.ParamKey
	// supports
	// 1. User/Password
	// 2. X-Access-Key/X-Access-Secret
	// 3. Realm auth (Basic auth)
	passport, secret, verifyCode, credType, ok := parseCredentials(s, c)
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
	user, err := s.AuthManager.Login(s.Store, passport, secret, credType, verifyCode)
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
	if err := s.SessionStateManager.Save(s.Store, user.Passport, user); err != nil {
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

// GW framework logout API.
func gwLogout(c *gin.Context) {
	s := hostServer(c)
	reqId := getRequestID(s, c)
	user := getUser(c)
	cks := s.conf.Service.Security.Auth.Cookie
	ok := s.AuthManager.Logout(s.Store, user)
	if !ok {
		respBody(500, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		respBody(500, reqId, "session store logout fail", nil)
		return
	}
	_ = s.SessionStateManager.Remove(s.Store, sid)
	domain := cks.Domain
	if domain == ":host" || domain == "" {
		domain = c.Request.Host
	}
	c.SetCookie(cks.Key, "", -1, cks.Path, domain, cks.Secure, cks.HttpOnly)
}

// GW framework auth Check Middleware
func gwAuthChecker(urls []conf.AllowUrl) gin.HandlerFunc {
	var allowUrls = make(map[string]bool)
	for _, url := range urls {
		for _, p := range url.Urls {
			s := p
			allowUrls[s] = true
		}
	}
	return func(c *gin.Context) {
		s := hostServer(c)
		user := getUser(c)
		path := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		requestId := getRequestID(s, c)
		//
		// No auth and request URI not in allowed urls.
		// UnAuthorized
		//
		if (user == nil || !user.IsAuth()) && !allowUrls[path] {
			auth := hostServer(c).conf.Service.Security.Auth.Server
			// Check url are allow dict.
			payload := gin.H{
				"Auth": gin.H{
					"LogIn": gin.H{
						"Url": fmt.Sprintf("%s/%s",
							strings.TrimRight(auth.Addr, "/"), strings.TrimLeft(auth.LogIn.Url, "/")),
						"Methods":   auth.LogIn.Methods,
						"AuthTypes": auth.LogIn.AuthTypes,
					},
					"LogOut": gin.H{
						"Url": fmt.Sprintf("%s/%s",
							strings.TrimRight(auth.Addr, "/"), strings.TrimLeft(auth.LogOut.Url, "/")),
						"Methods": auth.LogOut.Methods,
					},
				},
			}
			body := respBody(http.StatusUnauthorized, requestId, errDefault401Msg, payload)
			c.JSON(http.StatusUnauthorized, body)
			c.Abort()
			return
		}
		c.Next()
	}
}

// helpers
func parseCredentials(s *HostServer, c *gin.Context) (passport, secret string, verifyCode string, credType string, result bool) {
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
	credType = "passport"
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

	credType = "aks"
	// 2. X-Access-Key/X-Access-Secret
	passport = c.GetHeader("X-Access-Key")
	secret = c.GetHeader("X-Access-Secret")
	verifyCode = c.GetHeader("X-Access-VerifyCode")
	result = passport != "" && secret != ""
	if result {
		return
	}

	// 3. Basic auth
	credType = "account"
	passport, secret, result = c.Request.BasicAuth()
	if verifyCode == "" {
		verifyCode = c.Query(param.VerifyCode)
	}
	return
}

//
// Auth User
//
type User struct {
	Id              uint64
	TenantId        uint64
	Passport        string
	MainRoleId      int // Platform Admin:1, Tenant Admin:2, roleId >= 10000 are custom role.
	ExtraRoleIdList []string
	Permissions     []Permission
}

func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}

func (user *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(user)
}

func (user User) IsAuth() bool {
	return &user != nil && user.Id > 0
}

func (user User) IsAdmin() bool {
	return user.IsAuth() && user.MainRoleId == 1
}

func (user User) IsTenantAdmin() bool {
	return user.IsAuth() && user.MainRoleId == 2
}

func getUser(c *gin.Context) *User {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(*User)
	}
	return nil
}
