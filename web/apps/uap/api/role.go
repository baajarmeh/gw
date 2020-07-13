package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetRole(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}
func CreateRole(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyRole(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteRole(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryRole(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestId, c.Query("uid")),
	})
}
