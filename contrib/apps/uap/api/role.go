package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetRole(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestID, c.Query("uid")),
	})
}
func CreateRole(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyRole(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteRole(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryRole(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Role ID is %s", c.RequestID, c.Query("uid")),
	})
}
