package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"net/http"
	"strings"
	"time"
)

//
// GW framework login API.
func gwLogin(c *gin.Context) {
	s := getHostServer(c)
	reqId := getRequestId(s, c)
	var err error
	var hasCheckPass = false
	var checker = s.AuthParamChecker
	var authParam AuthParameter
	for _, resolver := range s.AuthParamResolvers {
		authParam = resolver.Resolve(c)
		if err = checker.Check(authParam); err == nil {
			hasCheckPass = true
			break
		}
	}

	if !hasCheckPass {
		c.JSON(http.StatusBadRequest, s.RespBodyBuildFunc(c, http.StatusBadRequest, reqId, err, nil))
		c.Abort()
		return
	}

	// Login
	user, err := s.AuthManager.Login(authParam)
	if err != nil || user.IsEmpty() {
		c.JSON(http.StatusNotFound, s.RespBodyBuildFunc(c, http.StatusNotFound, reqId, err, nil))
		c.Abort()
		return
	}
	sid, credential, ok := encryptSid(s, authParam)
	if !ok {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(c, http.StatusInternalServerError, reqId, "Create session ID fail.", nil))
		c.Abort()
		return
	}
	if err := s.SessionStateManager.Save(sid, user); err != nil {
		c.JSON(http.StatusInternalServerError, s.RespBodyBuildFunc(c, http.StatusInternalServerError, reqId, "Save session fail.", err.Error()))
		c.Abort()
		return
	}
	var userPerms []gin.H
	for _, p := range user.Permissions {
		userPerms = append(userPerms, gin.H{
			"Key":  p.Key,
			"Name": p.Name,
			"Desc": p.Descriptor,
		})
	}
	cks := s.conf.Security.Auth.Cookie
	expiredAt := time.Duration(cks.MaxAge) * time.Second
	var userRoles = gin.H{
		"Id":   0,
		"name": "",
		"desc": "",
	}
	payload := gin.H{
		"Credentials": gin.H{
			"Token":     credential,
			"ExpiredAt": time.Now().Add(expiredAt).Unix(),
		},
		"Roles":       userRoles,
		"Permissions": userPerms,
	}
	body := s.RespBodyBuildFunc(c, 0, reqId, nil, payload)
	c.SetCookie(cks.Key, credential, cks.MaxAge, cks.Path, cks.Domain, cks.Secure, cks.HttpOnly)
	c.JSON(http.StatusOK, body)
}

// GW framework logout API.
func gwLogout(c *gin.Context) {
	s := getHostServer(c)
	reqId := getRequestId(s, c)
	user := getUser(c)
	cks := s.conf.Security.Auth.Cookie
	ok := s.AuthManager.Logout(user)
	if !ok {
		s.RespBodyBuildFunc(c, http.StatusInternalServerError, reqId, "auth logout fail", nil)
		return
	}
	sid, ok := getSid(s, c)
	if !ok {
		s.RespBodyBuildFunc(c, http.StatusInternalServerError, reqId, "session store logout fail", nil)
		return
	}
	_ = s.SessionStateManager.Remove(sid)
	c.SetCookie(cks.Key, "", -1, cks.Path, cks.Domain, cks.Secure, cks.HttpOnly)
}

// GW framework auth Check Middleware
func gwAuthChecker(urls []conf.AllowUrl) gin.HandlerFunc {
	var allowUrls = make(map[string]bool)
	for _, url := range urls {
		for _, p := range url.Urls {
			s := p
			allowUrls[s] = true
		}
	}
	return func(c *gin.Context) {
		s := getHostServer(c)
		user := getUser(c)
		path := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		requestId := getRequestId(s, c)
		//
		// No auth and request URI not in allowed urls.
		// UnAuthorized
		//
		if (user.IsEmpty() || !user.IsAuth()) && !allowUrls[path] {
			auth := s.conf.Security.AuthServer
			// Check url are allow dict.
			payload := gin.H{
				"Auth": gin.H{
					"LogIn": gin.H{
						"Url": fmt.Sprintf("%s/%s",
							strings.TrimRight(auth.Addr, "/"), strings.TrimLeft(auth.LogIn.Url, "/")),
						"Methods":   auth.LogIn.Methods,
						"AuthTypes": auth.LogIn.AuthTypes,
					},
					"LogOut": gin.H{
						"Url": fmt.Sprintf("%s/%s",
							strings.TrimRight(auth.Addr, "/"), strings.TrimLeft(auth.LogOut.Url, "/")),
						"Methods": auth.LogOut.Methods,
					},
				},
			}
			body := s.RespBodyBuildFunc(c, http.StatusUnauthorized, requestId, errDefault401Msg, payload)
			c.JSON(http.StatusUnauthorized, body)
			c.Abort()
			return
		}
		c.Next()
	}
}
