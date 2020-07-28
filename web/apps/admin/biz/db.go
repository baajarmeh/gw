package biz

import (
	"github.com/oceanho/gw/web/apps/admin/dto"
	"gorm.io/gorm"
)

func GetTables(db *gorm.DB, tables []dto.Table) error {
	return nil
}

func CreateTableRow(db *gorm.DB, tableName string, id interface{}, value map[string]interface{}) error {
	return db.Table(tableName).Create(value).Error
}

func UpdateTableRow(db *gorm.DB, tableName string, values map[string]interface{}) error {
	return db.Table(tableName).Updates(values).Error
}
