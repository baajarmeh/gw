package Db

import "github.com/oceanho/gw/backend/gwdb"

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
	Type         PermissionType `gorm:"not null"` // 1. User Permission, 2. Role/Group Permission
	ObjectID     uint64         `gorm:"index:idx_tenant_expr; not null"`
	PermissionID uint64         `gorm:"not null"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (ObjectPermission) TableName() string {
	return getTableName("object_permission")
}
