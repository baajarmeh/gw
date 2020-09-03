package Db

//go:generate stringer -type=ArtifactCategory,StorageEngine

import "github.com/oceanho/gw/backend/gwdb"

const artifactTableName = "gw_pvm_artifact"

type ArtifactCategory uint8

const (
	GzipFile    ArtifactCategory = iota
	GenericFile ArtifactCategory = 1
	DockerImage ArtifactCategory = 2
)

type Artifact struct {
	gwdb.Model
	gwdb.HasTenantState
	ComponentId uint64
	Location    string
	Category    ArtifactCategory
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Artifact) TableName() string {
	return artifactTableName
}
