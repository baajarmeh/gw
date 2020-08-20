package dbModel

import "github.com/oceanho/gw"

func Migrate(ctx gw.MigrationContext) {
	db := ctx.Store().GetDbStore()
	_ = db.AutoMigrate(&User{}, &Role{}, &UserAccessKeySecret{}, &UserProfile{}, &UserRoleMapping{})
}
