package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetUserPerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func CreateUserPerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyUserPerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteUserPerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryUserPerms(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestId, c.Query("uid")),
	})
}
