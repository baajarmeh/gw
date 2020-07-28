package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"net/http"
	"strings"
)

type IAuth interface {
	Auth(passport, secret string, store Store) (User, error)
}

type IPermission interface {
	HasPerms(user User, store Store, perms ...string) (bool, error)
	HasAnyPerms(user User, store Store, perms ...string) (bool, error)
}

const UserKey = "gw-user"

func auth(vars map[string]string, urls []conf.AllowUrl) gin.HandlerFunc {
	var allowUrls = make(map[string]bool)
	for _, url := range urls {
		for _, p := range url.Urls {
			s := p
			for k, v := range vars {
				s = strings.Replace(s, k, v, 4)
			}
			allowUrls[s] = true
		}
	}
	return func(c *gin.Context) {
		user := getUser(c)
		path := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		//
		// No auth and request URI not in allowed urls.
		// UnAuthorized
		//
		if (user == nil || !user.Auth()) && !allowUrls[path] {
			// Check url are allow dict.
			payload := resp(http.StatusUnauthorized, getRequestID(c), errDefault401Msg, errDefaultPayload)
			c.JSON(http.StatusUnauthorized, payload)
			c.Abort()
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
