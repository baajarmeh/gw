package Api

import "github.com/oceanho/gw"

func GetMenu(c *gw.Context) {
	var id uint64
	if c.MustGetUint64IDFromParam(&id) != nil {
		return
	}
}

func GetMenuByName(c *gw.Context) {

}

func CreateMenu(c *gw.Context) {

}

func BatchCreateMenu(c *gw.Context) {

}

func ModifyMenu(c *gw.Context) {

}

func SearchMenuPageList(c *gw.Context) {

}

func QueryMenuPageList(c *gw.Context) {

}
