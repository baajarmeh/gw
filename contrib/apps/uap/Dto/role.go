package Dto

import (
	"github.com/oceanho/gw/backend/gwdb"
)

//go:generate gomodifytags -w -file ./role.go -add-tags json -transform snakecase -all

type Role struct {
	gwdb.Model
	Name       string `gorm:"type:varchar(32);not null" json:"name"`
	Descriptor string `gorm:"type:varchar(128);not null" json:"descriptor"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}
