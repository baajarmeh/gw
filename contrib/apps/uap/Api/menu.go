package Api

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
)

func GetMenu(c *gw.Context) {
	var id uint64
	if c.MustGetUint64IDFromParam(&id) != nil {
		return
	}
}

func GetMenuByName(c *gw.Context) {

}

func CreateMenu(c *gw.Context) {
	var appName string
	if c.MustParam("app", &appName) != nil {
		return
	}
	var dto Dto.Menu
	if c.Bind(&dto) != nil {
		return
	}
}

func BatchCreateMenu(c *gw.Context) {
	var dto Dto.Menu
	if c.Bind(&dto) != nil {
		return
	}
}

func ModifyMenu(c *gw.Context) {

}

func SearchMenuPageList(c *gw.Context) {

}

func QueryMenuPageList(c *gw.Context) {

}
