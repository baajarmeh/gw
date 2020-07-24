package biz

import (
	"github.com/oceanho/gw/web/apps/uap/entities"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *entities.User) error {
	return db.Create(user).Error
}
