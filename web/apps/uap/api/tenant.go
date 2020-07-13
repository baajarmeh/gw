package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/contrib/app"
)

func GetTenant(ctx *app.ApiContext) interface{} {
	return gin.H{
		"payload": fmt.Sprintf("request id is: %s, user ID is %s", ctx.RequestId, ctx.Query("uid")),
	}
}
