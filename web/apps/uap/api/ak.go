package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetAK(c *gw2.ApiContext) {
	db := c.Store.GetDbStore()
	row := db.Raw("select 1 from DUAL").Row()
	var result uint64 = 0
	err := row.Scan(&result)

	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s, db result: %d, db err: %v",
			c.RequestId, c.Query("uid"), result, err),
	})
}
func CreateAK(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyAK(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteAK(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryAK(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
