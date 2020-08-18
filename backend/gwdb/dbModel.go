package gwdb

import "time"

type Model struct {
	ID uint64 `gorm:"primary_key;auto_increment:true;not null"`
}

type HasTenantState struct {
	TenantId uint64 `gorm:"default:0;not null;index:idx_tenant_expr"`
}

type HasCreationState struct {
	CreatedAt *time.Time `gorm:"not null"`
}

type HasModificationState struct {
	ModifiedAt *time.Time
}

type HasSoftDeletionState struct {
	IsDeleted *bool `gorm:"default:0"`
	DeletedAt *time.Time
}

type HasLockState struct {
	IsLocked     *bool `gorm:"default:0;not null"`
	LockedAt     *time.Time
	LockedReason string `gorm:"type:varchar(128)"`
}

type HasActivationState struct {
	IsActive *bool `gorm:"default:1;not null"`
}

type HasEffectivePeriodState struct {
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}
