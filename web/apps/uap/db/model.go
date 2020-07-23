package db

import (
	"github.com/oceanho/gw/backend"
)

type Tenant struct {
	backend.Model
	backend.HasCreationState
	backend.HasModificationState
	backend.HasSoftDeletionState
	backend.HasActivationState
	Name string
}

func (Tenant) TableName() string {
	return "uap_tenant"
}
