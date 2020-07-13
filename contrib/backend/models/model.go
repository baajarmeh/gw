package models

import "time"

type Model struct {
	ID uint64 `gorm:"primary_key;auto_increment:true"`
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

type HasActivationState struct {
	IsActive *bool
}
