package req

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func GetRequestId(c *gin.Context) string {
	requestId := c.GetHeader("X-Request-Id")
	if requestId == "" {
		requestId = internalGenRequestId()
	}
	return requestId
}

func internalGenRequestId() string {
	return fmt.Sprintf("gw-%d", time.Now().UnixNano())
}
