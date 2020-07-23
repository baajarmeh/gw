package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetRolePerms(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func CreateRolePerms(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyRolePerms(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteRolePerms(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryRolePerms(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, RolePerms ID is %s", c.RequestID, c.Query("uid")),
	})
}
