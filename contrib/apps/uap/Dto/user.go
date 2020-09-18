package Dto

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/backend/gwdb"
	"time"
)

//go:generate gomodifytags -w -file ./user.go -add-tags json -transform snakecase -all

type UserDto struct {
	gwdb.Model
	Passport string         `binding:"required" json:"passport"`
	Secret   string         `binding:"required" json:"secret"`
	IsActive bool           `json:"is_active"`
	UserType gw.UserType    `gorm:"-" json:"type"`
	Profile  UserProfileDto `gorm:"-" json:"profile"`
}

type UserProfileDto struct {
	UserID   uint64     `json:"user_id"`
	Gender   uint8      `json:"gender"` // 1.man, 2.woman, 3.custom, 4.unknown
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Avatar   string     `json:"avatar"`
	Address  string     `json:"address"`
	PostCode string     `json:"post_code"`
	BirthDay *time.Time `json:"birth_day"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}
