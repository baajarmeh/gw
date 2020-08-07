package tester

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/rest"
)

func init() {
}

type AppRestOnly struct {
}

func NewAppRestOnly() AppRestOnly {
	return AppRestOnly{}
}

func (a AppRestOnly) Name() string {
	return "gw-rest-app"
}

func (a AppRestOnly) Router() string {
	return "tester"
}

func (a AppRestOnly) Register(router *gw.RouterGroup) {
	gw.RegisterRestAPI(router, &rest.MyTesterRestAPI{})
}

func (a AppRestOnly) Migrate(store gw.Store) {
}

func (a AppRestOnly) Use(opt *gw.ServerOption) {
}
