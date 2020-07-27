package entities

import (
	"fmt"
	"github.com/oceanho/gw/backend"
	"time"
)

var tablePrefix = "uap"

type User struct {
	backend.Model
	Passport  string `gorm:"type:varchar(32);unique_index;not null"`
	Secret    string `gorm:"type:varchar(128)"`
	IsTenancy bool   `gorm:"default:0;not null"`
	backend.HasTenantState
	backend.HasLockState
	backend.HasCreationState
	backend.HasModificationState
	backend.HasSoftDeletionState
	backend.HasActivationState
}

func (User) TableName() string {
	return fmt.Sprintf("%s_%s", tablePrefix, "user")
}

type Profile struct {
	backend.Model
	Gender   uint8  `gorm:"default:1"` // 1.man, 2.woman, 3.custom, 4.unknown
	UserID   uint64 `gorm:"index"`
	Name     string `gorm:"type:varchar(64);index"`
	Email    string `gorm:"type:varchar(128);index"`
	Phone    string `gorm:"type:varchar(16);index"`
	Avatar   string `gorm:"type:varchar(256)"`
	Address  string `gorm:"type:varchar(256)"`
	PostCode string `gorm:"type:varchar(16)"`
	BirthDay *time.Time
	backend.HasTenantState
	backend.HasCreationState
	backend.HasModificationState
}

func (Profile) TableName() string {
	return fmt.Sprintf("%s_%s", tablePrefix, "profile")
}
