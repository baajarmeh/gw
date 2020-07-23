package dto

import "github.com/oceanho/gw/backend"

type MyTester struct {
	backend.Model
	backend.HasCreationState
	backend.HasActivationState
	backend.HasModificationState
	backend.HasSoftDeletionState
	backend.HasEffectivePeriodState
	backend.HasTenantState
}
