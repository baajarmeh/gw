package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/contrib/app/logger"
)

func GetTenant(c *app.ApiContext) {
	db := c.Store.GetDbStore()
	sqlDb, err := db.DB()
	if err != nil {
		logger.Error("get sql.Db fail.")
		return
	}
	err = sqlDb.Ping()
	if err != nil {
		logger.Error("sql.Db not pong.")
		return
	}
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
func CreateTenant(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyTenant(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteTenant(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryTenant(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
