package app

import "github.com/gin-gonic/gin"

type User struct {
	Id       uint64
	TenantId uint64
}

func GetUser(c *gin.Context) User {
	return User{}
}
