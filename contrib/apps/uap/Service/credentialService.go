package Service

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Db"
)

type ICredentialProtectService interface {
	Encrypt(cred *Db.Credential) error
	Decrypt(cred *Db.Credential) error
	Validate(cred *Db.Credential) error
}

type ICredentialService interface {
	QueryById(id uint64) (Db.Credential, error)
}

type DefaultCredentialProtectService struct {
	hash    gw.ICryptoHash
	protect gw.ICryptoProtect
}

func (dcsi DefaultCredentialProtectService) New(
	hash gw.ICryptoHash,
	protect gw.ICryptoProtect) ICredentialProtectService {
	dcsi.hash = hash
	dcsi.protect = protect
	return dcsi
}

func (dcsi DefaultCredentialProtectService) Encrypt(cred *Db.Credential) error {
	return nil
}

func (dcsi DefaultCredentialProtectService) Decrypt(cred *Db.Credential) error {
	return nil
}

func (dcsi DefaultCredentialProtectService) Validate(cred *Db.Credential) error {
	return nil
}

type DefaultCredentialService struct {
	store                    gw.IStore
	credentialProtectService ICredentialProtectService
}

// DI
func (dcs DefaultCredentialService) New(store gw.IStore,
	credentialProtectService ICredentialProtectService) ICredentialService {
	dcs.store = store
	dcs.credentialProtectService = credentialProtectService
	return dcs
}

func (dcs DefaultCredentialService) QueryById(id uint64) (Db.Credential, error) {
	var model Db.Credential
	err := dcs.store.GetDbStore().First(&model, "id = ?", id).Error
	return model, err
}
