package GwImpl

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
	"gorm.io/gorm"
	"sync"
)

type PermissionManager struct {
	_state            int
	locker            sync.Mutex
	store             gw.IStore
	state             *gw.ServerState
	permissionChecker gw.IPermissionChecker
	queryPermSQL      string
	delPermMapSQL     string
	queryPermMapSQL   string
	conf              *conf.ApplicationConfig
}

// DI

func (pm PermissionManager) New(store gw.IStore) gw.IPermissionManager {
	pm.store = store
	return &pm
}

func (pm *PermissionManager) Store() *gorm.DB {
	return pm.store.GetDbStoreByName(pm.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (pm *PermissionManager) Checker() gw.IPermissionChecker {
	return pm.permissionChecker
}

func (pm *PermissionManager) Create(perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	pm.locker.Lock()
	defer pm.locker.Unlock()
	db := pm.Store()
	tx := db.Begin()
	var registered = make(map[string]bool)
	for i := 0; i < len(perms); i++ {
		p := perms[i]
		var model Db.Permission
		uk := fmt.Sprintf("%d-%d-%s", p.TenantID, p.AppID, p.Key)
		if registered[uk] {
			continue
		}
		err := db.Model(Db.Permission{}).Where(pm.queryPermMapSQL, p.TenantID, p.AppID, p.Key).Take(&model).Error
		if err != nil && err.Error() != "record not found" {
			panic(fmt.Sprintf("check perms has exist fail, err: %v", err))
		}
		if model.ID == 0 {
			model.Key = p.Key
			model.Name = p.Name
			model.TenantID = p.TenantID
			model.AppID = p.AppID
			model.Scope = p.Scope
			model.Descriptor = p.Descriptor
			_ = tx.Create(&model)
			perms[i].ID = model.ID
		}
		// FIXME(OceanHo): needs?
		// tx.Updates(model)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) Modify(perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	pm.locker.Lock()
	defer pm.locker.Unlock()
	tx := pm.Store().Begin()
	for _, p := range perms {
		tx.Updates(p)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) Drop(perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		tx.Delete(p)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) Query(tenantId, appId uint64, expr gw.PagerExpr) (
	total int64, result []*gw.Permission, error error) {
	error = pm.Store().Where("tenant_id = ? and app_id = ?",
		tenantId, appId).Count(&total).Offset(expr.PageOffset()).Limit(expr.PageSize).Scan(result).Error
	return
}

func (pm *PermissionManager) QueryByUser(tenantId,
	userId uint64, expr gw.PagerExpr) (total int64, result []*gw.Permission, error error) {
	var sql = fmt.Sprintf(pm.queryPermSQL, tenantId, userId, gw.UserPermission)
	var countSql = fmt.Sprintf("select t1.id as total %s", sql)
	var dataSql = fmt.Sprintf("select t1.* %s limit %d offset %d", sql, expr.PageSize, expr.PageOffset())
	tx := pm.Store().Begin()
	err := tx.Raw(countSql).Scan(&total).Error
	if err != nil {
		tx.Rollback()
		return total, nil, err
	}
	err = tx.Raw(dataSql).Scan(&result).Error
	if err != nil {
		tx.Rollback()
		return 0, nil, err
	}
	err = tx.Commit().Error
	return total, result, err
}

func (pm *PermissionManager) GrantToUser(uid uint64, perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		var pm Db.PermissionRelation
		pm.PermissionID = p.ID
		pm.TenantID = p.TenantID
		pm.Type = gw.UserPermission
		pm.ObjectID = uid
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) GrantToRole(roleId uint64, perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		var pm Db.PermissionRelation
		pm.PermissionID = p.ID
		pm.TenantID = p.TenantID
		pm.Type = gw.RolePermission
		pm.ObjectID = roleId
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) RevokeFromUser(uid uint64, perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		tx.Delete(Db.PermissionRelation{}, pm.delPermMapSQL, uid, p.TenantID, p.ID, gw.UserPermission)
	}
	return tx.Commit().Error
}

func (pm *PermissionManager) RevokeFromRole(roleId uint64, perms ...*gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		tx.Delete(Db.PermissionRelation{}, pm.delPermMapSQL, roleId, p.TenantID, p.ID, gw.RolePermission)
	}
	return tx.Commit().Error
}

func DefaultPermissionManager(state *gw.ServerState) gw.IPermissionManager {
	cnf := state.ApplicationConfig()
	ptn := Db.Permission{}.TableName()
	pmtn := Db.PermissionRelation{}.TableName()
	queryPermSQL := fmt.Sprintf(" from %s t1 inner join %s t2 on t1.id = t2.permission_id", ptn, pmtn)
	queryPermSQL = queryPermSQL + " where t2.tenant_id=%d and t2.object_id=%d and t2.type=%d"
	return &PermissionManager{
		conf:              cnf,
		state:             state,
		store:             state.Store(),
		queryPermSQL:      queryPermSQL,
		delPermMapSQL:     " object_id = ? and tenant_id = ? and permission_id = ? and type = ? ",
		queryPermMapSQL:   " tenant_id = ? and app_id = ? and `key` = ? ",
		permissionChecker: state.PermissionChecker(),
	}
}