package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetRolePerms(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func CreateRolePerms(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyRolePerms(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteRolePerms(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryRolePerms(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}
