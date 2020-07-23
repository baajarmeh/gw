package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetRole(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}
func CreateRole(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyRole(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteRole(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryRole(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}
