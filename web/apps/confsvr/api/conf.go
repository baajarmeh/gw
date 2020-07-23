package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetConf(c *gw.Context) {
	//user := c.User
	c.OK(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
