package Service

import (
	"github.com/oceanho/gw"
	. "github.com/oceanho/gw/contrib/apps/crond/app/Dto"
	"time"
)

type JobManager interface {
	Create(job Job) error
	Modify(job Job) error
	Query(from time.Time, to time.Time, expr gw.PagerExpr) (*JobPager, error)
	Publish(job Job) error
}

func DefaultJobManager() DefaultJobManagerImpl {
	return DefaultJobManagerImpl{}
}

type DefaultJobManagerImpl struct {
	store gw.IStore
}

func (d DefaultJobManagerImpl) Create(job Job) error {
	return nil
}

func (d DefaultJobManagerImpl) Modify(job Job) error {
	return nil
}

func (d DefaultJobManagerImpl) Query(from time.Time, to time.Time, expr gw.PagerExpr) (*JobPager, error) {
	return nil, nil
}

func (d DefaultJobManagerImpl) Publish(job Job) error {
	return nil
}

// DI interface
func (d DefaultJobManagerImpl) New(store gw.IStore) JobManager {
	d.store = store
	return d
}
