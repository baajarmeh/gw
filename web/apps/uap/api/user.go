package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetUser(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
func CreateUser(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyUser(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteUser(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryUser(c *gw.Context) {
	c.OK( gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
