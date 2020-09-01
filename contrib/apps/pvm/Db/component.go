package Db

import "github.com/oceanho/gw/backend/gwdb"

const (
	componentTableName = "pvm_component"
)

type Component struct {
	gwdb.Model
	gwdb.HasTenantState
	Name       string `gorm:"type:varchar(256)"`
	Descriptor string `gorm:"type:varchar(512)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Component) TableName() string {
	return componentTableName
}
