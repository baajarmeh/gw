package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"github.com/oceanho/gw/utils/secure"
	"strings"
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
			user, err := s.sessionStateManager.Query(s.store, sid)
			if err == nil && user != nil {
				// Set User State.
				c.Set(gwUserKey, user)
			}
		}
		c.Next()
	}
}

func encryptSid(s *HostServer, passport string) (secureSid string, ok bool) {

	rdnKey := secure.RandomStr(32)

	block := secure.AesBlock(rdnKey)
	encryptor := secure.AesEncryptCFB(rdnKey, block)
	passportSrc := []byte(passport)
	passportDst := make([]byte, len(passportSrc))
	encryptor.XORKeyStream(passportDst, passportSrc)
	encryptedPassport := string(passportDst)

	rdnKeySrc := []byte(rdnKey)
	rdnKeyDst := make([]byte, len(rdnKeySrc))
	_ = s.protect.Encrypt(rdnKeyDst, rdnKeySrc)
	encryptedRdnKey := string(rdnKeyDst)

	sid := fmt.Sprintf("%s,%s", encryptedPassport, encryptedRdnKey)

	sidSrc := []byte(sid)
	sidDst := make([]byte, len(sidSrc))
	_ = s.hash.Hash(sidDst, sidSrc)
	sigSum := string(sidDst)

	// Append sid Hash.
	sid = fmt.Sprintf("%s,%s", sid, sigSum)

	src := []byte(sid)
	dst := make([]byte, len(src))
	err := s.protect.Encrypt(dst, src)
	if err != nil {
		logger.Error("encryptSid() -> s.crypto.Encrypt(dst,b) fail, err: %v.", err)
		return "", false
	}
	secureSid = secure.EncodeBase64URL(dst)
	return secureSid, true
}

func decryptSid(s *HostServer, secureSid string, client string) (passport string, ok bool) {
	sid, err := secure.DecodeBase64URL(secureSid)
	if err != nil {
		logger.Error("decryptSid() -> security.DecodeBase64URL(_sid) fail, err:%v . sid=%s", err, sid)
		return "", false
	}
	dst := make([]byte, len(sid))
	err = s.protect.Decrypt(dst, sid)
	if err != nil {
		logger.Error("decryptSid() -> s.crypto.Decrypt(dst,[]byte(sid)) fail, err:%v . sid=%s", err, sid)
		return "", false
	}
	sids := strings.Split(string(dst), ",")
	if len(sids) != 3 {
		logger.Warn("got a invalid secureSid data from %s", client)
		return "", false
	}

	encryptedSid := sids[0]
	encryptedRdnKey := sids[1]
	// Check data sign sum.
	sidSumOri := fmt.Sprintf("%s,%s", encryptedSid, encryptedRdnKey)
	sidSumSrc := []byte(sidSumOri)
	sidSumDst := make([]byte, len(sidSumSrc))
	_ = s.hash.Hash(sidSumDst, sidSumSrc)
	// data sum Not match, maybe has modified.
	if sids[2] != string(sidSumDst) {
		logger.Warn("got a invalid secureSid data sum from %s", client)
		return "", false
	}

	rdnKeySrc := []byte(encryptedRdnKey)
	rdnKeyDst := make([]byte, len(rdnKeySrc))
	_ = s.protect.Decrypt(rdnKeyDst, rdnKeySrc)
	rdnKey := string(rdnKeyDst)

	block := secure.AesBlock(rdnKey)
	decryptor := secure.AesDecryptCFB(rdnKey, block)
	sidSrc := []byte(encryptedSid)
	sidDst := make([]byte, len(sidSrc))
	decryptor.XORKeyStream(sidDst, sidSrc)
	return string(sidDst), true
}

func getClient(c *gin.Context) string {
	return c.Request.RemoteAddr
}

func getSid(s *HostServer, c *gin.Context) (string, bool) {
	sid, ok := c.Get(sidStateKey)
	if ok {
		return sid.(string), true
	}
	// Trust Sid header.
	// sid can be decrypt AT gateway (OpenResty, Nginx etc.) layout.
	sidKey := s.conf.Service.Security.Auth.TrustSidKey
	if sidKey != "" {
		sid := c.GetHeader(sidKey)
		if sid != "" {
			return sid, ok
		}
	}
	client := getClient(c)
	cks := s.conf.Service.Security.Auth.Cookie
	// 1. Query cookie get Sid.
	_sid, err := c.Cookie(cks.Key)
	if _sid == "" || err != nil {
		// 2. Query http header get token.
		_sid = c.GetHeader("X-Auth-Token")
		if _sid == "" {
			return "", false
		}
	}
	return decryptSid(s, _sid, client)
}

func hostServer(c *gin.Context) *HostServer {
	serverName := c.MustGet(stateKey).(string)
	return servers[serverName]
}

func config(c *gin.Context) conf.Config {
	return *hostServer(c).conf
}
