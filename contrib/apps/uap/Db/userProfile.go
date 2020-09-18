package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
	"time"
)

type UserProfile struct {
	gwdb.Model
	UserID   uint64 `gorm:"index;not null"`
	Gender   uint8  `gorm:"default:4"` // 1.man, 2.woman, 3.custom, 4.unknown
	Name     string `gorm:"type:varchar(64)"`
	Email    string `gorm:"type:varchar(128);index"`
	Phone    string `gorm:"type:varchar(16);index"`
	Avatar   string `gorm:"type:varchar(256)"`
	Address  string `gorm:"type:varchar(256)"`
	WebSite  string `gorm:"type:varchar(256)"`
	PostCode string `gorm:"type:varchar(16)"`
	BirthDay *time.Time
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (UserProfile) TableName() string {
	return getTableName("user_profile")
}
