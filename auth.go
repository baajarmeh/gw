package gw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/libs/gwjsoner"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	b, _ := gwjsoner.Marshal(p)
	return string(b)
}

func (p Permission) IdStr() string {
	if p.TenantId > 0 {
		return fmt.Sprintf("%d.%s.%s", p.TenantId, p.Category, p.Key)
	}
	return fmt.Sprintf("%s.%s", p.Category, p.Key)
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
	Create(user *User) error
	Modify(user User) error
	Delete(tenantId, userId uint64) error
	Query(tenantId, userId uint64) (User, error)
	QueryByUser(tenantId uint64, passport, password string) (User, error)
	QueryByAKS(tenantId uint64, accessKey, accessSecret string) (User, error)
	QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error
}

type IAuthManager interface {
	Login(tenantId uint64, passport, secret, verifyCode string, credType CredentialType) (User, error)
	Logout(user User) bool
}

type ISessionStateManager interface {
	Save(sid string, user User) error
	Query(sid string) (User, error)
	Remove(sid string) error
}

type IPermissionChecker interface {
	Check(user User, perms ...Permission) bool
}
type IPermissionManager interface {
	Initial()
	Checker() IPermissionChecker
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
	err = gwjsoner.Unmarshal(bytes, &user)
	if err != nil {
		return EmptyUser, err
	}
	return user, nil
}

var (
	EmptyUser            = User{ID: 0}
	ErrorEmptyInput      = fmt.Errorf("empty input")
	ErrorUserNotFound    = fmt.Errorf("user not found")
	ErrorSessionNotFound = fmt.Errorf("session not found")
	ErrorUserHasExists   = fmt.Errorf("user object has exists")
	defaultAm            *DefaultAuthManagerImpl
	defaultSs            *DefaultSessionStateManagerImpl
	defaultPm            *EmptyPermissionManagerImpl
	defaultChecker       *DefaultPassPermissionChecker
)

type CredentialType string

const (
	UserPassword    CredentialType = "user"
	AccessKeySecret CredentialType = "aks"
)

func (d *DefaultAuthManagerImpl) Login(tenantId uint64, passport, secret, verifyCode string, credType CredentialType) (User, error) {
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
	state             int
	store             IStore
	conf              conf.ApplicationConfig
	locker            sync.Mutex
	perms             map[string]map[string]Permission
	permissionChecker IPermissionChecker
}

func DefaultPermissionManager(state ServerState) *EmptyPermissionManagerImpl {
	if defaultPm == nil {
		defaultPm = &EmptyPermissionManagerImpl{
			conf:  *state.ApplicationConfig(),
			store: state.Store(),
			perms: make(map[string]map[string]Permission),
			permissionChecker: DefaultPassPermissionChecker{
				State: state,
			},
		}
	}
	return defaultPm
}

type DefaultPassPermissionChecker struct {
	State           ServerState
	CustomCheckFunc func(user User, perms ...Permission) bool
}

func (a DefaultPassPermissionChecker) Check(user User, perms ...Permission) bool {
	if a.CustomCheckFunc != nil {
		return a.CustomCheckFunc(user, perms...)
	}
	if user.IsAdmin() {
		return true
	}
	if !user.IsAuth() || user.PermMaps == nil {
		return false
	}
	for _, p := range perms {
		if _, ok := user.PermMaps[p.IdStr()]; ok {
			return true
		}
	}
	return false
}

func (p *EmptyPermissionManagerImpl) getStore() *gorm.DB {
	return p.store.GetDbStoreByName(p.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (p *EmptyPermissionManagerImpl) Initial() {
}

func (p *EmptyPermissionManagerImpl) Checker() IPermissionChecker {
	return p.permissionChecker
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
	return pm.GrantToUser(user.ID, user.Permissions...)
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

func (d EmptyUserManagerImpl) QueryByUser(tenantId uint64, user, password string) (User, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryByAKS(tenantId uint64, accessKey, accessSecret string) (User, error) {
	return EmptyUser, nil
}

func (d EmptyUserManagerImpl) QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error {
	return nil
}

var DefaultPageExpr = DefaultPagerExpr(1024, 1)

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
		c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(http.StatusBadRequest, reqId, "Missing Credentials.", nil))
		c.Abort()
		return
	}

	// check params
	if p, ok := s.authParamValidators[pKey.Passport]; ok {
		if !p.MatchString(passport) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(http.StatusBadRequest, reqId, "Invalid Credentials/Passport Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.Secret]; ok {
		if !p.MatchString(secret) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(http.StatusBadRequest, reqId, "Invalid Credentials/Secret Formatter.", nil))
			c.Abort()
			return
		}
	}
	if p, ok := s.authParamValidators[pKey.VerifyCode]; ok {
		if !p.MatchString(verifyCode) {
			c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(http.StatusBadRequest, reqId, "Invalid VerifyCode Formatter.", nil))
			c.Abort()
			return
		}
	}

	// Login
	user, err := s.AuthManager.Login(tenantId, passport, secret, verifyCode, credType)
	if err != nil || user.IsEmpty() {
		c.JSON(http.StatusNotFound, s.RespBodyBuildFunc(http.StatusNotFound, reqId, err.Error(), nil))
		c.Abort()
		return
	}
	sid, credential, ok := encryptSid(s, passport)
	if !ok {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(http.StatusInternalServerError, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	if err := s.SessionStateManager.Save(sid, user); err != nil {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(http.StatusInternalServerError, reqId, "Save session fail.", err.Error()))
		c.Abort()
		return
	}
	var userPerms []gin.H
	for _, p := range user.Permissions {
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
	payload := gin.H{
		"Credentials": gin.H{
			"Token":     credential,
			"ExpiredAt": time.Now().Add(expiredAt).Unix(),
		},
		"Roles":       userRoles,
		"Permissions": userPerms,
	}
	body := s.RespBodyBuildFunc(0, reqId, nil, payload)
	c.SetCookie(cks.Key, credential, cks.MaxAge, cks.Path, cks.Domain, cks.Secure, cks.HttpOnly)
	c.JSON(http.StatusOK, body)
}

// GW framework logout API.
func gwLogout(c *gin.Context) {
	s := *hostServer(c)
	reqId := getRequestID(s, c)
	user := getUser(c)
	cks := s.conf.Security.Auth.Cookie
	ok := s.AuthManager.Logout(user)
	if !ok {
		s.RespBodyBuildFunc(http.StatusInternalServerError, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		s.RespBodyBuildFunc(http.StatusInternalServerError, reqId, "session store logout fail", nil)
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
func parseCredentials(s HostServer, c *gin.Context) (tenantId uint64,
	passport, secret string, verifyCode string, credType CredentialType, result bool) {
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
	var tenantIdStr = ""

	// 1. User/Password
	credType = UserPassword
	passport, _ = c.GetPostForm(param.Passport)
	secret, _ = c.GetPostForm(param.Secret)
	tenantIdStr, _ = c.GetPostForm(param.TenantId)
	verifyCode, _ = c.GetPostForm(param.VerifyCode)
	result = passport != "" && secret != ""
	if result {
		tenantId, _ = strconv.ParseUint(tenantIdStr, 10, 64)
		return
	}

	// JSON
	if c.ContentType() == "application/json" {
		// json decode
		cred := gin.H{
			param.Passport:   "",
			param.Secret:     "",
			param.VerifyCode: "",
			param.TenantId:   "",
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

	tenantIdStr = c.GetHeader("X-GW-Tenant-ID")
	// 2. Basic auth
	credType = UserPassword
	passport, secret, result = c.Request.BasicAuth()
	if verifyCode == "" {
		verifyCode = c.Query(param.VerifyCode)
	}
	result = passport != "" && secret != ""
	if result {
		tenantId, _ = strconv.ParseUint(tenantIdStr, 10, 64)
		return
	}

	credType = AccessKeySecret
	// 3. X-Aks-Key/X-Aks-Secret
	passport = c.GetHeader("X-Aks-Key")
	secret = c.GetHeader("X-Aks-Secret")
	verifyCode = c.GetHeader("X-Aks-VerifyCode")
	tenantId, _ = strconv.ParseUint(tenantIdStr, 10, 64)
	return
}

type UserType uint8

const (
	Administrator UserType = 1
	Tenancy       UserType = 2
	NonUser       UserType = 3
)

func (ut UserType) IsAdmin() bool {
	return ut == Administrator
}

func (ut UserType) IsTenancy() bool {
	return ut == Tenancy
}

func (ut UserType) IsUser() bool {
	return ut == NonUser
}

//
// Gw User
//
type User struct {
	ID          uint64
	TenantId    uint64
	Passport    string
	Secret      string                `gorm:"-"`
	Password    string                `gorm:"-"` // Hash string
	UserType    UserType              `gorm:"-"` // gw.AuthUserType
	Roles       []string              `gorm:"-"`
	Permissions []Permission          `gorm:"-"`
	PermMaps    map[string]Permission `gorm:"-"`
}

//
// Gw Session
//
type Session struct {
	User
	AccessKey    string
	AccessSecret string
	SessionID    string
	CreatedAt    *time.Time
	LastActive   *time.Time
}

func (user User) MarshalBinary() (data []byte, err error) {
	return gwjsoner.Marshal(&user)
}

func (user User) IsAuth() bool {
	return &user != nil && user.ID > 0
}

func (user User) IsEmpty() bool {
	return user.ID == EmptyUser.ID
}

func (user User) IsAdmin() bool {
	return user.IsAuth() && user.UserType == Administrator
}

func (user User) IsTenancy() bool {
	return user.IsAuth() && user.UserType == Tenancy
}

func (user User) IsUser() bool {
	return user.IsAuth() && user.UserType == NonUser
}

func getUser(c *gin.Context) User {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(User)
	}
	return EmptyUser
}
