package api

import (
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetAuth(c *gw2.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateAuth(c *gw2.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyAuth(c *gw2.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyAuth(c *gw2.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
