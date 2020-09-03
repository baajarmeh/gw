package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/libs/gwjsoner"
	"gorm.io/gorm"
	"sync"
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

type AuthParameter struct {
	Passport   string
	Password   string
	TenantId   uint64
	VerifyCode string
	CredType   CredentialType
}

type IAuthManager interface {
	Login(param AuthParameter) (User, error)
	Logout(user User) bool
}

type IAuthParamResolver interface {
	Resolve(ctx *gin.Context) AuthParameter
}

type IAuthParamChecker interface {
	Check(param AuthParameter) error
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

var (
	EmptyUser            = User{ID: 0}
	ErrorEmptyInput      = fmt.Errorf("empty input")
	ErrorUserNotFound    = fmt.Errorf("user not found")
	ErrorInvalidParamID  = fmt.Errorf("invalid param id")
	ErrorSessionNotFound = fmt.Errorf("session not found")
	ErrorUserHasExists   = fmt.Errorf("user object has exists")
)

type CredentialType string

const (
	AksAuth          CredentialType = "aks"
	BasicAuth        CredentialType = "basic"
	UserPasswordAuth CredentialType = "user"
)

type DefaultPermissionManagerImpl struct {
	state             int
	store             IStore
	conf              conf.ApplicationConfig
	locker            sync.Mutex
	perms             map[string]map[string]Permission
	permissionChecker IPermissionChecker
}

func DefaultPermissionManager(state *ServerState) *DefaultPermissionManagerImpl {
	var defaultPm = &DefaultPermissionManagerImpl{
		conf:              *state.ApplicationConfig(),
		store:             state.Store(),
		perms:             make(map[string]map[string]Permission),
		permissionChecker: state.PermissionChecker(),
	}
	return defaultPm
}

type DefaultPassPermissionCheckerImpl struct {
	State           *ServerState
	CustomCheckFunc func(user User, perms ...Permission) bool
}

func DefaultPassPermissionChecker(state *ServerState) IPermissionChecker {
	return DefaultPassPermissionCheckerImpl{
		State:           state,
		CustomCheckFunc: nil,
	}
}

func (a DefaultPassPermissionCheckerImpl) Check(user User, perms ...Permission) bool {
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

func (p *DefaultPermissionManagerImpl) getStore() *gorm.DB {
	return p.store.GetDbStoreByName(p.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (p *DefaultPermissionManagerImpl) Initial() {
}

func (p *DefaultPermissionManagerImpl) Checker() IPermissionChecker {
	return p.permissionChecker
}

func (p *DefaultPermissionManagerImpl) Create(category string, perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) Modify(perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) Drop(perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) Query(tenantId uint64, category string, expr PagerExpr) (
	total int64, result []Permission, error error) {
	return 0, nil, nil
}

func (p *DefaultPermissionManagerImpl) QueryByUser(tenantId, userId uint64, expr PagerExpr) (
	total int64, result []Permission, error error) {
	return 0, nil, nil
}

func (p *DefaultPermissionManagerImpl) GrantToUser(uid uint64, perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) GrantToRole(roleId uint64, perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) RevokeFromUser(uid uint64, perms ...Permission) error {
	return nil
}

func (p *DefaultPermissionManagerImpl) RevokeFromRole(roleId uint64, perms ...Permission) error {
	return nil
}

func DefaultUserManager(state *ServerState) IUserManager {
	return EmptyUserManagerImpl{
		state: state,
	}
}

// EmptyUserManagerImpl ...
type EmptyUserManagerImpl struct {
	state *ServerState
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
