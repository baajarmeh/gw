package db

import "github.com/oceanho/gw/contrib/backend/models"

type Tenant struct {
	models.Model
	models.HasCreationState
	models.HasModificationState
	models.HasSoftDeletionState
	models.HasActivationState
	Name string
}

func (Tenant) TableName() string {
	return "uap_tenant"
}
