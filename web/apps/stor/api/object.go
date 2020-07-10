package api

import "github.com/gin-gonic/gin"

func CreateObject(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"payload": "Hello, world",
	})
}
