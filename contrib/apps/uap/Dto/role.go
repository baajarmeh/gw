package Dto

import (
	"github.com/oceanho/gw/backend/gwdb"
)

//go:generate gomodifytags -w -file ./role.go -add-tags json -transform snakecase -all

type Role struct {
	gwdb.Model
	Name       string `binding:"required,lt=128" json:"name" form:"name"`
	Descriptor string `binding:"lt=256" json:"desc" form:"desc"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}
