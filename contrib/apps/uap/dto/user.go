package dto

import "github.com/oceanho/gw/backend/gwdb"

type UserDto struct {
	gwdb.Model
	TenantId uint64
	Passport string
	Secret   string
	UserRole uint8 `gorm:"-"`
}

type ProfileDto struct {
	Name string `json:"name"`
}
