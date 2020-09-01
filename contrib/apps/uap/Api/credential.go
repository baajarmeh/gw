package Api

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Service"
	"gorm.io/gorm"
)

func QueryCredentialById(ctx *gw.Context) {
	var id uint64
	if ctx.MustGetIdUint64FromParam(&id) != nil {
		return
	}
	var services = Service.Services(ctx)
	_ = services
	model, err := services.CredentialService.QueryById(id)
	ctx.JSON(err, model)
}

func QueryCredentialByIdDecorators() gw.Decorator {
	return gw.NewStoreDbSetupDecorator(func(ctx *gw.Context, db *gorm.DB) *gorm.DB {
		var user = ctx.User()
		if user.IsTenancy() {
			//return db.Where("")
		}
		return db
	})
}
