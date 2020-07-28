package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var requestIdStateKey = "gwState-X-Request-Id"

func getRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = c.GetString(requestIdStateKey)
		if requestID == "" {
			requestID = internalGenRequestID()
		}
	}
	c.Set(requestIdStateKey, requestID)
	return requestID
}

func internalGenRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// PagerExpr represents a general pager query request model for gw framework.
type PagerExpr struct {
	PageSize   int `json:"pageSize" binding:"required" query:"pageSize" params:"pageSize" form:"pageSize"`
	PageNumber int `json:"pageNumber" binding:"required" query:"pageNumber" params:"pageNumber" form:"pageNumber"`
}

// SearcherExpr represents a general searcher query request model for gw framework.
type SearcherExpr struct {
	Field      string `json:"field" query:"field" params:"field" form:"field"`
	SearchMode string `json:"mode" query:"mode" params:"mode" form:"mode"`
}

// RangeExpr represents a general range query request model for gw framework.
type RangeExpr struct {
	Field string      `json:"field" query:"field" params:"field" form:"field"`
	Left  interface{} `json:"left" query:"left" params:"left" form:"left"`
	Right interface{} `json:"right" query:"right" params:"right" form:"right"`
}

// OrderlyExpr represents a general sort request model for gw framework.
type OrderlyExpr struct {
	Field     string `json:"field" query:"field" params:"field" form:"field"`
	Direction string `json:"direction" query:"direction" params:"direction" form:"direction"`
}

// QueryExpr represents a general query request model for gw framework.
type QueryExpr struct {
	PagerExpr
	Searcher []SearcherExpr `json:"s" query:"s" params:"s" form:"s"`
	Ranger   []RangeExpr    `json:"r" query:"r" params:"r" form:"r"`
	Orderly  []OrderlyExpr  `json:"o" query:"o" params:"o" form:"o"`
}

func (expr PagerExpr) PageOffset() int {
	p := expr.PageNumber - 1
	if p < 0 {
		p = 0
	}
	return p * expr.PageSize
}
