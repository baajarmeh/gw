package gw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw/conf"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Permission struct {
	ID         uint64
	TenantId   uint64
	Category   string
	Key        string
	Name       string
	Descriptor string
}

func (p Permission) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewPerm(key, name, descriptor string) Permission {
	return Permission{
		Key:        key,
		Name:       name,
		Descriptor: descriptor,
	}
}

var DefaultPermNames = []string{
	"ReadAll", "Creation",
	"Modification", "Deletion", "Disable", "ReadDetail",
}

//
// NewPermAll ...
// resource like: User, Role, Order etc.
// returns as
//  ReadAllUserPermission, CreationUserPermission,
//  ModificationUserPermission, DeletionUserPermission, DisableUserPermission, ReadDetailUserPermission etc.
//
func NewPermAll(resource string) []Permission {
	return NewPermByNames(resource, DefaultPermNames...)
}

//
// NewPermByNames ...
// resource like: User, Role, Order etc.
// permNames like: ReadAll, Creation, Modification,Deletion, Disable, ReadDetail etc.
//
func NewPermByNames(resource string, permNames ...string) []Permission {
	var perms []Permission
	for _, p := range permNames {
		kn := fmt.Sprintf("%s%sPermission", p, resource)
		desc := fmt.Sprintf("Define A %s %s permission", p, resource)
		perms = append(perms, Permission{
			Key:        kn,
			Name:       kn,
			Descriptor: desc,
		})
	}
	return perms
}

func NewPermSameKeyName(kn, descriptor string) Permission {
	return NewPerm(kn, kn, descriptor)
}

type IUserManager interface {
	Create(user *User) error
	Modify(user User) error
	Delete(tenantId, userId uint64) error
	Query(tenantId, userId uint64) (User, error)
	QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error
}

type IAuthManager interface {
	Login(passport, secret, credType, verifyCode string) (User, error)
	Logout(user User) bool
}

type ISessionStateManager interface {
	Save(sid string, user User) error
	Query(sid string) (User, error)
	Remove(sid string) error
}

type IPermissionManager interface {
	Initial()
	Has(user User, perms ...Permission) bool
	Create(category string, perms ...Permission) error
	Modify(perms ...Permission) error
	Drop(perms ...Permission) error
	Query(tenantId uint64, category string, expr PagerExpr) (total int64, result []Permission, error error)
	QueryByUser(tenantId, userId uint64, expr PagerExpr) (total int64, result []Permission, error error)
	GrantToUser(uid uint64, perms ...Permission) error
	GrantToRole(roleId uint64, perms ...Permission) error
	RevokeFromUser(uid uint64, perms ...Permission) error
	RevokeFromRole(roleId uint64, perms ...Permission) error
}

type DefaultAuthManagerImpl struct {
	store IStore
	cnf   conf.ApplicationConfig
	users map[string]*defaultUser
}

type defaultUser struct {
	User
	secret string
}

type DefaultSessionStateManagerImpl struct {
	store              IStore
	storeName          string
	storePrefix        string
	expirationDuration time.Duration
	cnf                conf.ApplicationConfig
	redisTimeout       time.Duration
}

func DefaultSessionStateManager(state ServerState) *DefaultSessionStateManagerImpl {
	defaultSs.store = state.Store()
	defaultSs.cnf = *state.ApplicationConfig()
	defaultSs.storeName = defaultSs.cnf.Security.Auth.Session.DefaultStore.Name
	defaultSs.storePrefix = defaultSs.cnf.Security.Auth.Session.DefaultStore.Prefix
	defaultSs.expirationDuration = time.Duration(defaultSs.cnf.Security.Auth.Cookie.MaxAge) * time.Second
	defaultSs.redisTimeout = time.Duration(defaultSs.cnf.Settings.TimeoutControl.Redis) * time.Millisecond
	return defaultSs
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
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return EmptyUser, err
	}
	return user, nil
}

var (
	EmptyUser       = User{Id: 0}
	ErrorEmptyInput = fmt.Errorf("empty input")
	defaultAm       *DefaultAuthManagerImpl
	defaultSs       *DefaultSessionStateManagerImpl
	defaultPm       *EmptyPermissionManagerImpl
)

func (d *DefaultAuthManagerImpl) Login(passport, secret, credType, verifyCode string) (User, error) {
	user, ok := d.users[passport]
	if ok && user.secret == secret {
		return user.User, nil
	}
	return EmptyUser, fmt.Errorf("user:%s not found or serect not match", passport)
}

func (d *DefaultAuthManagerImpl) Logout(user User) bool {
	// Nothing to do.
	return true
}

func DefaultAuthManager(state ServerState) *DefaultAuthManagerImpl {
	defaultAm.cnf = *state.ApplicationConfig()
	defaultAm.store = state.Store()
	return defaultAm
}

func init() {
	defaultAm = &DefaultAuthManagerImpl{
		users: map[string]*defaultUser{
			"admin": {
				User: User{
					Id:       10000,
					Passport: "admin",
					TenantId: 0,
					RoleId:   1,
				},
				secret: "123@456",
			},
			"gw-tenant-admin": {
				User: User{
					Id:       10001,
					Passport: "gw-tenant-admin",
					TenantId: 0,
					RoleId:   2,
				},
				secret: "123@456",
			},
			"gw-user1": {
				User: User{
					Id:       100000,
					Passport: "gw-user1",
					TenantId: 10001,
					RoleId:   100000,
				},
				secret: "123@456",
			},
		},
	}
	defaultSs = &DefaultSessionStateManagerImpl{}
}

type EmptyPermissionManagerImpl struct {
	state  int
	store  IStore
	conf   conf.ApplicationConfig
	locker sync.Mutex
	perms  map[string]map[string]Permission
}

func DefaultPermissionManager(state ServerState) *EmptyPermissionManagerImpl {
	if defaultPm == nil {
		defaultPm = &EmptyPermissionManagerImpl{
			conf:  *state.ApplicationConfig(),
			store: state.Store(),
			perms: make(map[string]map[string]Permission),
		}
	}
	return defaultPm
}

func (p *EmptyPermissionManagerImpl) getStore() *gorm.DB {
	return p.store.GetDbStoreByName(p.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (p *EmptyPermissionManagerImpl) Initial() {
}

func (p *EmptyPermissionManagerImpl) Has(user User, perms ...Permission) bool {
	return user.IsAuth()
}

func (p *EmptyPermissionManagerImpl) Create(category string, perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) Modify(perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) Drop(perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) Query(tenantId uint64, category string, expr PagerExpr) (
	total int64, result []Permission, error error) {
	return 0, nil, nil
}

func (p *EmptyPermissionManagerImpl) QueryByUser(tenantId, userId uint64, expr PagerExpr) (
	total int64, result []Permission, error error) {
	return 0, nil, nil
}

func (p *EmptyPermissionManagerImpl) GrantToUser(uid uint64, perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) GrantToRole(roleId uint64, perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) RevokeFromUser(uid uint64, perms ...Permission) error {
	return nil
}

func (p *EmptyPermissionManagerImpl) RevokeFromRole(roleId uint64, perms ...Permission) error {
	return nil
}

func DefaultUserManager(state ServerState) IUserManager {
	return EmptyUserManagerImpl{
		state: state,
	}
}

// EmptyUserManagerImpl ...
type EmptyUserManagerImpl struct {
	state ServerState
}

func (d EmptyUserManagerImpl) Create(user *User) error {
	pm := d.state.PermissionManager()
	return pm.GrantToUser(user.Id, user.Permissions...)
}

func (d EmptyUserManagerImpl) Modify(user User) error {
	return nil
}

func (d EmptyUserManagerImpl) Delete(tenantId, userId uint64) error {
	return nil
}

func (d EmptyUserManagerImpl) Query(tenantId, userId uint64) (User, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error {
	return nil
}

var defaultPageExpr = DefaultPagerExpr(1024, 1)

//
// GW framework login API.
func gwLogin(c *gin.Context) {
	s := *hostServer(c)
	reqId := getRequestID(s, c)
	pKey := s.conf.Security.Auth.ParamKey
	// supports
	// 1. User/Password
	// 2. X-Access-Key/X-Access-Secret
	// 3. Realm auth (Basic auth)
	passport, secret, verifyCode, credType, ok := parseCredentials(s, c)
	if !ok {
		c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(400, reqId, "Missing Credentials.", nil))
		c.Abort()
		return
	}

	// check params
	if p, ok := s.authParamValidators[pKey.Passport]; ok {
		if !p.MatchString(passport) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(400, reqId, "Invalid Credentials/Passport Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.Secret]; ok {
		if !p.MatchString(secret) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(400, reqId, "Invalid Credentials/Secret Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.VerifyCode]; ok {
		if !p.MatchString(verifyCode) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(400, reqId, "Invalid VerifyCode Formatter.", nil))
			c.Abort()
			return
		}
	}

	// Login
	user, err := s.AuthManager.Login(passport, secret, credType, verifyCode)
	if err != nil || user.IsEmpty() {
		c.JSON(http.StatusNotFound, s.RespBodyBuildFunc(404, reqId, err.Error(), nil))
		c.Abort()
		return
	}
	sid, ok := encryptSid(s, passport)
	if !ok {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(-1, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	if err := s.SessionStateManager.Save(user.Passport, user); err != nil {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(-1, reqId, "Save session fail.", err.Error()))
		c.Abort()
		return
	}
	_, perms, err := s.PermissionManager.QueryByUser(user.TenantId, user.Id, defaultPageExpr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(-1, reqId, "Query user's permission.", err.Error()))
		c.Abort()
		return
	}
	var userPerms []gin.H
	for _, p := range perms {
		userPerms = append(userPerms, gin.H{
			"Key":  p.Key,
			"Name": p.Name,
			"Desc": p.Descriptor,
		})
	}

	cks := s.conf.Security.Auth.Cookie
	expiredAt := time.Duration(cks.MaxAge) * time.Second
	var userRoles = gin.H{
		"Id":   user.RoleId,
		"name": "",
		"desc": "",
	}
	body := gin.H{
		"Credentials": gin.H{
			"Token":     sid,
			"ExpiredAt": time.Now().Add(expiredAt).Unix(),
		},
		"Roles":       userRoles,
		"Permissions": userPerms,
	}
	payload := s.RespBodyBuildFunc(0, reqId, nil, body)
	// token, header, X-Auth-Token
	c.Header("X-Auth-Token", sid)
	c.Header("Authorization", sid)
	c.SetCookie(cks.Key, sid, cks.MaxAge, cks.Path, cks.Domain, cks.Secure, cks.HttpOnly)
	c.JSON(http.StatusOK, payload)
}

// GW framework logout API.
func gwLogout(c *gin.Context) {
	s := *hostServer(c)
	reqId := getRequestID(s, c)
	user := getUser(c)
	cks := s.conf.Security.Auth.Cookie
	ok := s.AuthManager.Logout(user)
	if !ok {
		s.RespBodyBuildFunc(500, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		s.RespBodyBuildFunc(500, reqId, "session store logout fail", nil)
		return
	}
	_ = s.SessionStateManager.Remove(sid)
	c.SetCookie(cks.Key, "", -1, cks.Path, cks.Domain, cks.Secure, cks.HttpOnly)
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
		s := *hostServer(c)
		user := getUser(c)
		path := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		requestId := getRequestID(s, c)
		//
		// No auth and request URI not in allowed urls.
		// UnAuthorized
		//
		if (user.IsEmpty() || !user.IsAuth()) && !allowUrls[path] {
			auth := hostServer(c).conf.Security.AuthServer
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
			body := s.RespBodyBuildFunc(http.StatusUnauthorized, requestId, errDefault401Msg, payload)
			c.JSON(http.StatusUnauthorized, body)
			c.Abort()
			return
		}
		c.Next()
	}
}

// helpers
func parseCredentials(s HostServer, c *gin.Context) (passport, secret string, verifyCode string, credType string, result bool) {
	//
	// Auth Param configuration.
	param := s.conf.Security.Auth.ParamKey
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
	SecretHash      string
	RoleId          int // Platform Admin:1, Tenant Admin:2, roleId >= 10000 are custom role.
	ExtraRoleIdList []string
	Permissions     []Permission
}

func (user User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(&user)
}

func (user User) IsAuth() bool {
	return &user != nil && user.Id > 0
}

func (user User) IsEmpty() bool {
	return user.Id == EmptyUser.Id
}

func (user User) IsAdmin() bool {
	return user.IsAuth() && user.RoleId == 1
}

func (user User) IsTenantAdmin() bool {
	return user.IsAuth() && user.RoleId == 2
}

func getUser(c *gin.Context) User {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(User)
	}
	return EmptyUser
}
