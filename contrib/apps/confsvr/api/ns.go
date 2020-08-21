package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func GetNS(c *gw.Context) {
	//user := c.User()
	c.JSON200(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateNS(c *gw.Context) {
	//user := c.User()
	c.JSON200(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyNS(c *gw.Context) {
	//user := c.User()
	c.JSON200(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyNS(c *gw.Context) {
	//user := c.User()
	c.JSON200(gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
