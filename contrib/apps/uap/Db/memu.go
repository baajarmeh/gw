package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type Menu struct {
	gwdb.Model
	AppID            uint64
	ParentID         uint64
	Name             string `gorm:"type:varchar(32)"`
	Icon             string `gorm:"type:varchar(64)"`
	Link             string `gorm:"type:varchar(256)"`
	LinkOpenBehavior string `gorm:"type:varchar(16)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
}

func (Menu) TableName() string {
	return getTableName("menu")
}
