package Db

import (
	"fmt"
	"github.com/oceanho/gw/backend/gwdb"
)

var tablePrefix = "gw_uap"

func getTableName(name string) string {
	return fmt.Sprintf("%s_%s", tablePrefix, name)
}

type User struct {
	gwdb.Model
	TenantId  uint64 `gorm:"default:0;not null;UNIQUEINDEX:idx_tenant_id_passport"`
	Passport  string `gorm:"type:varchar(32);UNIQUEINDEX:idx_tenant_id_passport;not null"`
	Secret    string `gorm:"type:varchar(128);not null"`
	IsUser    bool   `gorm:"default:0;not null"`
	IsAdmin   bool   `gorm:"default:0;not null"`
	IsTenancy bool   `gorm:"default:0;not null"`
	gwdb.HasLockState
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
	gwdb.HasActivationState
}

func (User) TableName() string {
	return getTableName("user")
}
