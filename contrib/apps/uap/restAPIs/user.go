package restAPIs

import "github.com/oceanho/gw"

type User struct {
}

func (u User) Name() string {
	return "user"
}

func UserDynamicRestAPI() gw.IDynamicRestAPI {
	return &User{}
}

//
// APIs
//
func (u User) Get(ctx *gw.Context) {
	ctx.State.Store()
}

func (u User) Post(ctx *gw.Context) {

}

func (u User) Put(ctx *gw.Context) {

}

func (u User) Delete(ctx *gw.Context) {

}

func (u User) QueryList(ctx *gw.Context) {

}
