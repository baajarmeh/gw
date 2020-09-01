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

type DefaultCredentialProtectServiceImpl struct {
	hash    gw.ICryptoHash
	protect gw.ICryptoProtect
}

func (dcsi DefaultCredentialProtectServiceImpl) New(
	hash gw.ICryptoHash,
	protect gw.ICryptoProtect) ICredentialProtectService {
	dcsi.hash = hash
	dcsi.protect = protect
	return dcsi
}

func (dcsi DefaultCredentialProtectServiceImpl) Encrypt(cred *Db.Credential) error {
	return nil
}

func (dcsi DefaultCredentialProtectServiceImpl) Decrypt(cred *Db.Credential) error {
	return nil
}

func (dcsi DefaultCredentialProtectServiceImpl) Validate(cred *Db.Credential) error {
	return nil
}

type DefaultCredentialServiceImpl struct {
	store                    gw.IStore
	credentialProtectService ICredentialProtectService
}

func (dcs DefaultCredentialServiceImpl) New(
	store gw.IStore,
	credentialProtectService ICredentialProtectService) ICredentialService {
	dcs.store = store
	dcs.credentialProtectService = credentialProtectService
	return dcs
}

func (dcs DefaultCredentialServiceImpl) QueryById(id uint64) (Db.Credential, error) {
	var model Db.Credential
	err := dcs.store.GetDbStore().First(&model, "id = ?", id).Error
	return model, err
}
