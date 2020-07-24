package entities

import (
	"fmt"
	"github.com/oceanho/gw/backend"
	"time"
)

var tablePrefix = "uap"

type User struct {
	backend.Model
	Passport  string `gorm:"type:varchar(128);unique_index"`
	Secret    string `gorm:"type:varchar(128)"`
	IsTenancy bool
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
	Gender   uint8
	UserID   uint64 `gorm:"INDEX"`
	Name     string `gorm:"type:varchar(64);INDEX"`
	Email    string `gorm:"type:varchar(128);INDEX"`
	Phone    string `gorm:"type:varchar(16);INDEX"`
	Avatar   string `gorm:"type:varchar(256);INDEX"`
	CardNo   string `gorm:"type:varchar(18);INDEX"`
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
