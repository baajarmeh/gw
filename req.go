package gw

import (
	"github.com/gin-gonic/gin"
)

var requestIdStateKey = "gw-req-id"

func getRequestId(s *HostServer, c *gin.Context) string {
	shouldSave := true
	requestID := c.GetHeader(s.conf.Settings.HeaderKey.RequestIDKey)
	if requestID == "" {
		requestID = c.GetString(requestIdStateKey)
		if requestID == "" {
			requestID = s.IDGenerator.NewStrID(32)
		} else {
			shouldSave = false
		}
	}
	if requestID != "" && shouldSave {
		c.Set(requestIdStateKey, requestID)
	}
	return requestID
}

// PagerExpr represents a general pager query request model for gw framework.
type PagerExpr struct {
	PageSize   int `json:"ps" binding:"required" query:"ps" params:"ps" form:"ps"`
	PageNumber int `json:"pn" binding:"required" query:"pn" params:"pn" form:"pn"`
}

// PagerResult represents a general pager query result model for gw framework.
type PagerResult struct {
	PagerExpr
	Total int64 `json:"total" query:"total" params:"total" form:"total"`
}

// DefaultPagerExpr returns a new PagerExpr
func DefaultPagerExpr(pageSize, pageNumber int) PagerExpr {
	return PagerExpr{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}
}

// SearcherExpr represents a general searcher query request model for gw framework.
type SearcherExpr struct {
	Field      string `json:"f" query:"f" params:"f" form:"f"`
	SearchMode string `json:"m" query:"m" params:"m" form:"m"`
}

// SearcherGroupExpr represents a general group searcher query request model for gw framework.
type SearcherGroupExpr struct {
	Left      SearcherExpr   `json:"l" query:"l" params:"l" form:"l"`
	LrMode    string         `json:"m" query:"m" params:"m" form:"m"`
	Right     SearcherExpr   `json:"r" query:"r" params:"r" form:"r"`
	GroupMode string         `json:"gm" query:"gm" params:"gm" form:"gm"`
	SubGroup  []SearcherExpr `json:"sg" query:"sg" params:"sg" form:"sg"`
}

// RangeExpr represents a general range query request model for gw framework.
type RangeExpr struct {
	Field string      `json:"f" query:"f" params:"f" form:"f"`
	Left  interface{} `json:"l" query:"l" params:"left" form:"l"`
	Right interface{} `json:"r" query:"r" params:"r" form:"r"`
}

// RangeGroupExpr represents a general group range query request model for gw framework.
type RangeGroupExpr struct {
	Left      RangeExpr        `json:"l" query:"l" params:"l" form:"l"`
	LRMode    string           `json:"lrm" query:"lrm" params:"lr_mode" form:"lrm"`
	Right     RangeExpr        `json:"r" query:"r" params:"r" form:"r"`
	GroupMode string           `json:"gm" query:"gm" params:"gm" form:"gm"`
	SubGroup  []RangeGroupExpr `json:"sg" query:"sg" params:"sg" form:"sg"`
}

// OrderlyExpr represents a general sort request model for gw framework.
type OrderlyExpr struct {
	Field     string `json:"f" query:"f" params:"f" form:"f"`
	Direction string `json:"d" query:"d" params:"d" form:"d"`
}

// QueryExpr represents a general query request model for gw framework.
type QueryExpr struct {
	PagerExpr
	Searcher      []SearcherExpr      `json:"s" query:"s" params:"s" form:"s"`
	Ranger        []RangeExpr         `json:"r" query:"r" params:"r" form:"r"`
	SearcherGroup []SearcherGroupExpr `json:"sg" query:"sg" params:"sg" form:"sg"`
	RangerGroup   []RangeGroupExpr    `json:"rg" query:"rg" params:"rg" form:"rg"`
	Orderly       []OrderlyExpr       `json:"o" query:"o" params:"o" form:"o"`
}

// PageOffset returns a Pager Offset Value.
func (expr PagerExpr) PageOffset() int {
	p := expr.PageNumber - 1
	if p < 0 {
		p = 0
	}
	return p * expr.PageSize
}
