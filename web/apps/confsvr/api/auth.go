package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetAuth(c *app.ApiContext) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateAuth(c *app.ApiContext) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyAuth(c *app.ApiContext) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyAuth(c *app.ApiContext) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
