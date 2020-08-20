package restAPIs

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/dbModel"
	"gorm.io/gorm"
)

type User struct {
}

func (u User) Name() string {
	return "user"
}

//
// APIs
//
func (u User) Get(ctx *gw.Context) {
	store := ctx.State.Store()
	db := store.GetDbStore()
	var user dbModel.User
	err := db.First(&user)
	ctx.JSON(err, user)
}

func (u User) OnGetBefore() gw.Decorator {
	return gw.NewStoreDbSetupDecorator(func(ctx gw.Context, db *gorm.DB) *gorm.DB {
		return db
	})
}

func (u User) Post(ctx *gw.Context) {

}

func (u User) Put(ctx *gw.Context) {

}

func (u User) Delete(ctx *gw.Context) {

}

func (u User) QueryList(ctx *gw.Context) {

}
