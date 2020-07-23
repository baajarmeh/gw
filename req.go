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

// QueryExpr represents a general query request model for gw framework.
type QueryExpr struct {
	Pager struct {
		PageSize   int `json:"pageSize" binding:"required" xml:"pageSize" query:"pageSize" params:"pageSize" from:"pageSize"`
		PageNumber int `json:"pageNumber" binding:"required" xml:"pageNumber" query:"pageNumber" params:"pageNumber" from:"pageNumber"`
	} `json:"pager" binding:"required" xml:"pager" query:"pager" params:"pager" from:"pager"`
	Expr struct {
		Equal struct {
			Field string `json:"field" query:"field" params:"field" from:"field"`
			Value string `json:"value" query:"value" params:"value" from:"value"`
			// Multiple
			Fields []string      `json:"fields" query:"fields" params:"fields" from:"fields"`
			Values []interface{} `json:"values" query:"values" params:"values" from:"values"`
		} `json:"equal" query:"equal" params:"equal" from:"equal"`
		Range struct {
			Field string `json:"field" query:"field" params:"field" from:"field"`
			Value string `json:"value" query:"value" params:"value" from:"value"`
			// Multiple
			Fields []string      `json:"fields" query:"fields" params:"fields" from:"fields"`
			Values []interface{} `json:"values" query:"values" params:"values" from:"values"`
		} `json:"range" query:"range" params:"range" from:"range"`
		Search struct {
			Field string `json:"field" query:"field" params:"field" from:"field"`
			Value string `json:"value" query:"value" params:"value" from:"value"`
			// Multiple
			Fields     []string      `json:"fields" query:"fields" params:"fields" from:"fields"`
			Values     []interface{} `json:"values" query:"values" params:"values" from:"values"`
			SearchMode string        `json:"mode" query:"mode" params:"mode" from:"mode"`
		} `json:"search" query:"search" params:"search" from:"search"`
	} `json:"expr" query:"expr" params:"expr" from:"expr"`
	Order []struct {
		Field     string `json:"field" query:"field" params:"field" from:"field"`
		Direction string `json:"direction" query:"direction" params:"direction" from:"direction"`
	} `json:"order" query:"order" params:"order" from:"order"`
}
