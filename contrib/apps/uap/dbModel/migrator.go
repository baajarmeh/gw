package dbModel

import "github.com/oceanho/gw"

func Migrate(state *gw.ServerState) {
	db := state.Store().GetDbStore()
	_ = db.AutoMigrate(&User{}, &Role{}, &UserAccessKeySecret{}, &UserProfile{}, &UserRoleMapping{})
}
