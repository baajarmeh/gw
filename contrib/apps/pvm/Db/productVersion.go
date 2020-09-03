package Db

import "github.com/oceanho/gw/backend/gwdb"

const productVersionTableName = "gw_pvm_product_version"

type ProductVersion struct {
	gwdb.Model
	gwdb.HasTenantState
	Name       string `gorm:"type:varchar(256)"`
	Descriptor string `gorm:"type:varchar(512)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (ProductVersion) TableName() string {
	return productVersionTableName
}
