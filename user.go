package gw

import "github.com/gin-gonic/gin"

type User struct {
	Id       uint64
	TenantId uint64
	Passport string
	RoleId   int // Platform Admin:1, Tenant Admin:2, roleId >= 10000 are custom role.
}

func (usr User) Auth() bool {
	return &usr != nil && usr.Id > 0
}

func (usr User) IsAdmin() bool {
	return usr.Auth() && usr.RoleId == 1
}

func (usr User) IsTenantAdmin() bool {
	return usr.Auth() && usr.RoleId == 2
}

func getUser(c *gin.Context) User {
	obj, ok := c.Get(UserKey)
	if !ok {
		obj = User{}
	}
	return obj.(User)
}
