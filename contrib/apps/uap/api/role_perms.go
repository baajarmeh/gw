package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetRolePerms(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func CreateRolePerms(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func ModifyRolePerms(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func DeleteRolePerms(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func QueryRolePerms(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId(), c.Query("uid")),
	})
}
