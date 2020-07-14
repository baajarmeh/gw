package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetEnv(c *app.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateEnv(c *app.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyEnv(c *app.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyEnv(c *app.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}