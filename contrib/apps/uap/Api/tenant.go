package Api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/logger"
)

func GetTenant(c *gw.Context) {
	db := c.Store().GetDbStore()
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
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user are %v", c.RequestId(), c.User()),
	})
}

func CreateTenant(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func ModifyTenant(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func DeleteTenant(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}

func QueryTenant(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId(), c.Query("uid")),
	})
}
