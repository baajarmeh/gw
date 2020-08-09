package biz

import (
	"github.com/oceanho/gw/contrib/apps/tester/dto"
	"gorm.io/gorm"
)

func CreateMyTester(db *gorm.DB, dto *dto.MyTester) error {
	return db.Create(dto).Error
}

func QueryMyTester(db *gorm.DB, outs *[]dto.MyTester) error {
	return db.Find(outs).Error
}
