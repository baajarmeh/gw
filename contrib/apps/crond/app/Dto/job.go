package Dto

import (
	"github.com/oceanho/gw"
	"time"
)

type Job struct {
	ID          string
	Name        string
	Command     string
	Args        []JobArgument
	CreatedAt   time.Time
	ModifiedAt  time.Time
	ExecPlanAt  time.Time
	ExecTimeout time.Duration
}

type JobArgument struct {
	Key      string
	Value    string
	Category string
}

type JobPager struct {
	gw.PagerResult
	Data []Job `json:"data"`
}
