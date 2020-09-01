package Dto

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/backend/gwdb"
)

type UserDto struct {
	gwdb.Model
	TenantId uint64
	Passport string
	Secret   string
	IsActive bool
	UserType gw.UserType `gorm:"-"`
}

type ProfileDto struct {
	Name string `json:"name"`
}
