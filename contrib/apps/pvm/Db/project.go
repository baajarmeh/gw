package Db

import "github.com/oceanho/gw/backend/gwdb"

const (
	projectTableName          = "pvm_project"
	projectVersionTableName   = "pvm_project_version"
	projectComponentTableName = "pvm_project_component"
)

type Project struct {
	gwdb.Model
	gwdb.HasTenantState
	Name       string `gorm:"type:varchar(256)"`
	Descriptor string `gorm:"type:varchar(512)"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Project) TableName() string {
	return projectTableName
}

type ProjectComponent struct {
	gwdb.Model
	ProjectId   uint64
	ComponentId uint64
	gwdb.HasCreationState
}

func (ProjectComponent) TableName() string {
	return projectComponentTableName
}

type ProjectVersion struct {
	gwdb.Model
	ProjectId uint64
	Version   string
	Remarks   string
	Publisher uint64
	gwdb.HasCreationState
}

func (ProjectVersion) TableName() string {
	return projectVersionTableName
}
