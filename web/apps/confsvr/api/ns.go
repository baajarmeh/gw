package api

import (
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetNS(c *gw2.Context) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateNS(c *gw2.Context) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyNS(c *gw2.Context) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyNS(c *gw2.Context) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
