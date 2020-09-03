package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type UserAccessKeySecret struct {
	gwdb.Model
	gwdb.HasTenantState
	Key    string
	Secret string
	gwdb.HasCreationState
	gwdb.HasActivationState
	gwdb.HasModificationState
}

func (UserAccessKeySecret) TableName() string {
	return getTableName("user_aks")
}
