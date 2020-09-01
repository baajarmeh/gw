package Db

//go:generate stringer -type=ArtifactCategory,StorageEngine

import "github.com/oceanho/gw/backend/gwdb"

const (
	storageTableName = "pvm_storage"
)

// docker/fs
type StorageEngine uint8

const (
	StorageEngineDocker     StorageEngine = 1 // docker
	StorageEngineFileSystem StorageEngine = 2 // file system
)

type Storage struct {
	gwdb.Model
	gwdb.HasTenantState
	Name         string `gorm:"type:varchar(128)"`
	Engine       StorageEngine
	Address      string `gorm:"type:varchar(512)"`
	CredentialId uint64
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Storage) TableName() string {
	return storageTableName
}
