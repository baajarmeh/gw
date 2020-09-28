package gwdb

import "time"

type Model struct {
	ID uint64 `gorm:"primary_key;auto_increment:true;not null" json:"id"`
}

type HasTenantState struct {
	TenantID uint64 `gorm:"default:0;not null;index:idx_tenant_expr" json:"tenant_id"`
}

type HasCreationState struct {
	CreatedAt *time.Time `gorm:"not null" json:"created_at"`
}

type HasModificationState struct {
	ModifiedAt *time.Time `json:"modified_at"`
}

type HasSoftDeletionState struct {
	IsDeleted bool       `gorm:"default:0" json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type HasLockState struct {
	IsLocked     bool       `gorm:"default:0;not null" json:"is_locked"`
	LockedAt     *time.Time `json:"locked_at"`
	LockedReason string     `gorm:"type:varchar(128)" json:"locked_reason"`
}

type HasActivationState struct {
	IsActive bool `gorm:"default:1;not null" json:"is_active"`
}

type HasEffectivePeriodState struct {
	EffectiveAt *time.Time `json:"effective_at"`
	PeriodAt    *time.Time `json:"period_at"`
}
