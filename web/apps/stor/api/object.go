package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func CreateObject(ctx *app.ApiContext) interface{} {
	return gin.H{
		"payload": "CreateObject",
	}
}

func ModifyObject(ctx *app.ApiContext) interface{} {
	return gin.H{
		"payload": "ModifyObject",
	}
}
