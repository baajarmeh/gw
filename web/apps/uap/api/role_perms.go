package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetRolePerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func CreateRolePerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyRolePerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteRolePerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryRolePerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestId, c.Query("uid")),
	})
}
