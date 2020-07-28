package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetAK(c *gw.Context) {
	db := c.Store.GetDbStore()
	row := db.Raw("select 1 from DUAL").Row()
	var result uint64 = 0
	err := row.Scan(&result)

	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s, db result: %d, db err: %v",
			c.RequestID, c.Query("uid"), result, err),
	})
}

func CreateAK(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func ModifyAK(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func DeleteAK(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}

func QueryAK(c *gw.Context) {
	c.JSON200(gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestID, c.Query("uid")),
	})
}
