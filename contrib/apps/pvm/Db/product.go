package Db

import "github.com/oceanho/gw/backend/gwdb"

const productTableName = "gw_pvm_product"

type Product struct {
	gwdb.Model
	gwdb.HasTenantState
	Name       string `gorm:"type:varchar(256)"`
	Descriptor string `gorm:"type:varchar(512)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Product) TableName() string {
	return productTableName
}
