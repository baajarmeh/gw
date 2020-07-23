package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetUserPerms(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func CreateUserPerms(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyUserPerms(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteUserPerms(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryUserPerms(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}
