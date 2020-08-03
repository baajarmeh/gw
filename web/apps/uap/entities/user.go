package entities

import (
	"fmt"
	"github.com/oceanho/gw/backend/gwDb"
	"time"
)

var tablePrefix = "uap"

type User struct {
	gwDb.Model
	Passport  string `gorm:"type:varchar(32);unique_index;not null"`
	Secret    string `gorm:"type:varchar(128);not null"`
	IsTenancy bool   `gorm:"default:0;not null"`
	IsAdmin   bool   `gorm:"default:0;not null"`
	gwDb.HasTenantState
	gwDb.HasLockState
	gwDb.HasCreationState
	gwDb.HasModificationState
	gwDb.HasSoftDeletionState
	gwDb.HasActivationState
}

func (User) TableName() string {
	return fmt.Sprintf("%s_%s", tablePrefix, "user")
}

type Profile struct {
	gwDb.Model
	Gender   uint8  `gorm:"default:1"` // 1.man, 2.woman, 3.custom, 4.unknown
	UserID   uint64 `gorm:"index"`
	Name     string `gorm:"type:varchar(64);index"`
	Email    string `gorm:"type:varchar(128);index"`
	Phone    string `gorm:"type:varchar(16);index"`
	Avatar   string `gorm:"type:varchar(256)"`
	Address  string `gorm:"type:varchar(256)"`
	PostCode string `gorm:"type:varchar(16)"`
	BirthDay *time.Time
	gwDb.HasTenantState
	gwDb.HasCreationState
	gwDb.HasModificationState
}

func (Profile) TableName() string {
	return fmt.Sprintf("%s_%s", tablePrefix, "profile")
}
