package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetDept(c *app.ApiContext) {
	db := c.Store.GetDbStore()
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s, db: %v", c.RequestId, c.Query("uid"), db),
	})
}
func CreateDept(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyDept(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteDept(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryDept(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, Dept ID is %s", c.RequestId, c.Query("uid")),
	})
}
