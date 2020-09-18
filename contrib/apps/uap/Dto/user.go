package Dto

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/backend/gwdb"
)

type UserDto struct {
	gwdb.Model
	Passport string `binding:"required"`
	Secret   string `binding:"required"`
	IsActive bool
	UserType gw.UserType `gorm:"-"`
}

type ProfileDto struct {
	Name string `json:"name"`
}
