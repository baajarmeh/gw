package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func getRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = internalGenRequestID()
	}
	return requestID
}

func internalGenRequestID() string {
	return fmt.Sprintf("gw-%d", time.Now().UnixNano())
}
