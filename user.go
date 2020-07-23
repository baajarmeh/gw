package gw

import "github.com/gin-gonic/gin"

type User struct {
	Id       uint64
	TenantId uint64
}

func getUser(c *gin.Context) User {
	return User{}
}
