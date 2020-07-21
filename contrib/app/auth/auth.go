package auth

import (
	"github.com/gin-gonic/gin"
	// "strconv"
)

type IAuth interface {
	Auth(c *gin.Context) (User, error)
}

const UserKey = "gw-user"

// func ApiAuth(a IAuth, realm string) gin.HandlerFunc {
//
// 	if realm == "" {
// 		realm = "Authorization Required"
// 	}
// 	realm = "Basic realm=" + strconv.Quote(realm)
// 	return func(c *gin.Context) {
// 		basicAuth,ok := c.Get(gin.AuthUserKey)
// 		if !ok {
// 			return
// 		}
// 		c.Set(UserKey, User{})
// 		c.Next()
// 	}
// }

type SessionAuth struct {
}

type CookieAuth struct {
}
