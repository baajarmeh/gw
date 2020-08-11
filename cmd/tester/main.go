package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/contrib/apps/tester"
	"github.com/oceanho/gw/logger"
	"strings"
	"time"
)

func main() {
	bcs := conf.DefaultBootConfig()
	opts := gw.NewServerOption(bcs)
	opts.Name = "my-tester-api"
	server := gw.New(opts)

	server.AddHook(gw.NewHook("my-tester-before", func(c *gin.Context) {
		c.Set("my-tester-id", 1000)
		c.Set("my-tester-start-at", time.Now().UnixNano())
	}, func(c *gin.Context) {
		id := c.MustGet("my-tester-id").(int)
		startAt, _ := c.MustGet("my-tester-start-at").(int64)
		nanoSeconds := time.Now().UnixNano() - startAt
		logger.Info("mytestid: %d, cost Nano Second: %d", id, nanoSeconds)
	}))

	server.HandleError(500, func(requestId, httpRequest string, headers []string, stack string, err []*gin.Error) {
		msgs := make([]string, 0)
		msgs = append(msgs, fmt.Sprintf("GW-500-Custom Handler"))
		msgs = append(msgs, fmt.Sprintf("requestId: %s", requestId))
		msgs = append(msgs, fmt.Sprintf("request: %s", httpRequest))
		msgs = append(msgs, fmt.Sprintf("headers: %s", strings.Join(headers, "\r\n")))
		msgs = append(msgs, fmt.Sprintf("stack: %s", stack))
		msgs = append(msgs, fmt.Sprintf("errors: %v", err))
		logger.Error("got a 500 error. \r\n%s", strings.Join(msgs, "\r\n=====================\r\n"))
	})

	registerApps(server)
	server.Serve()
}

func registerApps(server *gw.HostServer) {
	server.Register(tester.NewAppRestOnly())
	server.Register(tester.New())
}
