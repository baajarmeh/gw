package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetUser(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
func CreateUser(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyUser(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteUser(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryUser(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
