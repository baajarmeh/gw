package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/pvm"
	"github.com/oceanho/gw/contrib/apps/uap"
)

func main() {
	server := gw.DefaultServer()
	server.OnStart(registerSentry)
	server.Register(uap.New(), pvm.New())
	server.HandleErrors(serverErrorHandler, 400, 401, 403, 404)
	server.HandleErrors(server5XXErrorHandler, 500, 502, 503, 504)
	server.Serve()
}

func server5XXErrorHandler(requestId string, httpRequest string, headers []string, stack string, errOriBody string, err []*gin.Error) {
	sentry.CaptureMessage(stack)
}

func serverErrorHandler(requestId string, httpRequest string, headers []string, stack string, errOriBody string, err []*gin.Error) {
	sentry.CaptureMessage(errOriBody)
}

func registerSentry(state *gw.ServerState) {
	dsn := state.ApplicationConfig().Settings.Monitor.SentryDNS
	sentry.Init(sentry.ClientOptions{Dsn: dsn})
}
