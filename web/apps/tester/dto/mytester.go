package dto

import (
	"github.com/oceanho/gw/backend/gwDb"
)

type MyTester struct {
	gwdb.Model
	gwdb.HasCreationState
	gwdb.HasActivationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
	gwdb.HasEffectivePeriodState
	gwdb.HasTenantState
}
