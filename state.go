package gw

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"github.com/oceanho/gw/utils/secure"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	gwAppNameKey  = "gw-app"
	gwSidStateKey = "gw-sid-state"
)

func globalState(serverName string) gin.HandlerFunc {
	// code copies from gin framework.
	var out io.Writer = os.Stderr
	var ginLogger *log.Logger
	if out != nil {
		ginLogger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *gin.Context) {
		//
		// host server state.
		// 1. register the HostServer state into gin.Context
		// 2. process request, try got User from the http requests.
		//
		c.Set(gwAppNameKey, serverName)
		s := hostServer(c)
		requestId := getRequestID(c)

		defer func() {
			var stacks []byte
			var httpRequest []byte
			var headers []string
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if ginLogger != nil {
					stacks = stack(3)
					httpRequest, _ = httputil.DumpRequest(c.Request, false)
					headers = strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					if brokenPipe {
						ginLogger.Printf("requestId: %s, %s\n%s%s", requestId, err, string(httpRequest), reset)
					} else if gin.IsDebugging() {
						ginLogger.Printf("[Recovery] requestId: %s, %s panic recovered:\n%s\n%s\n %s %s", requestId,
							timeFormat(time.Now()), strings.Join(headers, "\r\n"), err, stacks, reset)
					} else {
						ginLogger.Printf("[Recovery] requestId: %s, %s panic recovered:\n%s\n%s%s", requestId,
							timeFormat(time.Now()), err, stacks, reset)
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					body := respBody(500, requestId, errDefault500Msg, nil)
					c.JSON(http.StatusInternalServerError, body)
				}
			}
			// handle panic Errors
			status := c.Writer.Status()
			handlers, ok := s.httpErrHandlers[status]
			defer func() {
				if err := recover(); err != nil {
					stacks = stack(3)
					logger.Error("handler or hooks errors fail, err: %s", string(stacks))
				}
			}()
			if ok {
				for _, handler := range handlers {
					handler(requestId, string(httpRequest), headers, string(stacks), c.Errors)
				}
			}
			// After handlers.
			for _, hook := range s.afterHooks {
				hook.Handler(c)
			}
		}()

		// before handlers
		for _, hook := range s.beforeHooks {
			hook.Handler(c)
		}

		//// has processed, return.
		//if c.Writer.Status() != 0 {
		//	return
		//}

		// gw framework handler.
		sid, ok := getSid(s, c)
		if ok {
			c.Set(gwSidStateKey, sid)
			user, err := s.sessionStateManager.Query(s.store, sid)
			if err == nil && user != nil {
				// set User State.
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
	sid, ok := c.Get(gwSidStateKey)
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
	serverName := c.MustGet(gwAppNameKey).(string)
	return servers[serverName]
}

func config(c *gin.Context) conf.Config {
	return *hostServer(c).conf
}

//
// code copies from gin framework.

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

const (
	reset = "\033[0m"
)

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
