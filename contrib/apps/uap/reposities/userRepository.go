package reposities

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(model *dbModel.User) error
}

type UserRepositoryImpl struct {
	Db          *gorm.DB
	User        *gorm.DB
	UserProfile *gorm.DB
}

func (u UserRepositoryImpl) New(store gw.IStore) UserRepository {
	u.Db = store.GetDbStore()
	// FIXME(oceanho): What's happened?
	//  If call u.Db.Model(xxx)
	//  the u.Db.Create(model).Error of model.name are dbModel.UserProfile{}.TableName() ?
	//u.User = u.Db.Model(dbModel.User{})
	//u.UserProfile = u.Db.Model(dbModel.UserProfile{})
	return u
}

func (u UserRepositoryImpl) Create(model *dbModel.User) error {
	return u.Db.Create(model).Error
}
