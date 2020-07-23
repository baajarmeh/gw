package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetPermission(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
func CreatePermission(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyPermission(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeletePermission(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryPermission(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func Grant(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func Revoke(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
