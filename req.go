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

// PagerExpr represents a general pager query request model for gw framework.
type PagerExpr struct {
	PageSize   int `json:"pageSize" binding:"required" xml:"pageSize" query:"pageSize" params:"pageSize" form:"pageSize"`
	PageNumber int `json:"pageNumber" binding:"required" xml:"pageNumber" query:"pageNumber" params:"pageNumber" form:"pageNumber"`
}

// QueryExpr represents a general query request model for gw framework.
type QueryExpr struct {
	PagerExpr
	Searcher struct {
		Equal struct {
			Field string `json:"field" query:"field" params:"field" form:"field"`
			Value string `json:"value" query:"value" params:"value" form:"value"`
			// Multiple
			Fields []string      `json:"fields" query:"fields" params:"fields" form:"fields"`
			Values []interface{} `json:"values" query:"values" params:"values" form:"values"`
		} `json:"equal" query:"equal" params:"equal" form:"equal"`
		Range struct {
			Field string `json:"field" query:"field" params:"field" form:"field"`
			Value string `json:"value" query:"value" params:"value" form:"value"`
			// Multiple
			Fields []string      `json:"fields" query:"fields" params:"fields" form:"fields"`
			Values []interface{} `json:"values" query:"values" params:"values" form:"values"`
		} `json:"range" query:"range" params:"range" form:"range"`
		Search struct {
			Field string `json:"field" query:"field" params:"field" form:"field"`
			Value string `json:"value" query:"value" params:"value" form:"value"`
			// Multiple
			Fields     []string      `json:"fields" query:"fields" params:"fields" form:"fields"`
			Values     []interface{} `json:"values" query:"values" params:"values" form:"values"`
			SearchMode string        `json:"mode" query:"mode" params:"mode" form:"mode"`
		} `json:"search" query:"search" params:"search" form:"search"`
	} `json:"expr" query:"expr" params:"expr" form:"expr"`
	Order []struct {
		Field     string `json:"field" query:"field" params:"field" form:"field"`
		Direction string `json:"direction" query:"direction" params:"direction" form:"direction"`
	} `json:"order" query:"order" params:"order" form:"order"`
}

func (expr PagerExpr) PageOffset() int {
	p := expr.PageNumber - 1
	if p < 0 {
		p = 0
	}
	return p * expr.PageSize
}
