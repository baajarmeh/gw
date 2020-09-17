package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type UserAccessKeySecret struct {
	gwdb.Model
	UserID uint64
	Key    string `gorm:"type:varchar(64);unique;not null"`
	Secret string `gorm:"type:varchar(128);not null"`
	gwdb.HasCreationState
	gwdb.HasActivationState
	gwdb.HasModificationState
}

func (UserAccessKeySecret) TableName() string {
	return getTableName("user_aks")
}
