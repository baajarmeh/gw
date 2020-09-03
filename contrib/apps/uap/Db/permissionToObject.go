package Db

import "github.com/oceanho/gw/backend/gwdb"

type PermissionMapping struct {
	gwdb.Model
	gwdb.HasTenantState
	ObjectID     uint64         `gorm:"index:idx_tenant_expr;not null"`
	Type         PermissionType `gorm:"not null"` // 1. User Permission, 2. Role/Group Permission
	PermissionID uint64         `gorm:"index;not null"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (PermissionMapping) TableName() string {
	return getTableName("permission_to_object")
}
