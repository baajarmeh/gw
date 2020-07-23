package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetUserPerms(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func CreateUserPerms(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyUserPerms(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteUserPerms(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryUserPerms(c *gw.Context) {
	c.OK(gin.H{
		"payload": fmt.Sprintf("request id is: %s, UserPerms ID is %s", c.RequestID, c.Query("uid")),
	})
}
