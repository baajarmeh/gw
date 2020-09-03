package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type User struct {
	gwdb.Model
	TenantID    uint64 `gorm:"default:0;not null;UNIQUEINDEX:idx_tenant_id_passport"`
	Passport    string `gorm:"type:varchar(32);UNIQUEINDEX:idx_tenant_id_passport;not null"`
	Secret      string `gorm:"type:varchar(128);not null"`
	IsUser      bool   `gorm:"default:0;not null"`
	IsAdmin     bool   `gorm:"default:0;not null"`
	IsTenancy   bool   `gorm:"default:0;not null"`
	LimitSuffix string `gorm:"type:varchar(32);not null;default:''"`
	gwdb.HasLockState
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
	gwdb.HasActivationState
}

func (User) TableName() string {
	return getTableName("user")
}
