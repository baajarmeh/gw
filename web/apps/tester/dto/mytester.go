package dto

import (
	"github.com/oceanho/gw/backend/gwDb"
)

type MyTester struct {
	gwDb.Model
	gwDb.HasCreationState
	gwDb.HasActivationState
	gwDb.HasModificationState
	gwDb.HasSoftDeletionState
	gwDb.HasEffectivePeriodState
	gwDb.HasTenantState
}
