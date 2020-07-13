package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
	"github.com/oceanho/gw/contrib/app/resp"
)

func CreateObject(ctx *app.ApiContext) resp.IApiResp {
	return resp.HttpApiRespJSON(200, gin.H{
		"payload": "CreateObject",
	})
}

func ModifyObject(ctx *app.ApiContext) interface{} {
	return gin.H{
		"payload": "ModifyObject",
	}
}
