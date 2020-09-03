package Db

import "github.com/oceanho/gw/backend/gwdb"

type UserRole struct {
	gwdb.Model
	gwdb.HasTenantState
	UserID uint64
	RoleId uint64
	gwdb.HasCreationState
}

func (UserRole) TableName() string {
	return getTableName("user_roles")
}
