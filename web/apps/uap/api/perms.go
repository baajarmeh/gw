package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetPermission(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
func CreatePermission(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyPermission(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeletePermission(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryPermission(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func Grant(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func Revoke(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
