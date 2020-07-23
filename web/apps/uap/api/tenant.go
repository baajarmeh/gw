package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/logger"
)

func GetTenant(c *gw.Context) {
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
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
func CreateTenant(c *gw.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyTenant(c *gw.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteTenant(c *gw.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryTenant(c *gw.Context) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
