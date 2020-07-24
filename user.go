package gw

import "github.com/gin-gonic/gin"

type User struct {
	Id       uint64
	TenantId uint64
	Passport string
}

func getUser(c *gin.Context) User {
	return User{}
}
