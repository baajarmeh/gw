package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
)

type CredentialCategory uint8

const (
	SshKeyCredential    CredentialCategory = 1 // ssh pub/key
	BasicAuthCredential CredentialCategory = 2 // basic auth
	DrCredential        CredentialCategory = 3 // docker registry
	TLSCertCredential   CredentialCategory = 4 // TLS Certificate
)

type Credential struct {
	gwdb.Model
	gwdb.HasTenantState
	UserId    uint64 `gorm:"default:0;not null"`
	Name      string `gorm:"type:varchar(32);"`
	Value     string
	Signature string
	Category  CredentialCategory
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
}

func (Credential) TableName() string {
	return getTableName("credential")
}
