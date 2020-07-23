package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetAuth(c *gw.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.OK(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateAuth(c *gw.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.OK(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyAuth(c *gw.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.OK(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyAuth(c *gw.Context) {
	//accessKeyId := c.GetHeader("X-Access-KeyId")
	//accessSecret := c.GetHeader("X-Access-Secret")
	c.OK(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
