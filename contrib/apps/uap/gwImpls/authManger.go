package gwImpls

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"time"
)

type UserManager struct {
	userCachePrefix  string
	cacheStoreName   string
	backendStoreName string
	expiration       time.Duration
	permPagerExpr    gw.PagerExpr
	state            gw.ServerState
}

var ErrUserHasExists = fmt.Errorf("user object has exists")

func (u UserManager) Create(user *gw.User) error {
	store := u.state.Store()
	db := store.GetDbStore()
	var model dbModel.User
	err := db.First(&model, "tenant_id=? and passport=?", user.TenantId, user.Passport).Error
	if err != nil && err.Error() != "record not found" {
		return err
	}
	if model.ID > 0 {
		return ErrUserHasExists
	}
	model.Passport = user.Passport
	model.TenantId = user.TenantId
	model.Secret = user.SecretHash
	model.IsAdmin = user.RoleId == 1
	model.IsTenancy = user.RoleId == 2
	model.IsUser = user.RoleId >= 10000
	tx := store.GetDbStore().Begin()
	err = tx.Create(&model).Error
	if err != nil {
		return err
	}
	err = tx.Commit().Error
	user.Id = model.ID
	return err
}

func (u UserManager) Modify(user gw.User) error {
	panic("implement me")
}

func (u UserManager) Delete(tenantId, userId uint64) error {
	panic("implement me")
}

func (u UserManager) Query(tenantId, userId uint64) (gw.User, error) {
	panic("implement me")
}

func (u UserManager) QueryList(tenantId uint64, expr gw.PagerExpr, total int64, out []gw.User) error {
	panic("implement me")
}

func DefaultUserManager(state gw.ServerState) UserManager {
	return UserManager{
		state:            state,
		cacheStoreName:   "primary",
		backendStoreName: "primary",
		expiration:       time.Hour * 168, // One week.
	}
}
