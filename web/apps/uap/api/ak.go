package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetAK(c *app.ApiContext) {
	db := c.Store.GetDbStore()
	row := db.Raw("select 1 from DUAL").Row()
	var result uint64 = 0
	err := row.Scan(&result)

	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s, db result: %d, db err: %v",
			c.RequestId, c.Query("uid"), result, err),
	})
}
func CreateAK(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func ModifyAK(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func DeleteAK(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}

func QueryAK(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", c.RequestId, c.Query("uid")),
	})
}
