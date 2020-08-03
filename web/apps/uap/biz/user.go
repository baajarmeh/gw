package biz

import (
	"github.com/oceanho/gw/web/apps/uap/dbModel"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *dbModel.User) error {
	return db.Create(user).Error
}
