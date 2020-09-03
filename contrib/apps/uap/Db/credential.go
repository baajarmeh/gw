package Db

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/backend/gwdb"
	"gorm.io/gorm"
)

type CredentialCategory uint8

const (
	SshKeyCredential    CredentialCategory = 1
	BasicAuthCredential CredentialCategory = 2
	DrCredential        CredentialCategory = 3
	TLSCertCredential   CredentialCategory = 4
	DatabaseCredential  CredentialCategory = 5
	RedisCredential     CredentialCategory = 6
	MongoDBCredential   CredentialCategory = 7
)

type Credential struct {
	gwdb.Model
	gwdb.HasTenantState
	UserID    uint64 `gorm:"not null"`
	Key       string `gorm:"type:varchar(32);unique;not null"`
	Name      string `gorm:"type:varchar(32);"`
	Value     string
	Signature string
	Category  CredentialCategory
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
}

func (cred *Credential) BeforeCreate(tx *gorm.DB) error {
	ctx, ok := gw.GetContextFromDB(tx)
	if ok {
		cred.Key = ctx.Server().IDGenerator.NewStrID(32)
	}
	return nil
}

func (Credential) TableName() string {
	return getTableName("credential")
}
