package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetDept(c *gw2.Context) {
	db := c.Store.GetDbStore()
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s, db: %v", c.RequestID, c.Query("uid"), db),
	})
}
func CreateDept(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyDept(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteDept(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryDept(c *gw2.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestID, c.Query("uid")),
	})
}
