package api

import (
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func CreateObject(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": "CreateObject",
	})
}

func ModifyObject(c *gw2.ApiContext) {
	c.JSON(200, gin.H{
		"payload": "ModifyObject",
	})
}
