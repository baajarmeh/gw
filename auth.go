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
// resource like: AuthUser, Role, Order etc.
// returns as
//  ReadAllUserPermission, CreationUserPermission,
//  ModificationUserPermission, DeletionUserPermission, DisableUserPermission, ReadDetailUserPermission etc.
//
func NewPermAll(resource string) []Permission {
	return NewPermByNames(resource, DefaultPermNames...)
}

//
// NewPermByNames ...
// resource like: AuthUser, Role, Order etc.
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
	Create(user *AuthUser) error
	Modify(user AuthUser) error
	Delete(tenantId, userId uint64) error
	Query(tenantId, userId uint64) (AuthUser, error)
	QueryByUser(tenantId uint64, passport, password string) (AuthUser, error)
	QueryByAKS(tenantId uint64, accessKey, accessSecret string) (AuthUser, error)
	QueryList(tenantId uint64, expr PagerExpr, total int64, out []AuthUser) error
}

type IAuthManager interface {
	Login(tenantId uint64, passport, secret, verifyCode string, credType CredentialType) (AuthUser, error)
	Logout(user AuthUser) bool
}

type ISessionStateManager interface {
	Save(sid string, user AuthUser) error
	Query(sid string) (AuthUser, error)
	Remove(sid string) error
}

type IPermissionManager interface {
	Initial()
	Has(user AuthUser, perms ...Permission) bool
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
	AuthUser
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

func (d *DefaultSessionStateManagerImpl) Save(sid string, user AuthUser) error {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := d.store.GetCacheStoreByName(d.storeName)
	return redis.Set(ctx, d.storeKey(sid), user, d.expirationDuration).Err()
}

func (d *DefaultSessionStateManagerImpl) Query(sid string) (AuthUser, error) {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	user := AuthUser{}
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
	EmptyUser          = AuthUser{ID: 0}
	ErrorEmptyInput    = fmt.Errorf("empty input")
	ErrorUserNotFound  = fmt.Errorf("user not found")
	ErrorUserHasExists = fmt.Errorf("user object has exists")
	defaultAm          *DefaultAuthManagerImpl
	defaultSs          *DefaultSessionStateManagerImpl
	defaultPm          *EmptyPermissionManagerImpl
)

type CredentialType string

const (
	UserPassword    CredentialType = "user"
	AccessKeySecret CredentialType = "aks"
)

func (d *DefaultAuthManagerImpl) Login(tenantId uint64, passport, secret, verifyCode string, credType CredentialType) (AuthUser, error) {
	user, ok := d.users[passport]
	if ok && user.secret == secret {
		return user.AuthUser, nil
	}
	return EmptyUser, fmt.Errorf("user:%s not found or serect not match", passport)
}

func (d *DefaultAuthManagerImpl) Logout(user AuthUser) bool {
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
				AuthUser: AuthUser{
					ID:       10000,
					Passport: "admin",
					TenantId: 0,
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

func (p *EmptyPermissionManagerImpl) Has(user AuthUser, perms ...Permission) bool {
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

func (d EmptyUserManagerImpl) Create(user *AuthUser) error {
	pm := d.state.PermissionManager()
	return pm.GrantToUser(user.ID, user.Permissions...)
}

func (d EmptyUserManagerImpl) Modify(user AuthUser) error {
	return nil
}

func (d EmptyUserManagerImpl) Delete(tenantId, userId uint64) error {
	return nil
}

func (d EmptyUserManagerImpl) Query(tenantId, userId uint64) (AuthUser, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryByUser(tenantId uint64, user, password string) (AuthUser, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryByAKS(tenantId uint64, accessKey, accessSecret string) (AuthUser, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryList(tenantId uint64, expr PagerExpr, total int64, out []AuthUser) error {
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
	tenantId, passport, secret, verifyCode, credType, ok := parseCredentials(s, c)
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
	user, err := s.AuthManager.Login(tenantId, passport, secret, verifyCode, credType)
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
	_, perms, err := s.PermissionManager.QueryByUser(user.TenantId, user.ID, defaultPageExpr)
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
		"Id":   0,
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
func parseCredentials(s HostServer, c *gin.Context) (tenantId uint64, passport, secret string, verifyCode string, credType CredentialType, result bool) {
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

	tenantId = 0

	// 1. User/Password
	credType = UserPassword
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

	credType = AccessKeySecret
	// 2. X-Access-Key/X-Access-Secret
	passport = c.GetHeader("X-Access-Key")
	secret = c.GetHeader("X-Access-Secret")
	verifyCode = c.GetHeader("X-Access-VerifyCode")
	result = passport != "" && secret != ""
	if result {
		return
	}

	// 3. Basic auth
	credType = UserPassword
	passport, secret, result = c.Request.BasicAuth()
	if verifyCode == "" {
		verifyCode = c.Query(param.VerifyCode)
	}
	return
}

type UserType uint8

const (
	Administrator       UserType = 1
	TenantAdministrator UserType = 2
	User                UserType = 3
)

//
// Gw Auth User
//
type AuthUser struct {
	ID          uint64
	TenantId    uint64
	Passport    string
	Secret      string       `gorm:"-"`
	UserType    UserType     `gorm:"-"` // gw.AuthUserType
	Roles       []string     `gorm:"-"`
	Permissions []Permission `gorm:"-"`
}

func (user AuthUser) MarshalBinary() (data []byte, err error) {
	return json.Marshal(&user)
}

func (user AuthUser) IsAuth() bool {
	return &user != nil && user.ID > 0
}

func (user AuthUser) IsEmpty() bool {
	return user.ID == EmptyUser.ID
}

func (user AuthUser) IsAdmin() bool {
	return user.IsAuth() && user.UserType == Administrator
}

func (user AuthUser) IsTenantAdmin() bool {
	return user.IsAuth() && user.UserType == TenantAdministrator
}

func getUser(c *gin.Context) AuthUser {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(AuthUser)
	}
	return EmptyUser
}
