package gw

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
)

const stateKey = "gw-app-state"

func gwState(serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(stateKey, serverName)
		c.Next()
	}
}

func getHostServer(c *gin.Context) *HostServer {
	serverName := c.MustGet(stateKey).(string)
	return servers[serverName]
}

func getConfig(c *gin.Context) conf.Config {
	return *getHostServer(c).conf
}
