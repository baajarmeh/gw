package Db

import "github.com/oceanho/gw/backend/gwdb"

type Permission struct {
	gwdb.Model
	TenantID   uint64 `gorm:"UNIQUEINDEX:idx_permission_tenant_id_app_id;not null"`
	AppID      uint64 `gorm:"UNIQUEINDEX:idx_permission_tenant_id_app_id;not null"`
	Key        string `gorm:"type:varchar(64);UNIQUEINDEX:idx_permission_tenant_id_app_id;not null"`
	Name       string `gorm:"type:varchar(128);not null"`
	Descriptor string `gorm:"type:varchar(256)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Permission) TableName() string {
	return getTableName("perm")
}

type PermissionType uint8

const (
	UserPermission PermissionType = 1
	RolePermission PermissionType = 2
)