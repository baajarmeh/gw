package gw

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	errDefault400Msg              = "Bad Request"
	errDefault401Msg              = "Unauthorized"
	errDefault403Msg              = "Access Denied"
	errDefault404Msg              = "Not Found"
	errDefault500Msg              = "Internal Server Error"
	errDefaultPayload interface{} = nil
)

// OK response a JSON formatter to client with http status = 200.
func (c *Context) OK(payload interface{}) {
	c.StatusJSON(http.StatusOK, 0, nil, payload)
}

// Err400 response a JSON formatter to client with http status = 400.
func (c *Context) Err400(status int) {
	c.Err400Msg(status, errDefault400Msg)
}

// Err400Msg response a JSON formatter to client with http status = 401.
func (c *Context) Err400Msg(status int, errMsg interface{}) {
	c.Err400PayloadMsg(status, errMsg, errDefaultPayload)
}

// Err400Payload response a has payload properties JSON formatter to client with http status = 400.
func (c *Context) Err400Payload(status int, payload interface{}) {
	c.Err400PayloadMsg(status, errDefault400Msg, payload)
}

// Err400PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err400PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusBadRequest, status, errMsg, payload)
}

// Err401 response a JSON formatter to client with http status = 401.
func (c *Context) Err401(status int) {
	c.Err401Msg(status, errDefault401Msg)
}

// Err401Msg response a a has errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err401Msg(status int, errMsg interface{}) {
	c.Err401PayloadMsg(status, errMsg, errDefaultPayload)
}

// Err401Payload response a has payload properties JSON formatter to client with http status = 401.
func (c *Context) Err401Payload(status int, payload interface{}) {
	c.Err401PayloadMsg(status, errDefault401Msg, payload)
}

// Err401PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err401PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusUnauthorized, status, errMsg, payload)
}

// Err403 response a JSON formatter to client with http status = 403.
func (c *Context) Err403(status int) {
	c.Err403Msg(status, errDefault403Msg)
}

// Err403Msg response a has errMsg properties JSON formatter to client with http status = 403.
func (c *Context) Err403Msg(status int, errMsg interface{}) {
	c.Err403PayloadMsg(status, errMsg, errDefaultPayload)
}

// Err403Payload response a has payload properties JSON formatter to client with http status = 403.
func (c *Context) Err403Payload(status int, payload interface{}) {
	c.Err403PayloadMsg(status, errDefault403Msg, payload)
}

// Err403PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 403.
func (c *Context) Err403PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusForbidden, status, errMsg, payload)
}

// Err404 response a JSON formatter to client with http status = 404.
func (c *Context) Err404(status int) {
	c.Err404Msg(status, errDefault404Msg)
}

// Err404Msg response a a has errMsg properties JSON formatter to client with http status = 404.
func (c *Context) Err404Msg(status int, errMsg interface{}) {
	c.Err404PayloadMsg(status, errMsg, errDefaultPayload)
}

// Err404Payload response a has payload properties JSON formatter to client with http status = 404.
func (c *Context) Err404Payload(status int, payload interface{}) {
	c.Err404PayloadMsg(status, payload, errDefault404Msg)
}

// Err404PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 404.
func (c *Context) Err404PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusNotFound, status, errMsg, payload)
}

// Err500 response a JSON formatter to client with http status = 500.
func (c *Context) Err500(status int) {
	c.Err500Msg(status, nil)
}

// Err500Msg response a has errMsg JSON formatter to client with http status = 500.
func (c *Context) Err500Msg(status int, errMsg interface{}) {
	c.Err500PayloadMsg(status, errMsg, errDefaultPayload)
}

// Err500Payload response a has payload properties JSON formatter to client with http status = 500.
func (c *Context) Err500Payload(status int, payload interface{}) {
	c.Err500PayloadMsg(status, errDefault500Msg, payload)
}

// JSON response a response JSON by status.
// response 200,payload if status=0, other response 400, errMsg
func (c *Context) JSON(errMsg interface{}, payload interface{}) {
	if errMsg == nil {
		c.OK(payload)
	} else {
		c.StatusJSON(http.StatusBadRequest, -1, errMsg, payload)
	}
}

// Fault response a response fail json by status.
func (c *Context) Fault(errMsg interface{}) {
	c.JSON(errMsg, nil)
}

// PagerJSON response a response Pager's json by status.
func (c *Context) PagerJSON(total int64, expr PagerExpr, data interface{}) {
	pager := gin.H{
		"Data":       data,
		"Total":      total,
		"PageSize":   expr.PageSize,
		"PageNumber": expr.PageNumber,
	}
	c.OK(pager)
}

// Err500PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 500.
func (c *Context) Err500PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusInternalServerError, status, errMsg, payload)
}

// StatusJSON response a JSON formatter to client.
// Auto call c.Abort() when code <= 200 || code >= 202.
func (c *Context) StatusJSON(code int, status int, errMsg interface{}, payload interface{}) {
	c.Context.JSON(code, resp(status, c.RequestID, errMsg, payload))
	if code <= 200 || code >= 202 {
		c.Abort()
	}
}

func resp(status int, requestID string, errMsg interface{}, payload interface{}) interface{} {
	return gin.H{
		"Status":    status,
		"ErrMsg":    errMsg,
		"RequestID": requestID,
		"Payload":   payload,
	}
}
