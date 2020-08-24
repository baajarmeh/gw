package gwImpls

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"gorm.io/gorm"
	"sync"
)

type PermissionManagerImpl struct {
	_state            int
	state             gw.ServerState
	store             gw.IStore
	conf              conf.ApplicationConfig
	locker            sync.Mutex
	queryPermSQL      string
	delPermMapSQL     string
	queryPermMapSQL   string
	permissionChecker gw.IPermissionChecker
}

func (pm *PermissionManagerImpl) Store() *gorm.DB {
	return pm.store.GetDbStoreByName(pm.conf.Security.Auth.Permission.DefaultStore.Name)
}

func (pm *PermissionManagerImpl) Initial() {
	pm.locker.Lock()
	defer pm.locker.Unlock()
	if pm._state > 0 {
		return
	}
	store := pm.Store()
	var tables []interface{}
	tables = append(tables, &dbModel.Permission{})
	tables = append(tables, &dbModel.PermissionMapping{})
	err := store.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}

func (pm *PermissionManagerImpl) Checker() gw.IPermissionChecker {
	return pm.permissionChecker
}

func (pm *PermissionManagerImpl) Create(category string, perms ...gw.Permission) error {
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
		var model dbModel.Permission
		uk := fmt.Sprintf("%d-%s-%s", p.TenantId, p.Category, p.Key)
		if registered[uk] {
			continue
		}
		if p.Category == "" {
			p.Category = category
		}
		err := db.Model(dbModel.Permission{}).Where(pm.queryPermMapSQL, p.TenantId, p.Category, p.Key).Take(&model).Error
		if err != nil && err.Error() != "record not found" {
			panic(fmt.Sprintf("check perms has exist fail, err: %v", err))
		}
		if model.ID == 0 {
			model.Key = p.Key
			model.Name = p.Name
			model.TenantId = p.TenantId
			model.Category = p.Category
			model.Descriptor = p.Descriptor
			tx.Create(&model)
			perms[i].ID = model.ID
		}
		// FIXME(OceanHo): needs?
		// tx.Updates(model)
	}
	return tx.Commit().Error
}

func (pm *PermissionManagerImpl) Modify(perms ...gw.Permission) error {
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

func (pm *PermissionManagerImpl) Drop(perms ...gw.Permission) error {
	if len(perms) < 1 {
		return gw.ErrorEmptyInput
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		tx.Delete(p)
	}
	return tx.Commit().Error
}

func (pm *PermissionManagerImpl) Query(tenantId uint64, category string, expr gw.PagerExpr) (
	total int64, result []gw.Permission, error error) {
	error = pm.Store().Where("tenant_id = ? and category = ?",
		tenantId, category).Count(&total).Offset(expr.PageOffset()).Limit(expr.PageSize).Scan(result).Error
	return
}

func (pm *PermissionManagerImpl) QueryByUser(tenantId,
	userId uint64, expr gw.PagerExpr) (total int64, result []gw.Permission, error error) {
	var sql = fmt.Sprintf(pm.queryPermSQL, tenantId, userId, dbModel.UserPermission)
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

func (pm *PermissionManagerImpl) GrantToUser(uid uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		var pm dbModel.PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = dbModel.UserPermission
		pm.ObjectID = uid
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (pm *PermissionManagerImpl) GrantToRole(roleId uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		var pm dbModel.PermissionMapping
		pm.PermissionID = p.ID
		pm.TenantId = p.TenantId
		pm.Type = dbModel.RolePermission
		pm.ObjectID = roleId
		tx.Create(&pm)
	}
	return tx.Commit().Error
}

func (pm *PermissionManagerImpl) RevokeFromUser(uid uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		tx.Delete(dbModel.PermissionMapping{}, pm.delPermMapSQL, uid, p.TenantId, p.ID, dbModel.UserPermission)
	}
	return tx.Commit().Error
}

func (pm *PermissionManagerImpl) RevokeFromRole(roleId uint64, perms ...gw.Permission) error {
	if len(perms) < 1 {
		return nil
	}
	tx := pm.Store().Begin()
	for _, p := range perms {
		p := p
		tx.Delete(dbModel.PermissionMapping{}, pm.delPermMapSQL, roleId, p.TenantId, p.ID, dbModel.RolePermission)
	}
	return tx.Commit().Error
}

func DefaultPermissionManager(state gw.ServerState) gw.IPermissionManager {
	cnf := state.ApplicationConfig()
	ptn := dbModel.Permission{}.TableName()
	pmtn := dbModel.PermissionMapping{}.TableName()
	queryPermSQL := fmt.Sprintf(" from %s t1 inner join %s t2 on t1.id = t2.permission_id", ptn, pmtn)
	queryPermSQL = queryPermSQL + " where t2.tenant_id=%d and t2.object_id=%d and t2.type=%d"
	return &PermissionManagerImpl{
		conf:            *cnf,
		state:           state,
		store:           state.Store(),
		queryPermSQL:    queryPermSQL,
		delPermMapSQL:   " object_id = ? and tenant_id = ? and permission_id = ? and type = ? ",
		queryPermMapSQL: " tenant_id = ? and category = ? and `key` = ? ",
		permissionChecker: gw.DefaultPassPermissionChecker{
			State:           state,
			CustomCheckFunc: nil,
		},
	}
}
