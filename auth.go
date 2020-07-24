package gw

type IAuth interface {
	Auth(passport, secret string, store Store) (User, error)
}

type IPermission interface {
	HasPerms(user User, store Store, perms ...string) (bool, error)
	HasAnyPerms(user User, store Store, perms ...string) (bool, error)
}

const UserKey = "gw-user"

// func Auth(a IAuth, realm string) gin.HandlerFunc {
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
