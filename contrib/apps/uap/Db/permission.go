package Db

import "github.com/oceanho/gw/backend/gwdb"

type Permission struct {
	gwdb.Model
	TenantId   uint64 `gorm:"UNIQUEINDEX:idx_permission_tenant_id_app_id;not null"`
	AppId      uint64 `gorm:"UNIQUEINDEX:idx_permission_tenant_id_app_id;not null"`
	Module     string `gorm:"type:varchar(64);not null"`
	Resource   string `gorm:"type:varchar(64);not null"`
	Category   string `gorm:"type:varchar(128);not null"`
	Key        string `gorm:"type:varchar(64);not null"`
	Name       string `gorm:"type:varchar(128);not null"`
	Descriptor string `gorm:"type:varchar(256)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Permission) TableName() string {
	return getTableName("permission")
}

type PermissionType uint8

const (
	UserPermission PermissionType = 1
	RolePermission PermissionType = 2
)

type ObjectPermission struct {
	gwdb.Model
	gwdb.HasTenantState
	ObjectID     uint64         `gorm:"index:idx_tenant_expr;not null"`
	Type         PermissionType `gorm:"not null"` // 1. User Permission, 2. Role/Group Permission
	PermissionID uint64         `gorm:"index;not null"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (ObjectPermission) TableName() string {
	return getTableName("object_permission")
}
