package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func CreateObject(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": "CreateObject",
	})
}

func ModifyObject(c *app.ApiContext) {
	c.JSON(200, gin.H{
		"payload": "ModifyObject",
	})
}
