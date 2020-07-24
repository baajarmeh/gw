package backend

import "time"

type Model struct {
	ID uint64 `gorm:"primary_key;auto_increment:true;NOT NULL"`
}

type HasTenantState struct {
	TenantId uint64
}

type HasCreationState struct {
	CreatedAt *time.Time
}

type HasModificationState struct {
	ModifiedAt *time.Time
}

type HasSoftDeletionState struct {
	IsDeleted *bool
	DeletedAt *time.Time
}

type HasLockState struct {
	IsLocked   *bool
	LockedAt   *time.Time
	LockReason string
}

type HasActivationState struct {
	IsActive *bool
}

type HasEffectivePeriodState struct {
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}
