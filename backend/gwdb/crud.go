package gwdb

import "gorm.io/gorm"

func Create(db *gorm.DB, value interface{}) error {
	return db.Create(value).Error
}

func Update(db *gorm.DB, values interface{}) error {
	return db.Updates(values).Error
}

func Delete(db *gorm.DB, value interface{}) error {
	return db.Delete(value).Error
}

func Get(db *gorm.DB, out interface{}, id uint64) error {
	return db.Where("id = ?", id).First(out).Error
}

func Query(db *gorm.DB, out interface{}, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Find(out).Error
}

func QueryByOrder(db *gorm.DB, out interface{}, order interface{}, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Order(order).Find(out).Error
}

func QueryList(db *gorm.DB, limit, offset int, out interface{}, total int64, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Count(&total).Offset(offset).Limit(limit).Order("id desc").Find(out).Error
}

func QueryListByOrder(db *gorm.DB, limit, offset int, out interface{}, total int64, order interface{}, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Count(&total).Offset(offset).Limit(limit).Order(order).Find(out).Error
}
