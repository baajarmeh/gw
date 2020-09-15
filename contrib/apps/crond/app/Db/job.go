package Db

import (
	"github.com/oceanho/gw/backend/gwdb"
	"time"
)

type Job struct {
	gwdb.Model
	Name       string
	Command    string
	CreatedAt  *time.Time
	ModifiedAt *time.Time
	ExecPlan   string
	Timeout    int
}

type JobArgument struct {
	gwdb.Model
	JobID    uint64
	Key      string `gorm:"type:varchar(128)"`
	Value    string `gorm:"type:varchar(256)"`
	Category string `gorm:"type:varchar(256)"`
}
