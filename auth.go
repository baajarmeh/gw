package gw

import (
	"github.com/gin-gonic/gin"
)

type IAuth interface {
	Auth(passport, secret string, store Store) (User, error)
}

type IPermission interface {
	HasPerms(user User, store Store, perms ...string) (bool, error)
	HasAnyPerms(user User, store Store, perms ...string) (bool, error)
}

const UserKey = "gw-user"

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)
		if !user.Auth() {
			// Check url are allow dict.
			return
		}
		c.Set(UserKey, c)
		c.Next()
	}
}

type SessionAuth struct {
}

type CookieAuth struct {
}
