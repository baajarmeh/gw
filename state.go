package gw

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"github.com/oceanho/gw/utils/secure"
)

const stateKey = "gw-app-state"
const sidStateKey = "gw-sid-state"

func hostState(serverName string) gin.HandlerFunc {
	//
	// host server state.
	// 1. register the HostServer state into gin.Context
	// 2. process request, try got User from the http requests.
	//
	return func(c *gin.Context) {
		c.Set(stateKey, serverName)
		s := hostServer(c)
		sid, ok := getSid(s, c)
		if ok {
			c.Set(sidStateKey, sid)
			user, err := s.sessionStore.Query(s.store, sid)
			if err == nil && user != nil {
				// Set User State.
				c.Set(gwUserKey, user)
			}
		}
		c.Next()
	}
}

func newSid(s *HostServer) (oriSig, secureSid string, ok bool) {
	sid := uuid.New().String()
	src := []byte(sid)
	dst := make([]byte, len(src))
	err := s.crypto.Encrypt(dst, src)
	if err != nil {
		logger.Error("newSid() -> s.crypto.Encrypt(dst,b) fail, err: %v.", err)
		return "", "", false
	}
	secureSid = secure.EncodeBase64URL(dst)
	return sid, secureSid, true
}

func getSid(s *HostServer, c *gin.Context) (string, bool) {
	sid, ok := c.Get(sidStateKey)
	if ok {
		return sid.(string), true
	}
	cks := s.conf.Service.Security.Auth.Cookie
	_sid, err := c.Cookie(cks.Key)
	if _sid == "" {
		return "", false
	}
	sids, err := secure.DecodeBase64URL(_sid)
	if err != nil {
		logger.Error("getSid() -> security.DecodeBase64URL(_sid) fail, err:%v . sid=%s", err, _sid)
		return "", false
	}
	dst := make([]byte, len(sids))
	err = s.crypto.Decrypt(dst, sids)
	if err != nil {
		logger.Error("getSid() -> s.crypto.Decrypt(dst,[]byte(sid)) fail, err:%v . sid=%s", err, _sid)
		return "", false
	}
	return string(dst), true
}

func hostServer(c *gin.Context) *HostServer {
	serverName := c.MustGet(stateKey).(string)
	return servers[serverName]
}

func config(c *gin.Context) conf.Config {
	return *hostServer(c).conf
}
