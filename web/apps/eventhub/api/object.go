package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
)

func CreateObject(c *gw.Context) {
	c.OK( gin.H{
		"payload": "CreateObject",
	})
}

func ModifyObject(c *gw.Context) {
	c.OK( gin.H{
		"payload": "ModifyObject",
	})
}
