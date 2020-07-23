package auth

import (
	"github.com/oceanho/gw/conf"
)

type IAuth interface {
	Auth(passport, secret string,store conf.Backend) (User, error)
}

type IPerm interface {
	HasPerms(user User, store conf.Backend, perm string) (bool, error)
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
