package gw

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/backend/gwdb"
	"github.com/oceanho/gw/conf"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Permission struct {
	gwdb.Model
	gwdb.HasTenantState
	Category   string `json:"category" gorm:"type:varchar(32);not null"`
	Key        string `json:"key" gorm:"type:varchar(64);not null"`
	Name       string `json:"name" gorm:"type:varchar(128); not null"`
	Descriptor string `json:"descriptor" gorm:"type:varchar(256)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Permission) TableName() string {
	return "gw_fw_permissions"
}

func (p Permission) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type PermissionMapping struct {
	gwdb.Model
	gwdb.HasTenantState
	gwdb.HasCreationState
	gwdb.HasModificationState
	Type         sql.NullInt32 `gorm:"not null"` // 1. User Permission, 2. Role/Group Permission
	ObjectID     uint64        `gorm:"not null"`
	PermissionID uint64        `gorm:"not null"`
}

var (
	userPermissionType = sql.NullInt32{
		Int32: 1,
	}
	rolePermissionType = sql.NullInt32{
		Int32: 1,
	}
)

func (PermissionMapping) TableName() string {
	return "gw_fw_permission_mappings"
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
func NewPermAll(resource string) []*Permission {
	return NewPermByNames(resource, DefaultPermNames...)
}

//
// NewPermByNames ...
// resource like: User, Role, Order etc.
// permNames like: ReadAll, Creation, Modification,Deletion, Disable, ReadDetail etc.
//
func NewPermByNames(resource string, permNames ...string) []*Permission {
	var perms []*Permission
	for _, p := range permNames {
		kn := fmt.Sprintf("%s%sPermission", p, resource)
		desc := fmt.Sprintf("Define A %s %s permission", p, resource)
		perms = append(perms, &Permission{
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

type IAuthManager interface {
	Login(store Store, passport, secret, credType, verifyCode string) (User, error)
	Logout(store Store, user User) bool
}

type ISessionStateManager interface {
	Save(store Store, sid string, user User) error
	Query(store Store, sid string) (User, error)
	Remove(store Store, sid string) error
}

type IPermissionManager interface {
	Initial()
	Has(user User, perms ...Permission) bool
	Create(category string, perms ...*Permission) error
	Modify(perms ...Permission) error
	Drop(perms ...Permission) error
	Query(tenantId uint64, category string, expr PagerExpr) (total int64, result []Permission, error error)
	GrantToUser(uid uint64, perms ...Permission) error
	GrantToRole(roleId uint64, perms ...Permission) error
	RevokeFromUser(uid uint64, perms ...Permission) error
	RevokeFromRole(roleId uint64, perms ...Permission) error
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
	cnf                conf.ApplicationConfig
	redisTimeout       time.Duration
}

func DefaultSessionStateManager(cnf conf.ApplicationConfig) *DefaultSessionStateManagerImpl {
	defaultSs.cnf = cnf
	timeout := cnf.Settings.TimeoutControl.Redis
	defaultSs.storeName = cnf.Security.Auth.Session.DefaultStore.Name
	defaultSs.storePrefix = cnf.Security.Auth.Session.DefaultStore.Prefix
	defaultSs.expirationDuration = time.Duration(cnf.Security.Auth.Cookie.MaxAge) * time.Second
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

func (d *DefaultSessionStateManagerImpl) Save(store Store, sid string, user User) error {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	redis := store.GetCacheStoreByName(d.storeName)
	return redis.Set(ctx, d.storeKey(sid), user, d.expirationDuration).Err()
}

func (d *DefaultSessionStateManagerImpl) Query(store Store, sid string) (User, error) {
	//ctx, cancel := d.context()
	//defer cancel()
	ctx := context.Background()
	user := User{}
	redis := store.GetCacheStoreByName(d.storeName)
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
)

func (d *DefaultAuthManagerImpl) Login(store Store, passport, secret, credType, verifyCode string) (User, error) {
	user, ok := d.users[passport]
	if ok && user.secret == secret {
		return user.User, nil
	}
	return EmptyUser, fmt.Errorf("user:%s not found or serect not match", passport)
}

func (d *DefaultAuthManagerImpl) Logout(store Store, user User) bool {
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
			"gw-tenant-admin": {
				User: User{
					Id:         10001,
					Passport:   "gw-tenant-admin",
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
}

type DefaultPermissionManagerImpl struct {
	state  int
	store  Store
	conf   conf.ApplicationConfig
	locker sync.Mutex
	perms  map[string]map[string]Permission
}

func DefaultPermissionManager(conf conf.ApplicationConfig, store Store) *DefaultPermissionManagerImpl {
	if defaultPm == nil {
		defaultPm = &DefaultPermissionManagerImpl{
			conf:  conf,
			store: store,
			perms: make(map[string]map[string]Permission),
		}
	}
	return defaultPm
}

func (p *DefaultPermissionManagerImpl) getStore() *gorm.DB {
	return p.store.GetDbStoreByName(p.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (p *DefaultPermissionManagerImpl) Initial() {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.state > 0 {
		return
	}
	store := p.getStore()
	var tables []interface{}
	tables = append(tables, &Permission{})
	tables = append(tables, &PermissionMapping{})
	err := store.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}

func (p *DefaultPermissionManagerImpl) Has(user User, perms ...Permission) bool {
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

func (p *DefaultPermissionManagerImpl) Create(category string, perms ...*Permission) error {
	if len(perms) < 1 {
		return ErrorEmptyInput
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	db := p.getStore()
	tx := db.Begin()
	var registered = make(map[string]bool)
	for i := 0; i < len(perms); i++ {
		p := perms[i]
		uk := fmt.Sprintf("%d-%s-%s", p.TenantId, p.Category, p.Key)
		if registered[uk] {
			continue
		}
		var count int64
		if p.Category == "" {
			p.Category = category
		}
		err := db.Model(Permission{}).Where("tenant_id = ? and category = ? and `key` = ?",
			p.TenantId, p.Category, p.Key).Count(&count).Error
		if err != nil {
			panic(fmt.Sprintf("check perms has exist fail, err: %v", err))
		}
		if count == 0 {
			tx.Create(p)
		}
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) Modify(perms ...Permission) error {
	if len(perms) < 1 {
		return ErrorEmptyInput
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		tx.Updates(p)
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) Drop(perms ...Permission) error {
	if len(perms) < 1 {
		return ErrorEmptyInput
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		tx.Delete(p)
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) Query(tenantId uint64, category string, expr PagerExpr) (
	total int64, result []Permission, error error) {
	error = p.getStore().Where("tenant_id = ? and category = ?",
		tenantId, category).Count(&total).Offset(expr.PageOffset()).Limit(expr.PageSize).Scan(result).Error
	return
}

func (p *DefaultPermissionManagerImpl) GrantToUser(uid uint64, perms ...Permission) error {
	if len(perms) < 1 {
		return nil
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		var pm PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = userPermissionType
		pm.ObjectID = uid
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) GrantToRole(roleId uint64, perms ...Permission) error {
	if len(perms) < 1 {
		return nil
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		var pm PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = rolePermissionType
		pm.ObjectID = roleId
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) RevokeFromUser(uid uint64, perms ...Permission) error {
	if len(perms) < 1 {
		return nil
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		tx.Delete(PermissionMapping{},
			"object_id = ? and tenant_id = ? and permission_id = ? and type = ?", uid, p.TenantId, p.ID, userPermissionType)
	}
	return tx.Commit().Error
}

func (p *DefaultPermissionManagerImpl) RevokeFromRole(roleId uint64, perms ...Permission) error {
	if len(perms) < 1 {
		return nil
	}
	p.locker.Lock()
	defer p.locker.Unlock()
	tx := p.getStore().Begin()
	for _, p := range perms {
		tx.Delete(PermissionMapping{},
			"object_id = ? and tenant_id = ? and permission_id = ? and type = ?", roleId, p.TenantId, p.ID, rolePermissionType)
	}
	return tx.Commit().Error
}

//
// GW framework login API.
func gwLogin(c *gin.Context) {
	s := hostServer(c)
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
	user, err := s.AuthManager.Login(s.Store, passport, secret, credType, verifyCode)
	if err != nil || user.Id == EmptyUser.Id {
		c.JSON(http.StatusOK, s.RespBodyBuildFunc(400, reqId, err.Error(), nil))
		c.Abort()
		return
	}
	sid, ok := encryptSid(s, passport)
	if !ok {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(-1, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	if err := s.SessionStateManager.Save(s.Store, user.Passport, user); err != nil {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(-1, reqId, "Save session fail.", err.Error()))
		c.Abort()
		return
	}
	cks := s.conf.Security.Auth.Cookie
	domain := cks.Domain
	// if domain == ":host" || domain == "" {
	// 	domain = ""
	// }
	expiredAt := time.Duration(cks.MaxAge) * time.Second
	token := gin.H{
		"Token":     sid,
		"ExpiredAt": time.Now().Add(expiredAt).Unix(),
	}
	payload := s.RespBodyBuildFunc(0, reqId, nil, token)
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
	cks := s.conf.Security.Auth.Cookie
	ok := s.AuthManager.Logout(s.Store, user)
	if !ok {
		s.RespBodyBuildFunc(500, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		s.RespBodyBuildFunc(500, reqId, "session store logout fail", nil)
		return
	}
	_ = s.SessionStateManager.Remove(s.Store, sid)
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
		s := hostServer(c)
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
	MainRoleId      int // Platform Admin:1, Tenant Admin:2, roleId >= 10000 are custom role.
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
	return user.IsAuth() && user.MainRoleId == 1
}

func (user User) IsTenantAdmin() bool {
	return user.IsAuth() && user.MainRoleId == 2
}

func getUser(c *gin.Context) User {
	obj, ok := c.Get(gwUserKey)
	if ok {
		return obj.(User)
	}
	return EmptyUser
}
