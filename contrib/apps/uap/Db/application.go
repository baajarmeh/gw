package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type App struct {
	gwdb.Model
	Key        string `gorm:"type:varchar(32);not null"`
	Name       string `gorm:"type:varchar(128);not null"`
	Router     string `gorm:"type:varchar(128);not null"`
	Descriptor string `gorm:"type:varchar(512);not null"`
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
}

func (App) TableName() string {
	return getTableName("application")
}
