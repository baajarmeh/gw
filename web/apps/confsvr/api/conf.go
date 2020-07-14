package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetConf(c *app.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
