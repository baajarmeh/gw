package restAPIs

import "github.com/oceanho/gw"

func Register(router *gw.RouterGroup) {
	router.RegisterRestAPIs(&User{})
}
