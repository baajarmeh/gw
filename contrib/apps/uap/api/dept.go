package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetDept(c *gw.Context) {
	db := c.Store().GetDbStore()
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s, db: %v", c.RequestID, c.Query("uid"), db),
	})
}
func CreateDept(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyDept(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteDept(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryDept(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}
