package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/pvm"
	"github.com/oceanho/gw/contrib/apps/uap"
	"github.com/oceanho/gw/libs/gwjsoner"
	"github.com/oceanho/gw/logger"
)

func main() {
	server := gw.DefaultServer()
	fixServer(server)
	server.Serve()
}

func fixServer(server *gw.HostServer) {
	server.OnStart(registerSentry).
		Register(uap.New(), pvm.New()).
		HandleErrors(server4XXErrorHandler, 400, 401, 403, 404).
		HandleErrors(server5XXErrorHandler, 500, 502, 503, 504)
}

func server5XXErrorHandler(requestId string, statusCode int, httpRequest string, headers []string, stack string, errOriBody string, err []*gin.Error) {
	sentry.CaptureMessage(warpSentryMessage(requestId, statusCode, httpRequest, headers, stack, errOriBody, err))
}

func server4XXErrorHandler(requestId string, statusCode int, httpRequest string, headers []string, stack string, errOriBody string, err []*gin.Error) {
	sentry.CaptureMessage(warpSentryMessage(requestId, statusCode, httpRequest, headers, stack, errOriBody, err))
}

func registerSentry(state *gw.ServerState) {
	dsn := state.ApplicationConfig().Settings.Monitor.SentryDNS
	sentry.Init(sentry.ClientOptions{Dsn: dsn})
}

func warpSentryMessage(requestId string, statusCode int, httpRequest string, headers []string, stack string, errOriBody string, err []*gin.Error) string {
	var body = gin.H{
		"requestId":         requestId,
		"httpRequest":       httpRequest,
		"headers":           headers,
		"stack":             stack,
		"oriRespBody":       errOriBody,
		"oriRespStatusCode": statusCode,
		"errs":              err,
	}
	b, e := gwjsoner.Marshal(body)
	if e != nil {
		logger.Error("warpSentryMessage fail, err: %v", e)
	}
	return string(b)
}
