package auth

import "github.com/gin-gonic/gin"

type User struct {
	ID       string
	Passport string
}

func GetUser(c *gin.Context) User {
	return User{}
}