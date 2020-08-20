package gw

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/logger"
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

// JSON200 response a JSON formatter to client with http status = 200.
func (c *Context) JSON200(payload interface{}) {
	c.StatusJSON(http.StatusOK, 0, nil, payload)
}

// JSON400 response a JSON formatter to client with http status = 400.
func (c *Context) JSON400(status int) {
	c.JSON400Msg(status, errDefault400Msg)
}

// JSON400Msg response a JSON formatter to client with http status = 401.
func (c *Context) JSON400Msg(status int, errMsg interface{}) {
	c.JSON400PayloadMsg(status, errMsg, errDefaultPayload)
}

// JSON400Payload response a has payload properties JSON formatter to client with http status = 400.
func (c *Context) JSON400Payload(status int, payload interface{}) {
	c.JSON400PayloadMsg(status, errDefault400Msg, payload)
}

// JSON400PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) JSON400PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusBadRequest, status, errMsg, payload)
}

// JSON401 response a JSON formatter to client with http status = 401.
func (c *Context) JSON401(status int) {
	c.JSON401Msg(status, errDefault401Msg)
}

// JSON401Msg response a a has errMsg properties JSON formatter to client with http status = 401.
func (c *Context) JSON401Msg(status int, errMsg interface{}) {
	c.JSON401PayloadMsg(status, errMsg, errDefaultPayload)
}

// JSON401Payload response a has payload properties JSON formatter to client with http status = 401.
func (c *Context) JSON401Payload(status int, payload interface{}) {
	c.JSON401PayloadMsg(status, errDefault401Msg, payload)
}

// JSON401PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) JSON401PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusUnauthorized, status, errMsg, payload)
}

// JSON403 response a JSON formatter to client with http status = 403.
func (c *Context) JSON403(status int) {
	c.JSON403Msg(status, errDefault403Msg)
}

// JSON403Msg response a has errMsg properties JSON formatter to client with http status = 403.
func (c *Context) JSON403Msg(status int, errMsg interface{}) {
	c.JSON403PayloadMsg(status, errMsg, errDefaultPayload)
}

// JSON403Payload response a has payload properties JSON formatter to client with http status = 403.
func (c *Context) JSON403Payload(status int, payload interface{}) {
	c.JSON403PayloadMsg(status, errDefault403Msg, payload)
}

// JSON403PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 403.
func (c *Context) JSON403PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusForbidden, status, errMsg, payload)
}

// JSON404 response a JSON formatter to client with http status = 404.
func (c *Context) JSON404(status int) {
	c.JSON404Msg(status, errDefault404Msg)
}

// JSON404Msg response a a has errMsg properties JSON formatter to client with http status = 404.
func (c *Context) JSON404Msg(status int, errMsg interface{}) {
	c.JSON404PayloadMsg(status, errMsg, errDefaultPayload)
}

// JSON404Payload response a has payload properties JSON formatter to client with http status = 404.
func (c *Context) JSON404Payload(status int, payload interface{}) {
	c.JSON404PayloadMsg(status, payload, errDefault404Msg)
}

// JSON404PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 404.
func (c *Context) JSON404PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusNotFound, status, errMsg, payload)
}

// JSON500 response a JSON formatter to client with http status = 500.
func (c *Context) JSON500(status int) {
	c.JSON500Msg(status, nil)
}

// JSON500Msg response a has errMsg JSON formatter to client with http status = 500.
func (c *Context) JSON500Msg(status int, errMsg interface{}) {
	c.JSON500PayloadMsg(status, errMsg, errDefaultPayload)
}

// JSON500Payload response a has payload properties JSON formatter to client with http status = 500.
func (c *Context) JSON500Payload(status int, payload interface{}) {
	c.JSON500PayloadMsg(status, errDefault500Msg, payload)
}

// JSON response a response JSON by status.
// response 200,payload if status=0, other response 400, errMsg
func (c *Context) JSON(errMsg interface{}, payload interface{}) {
	switch ty := errMsg.(type) {
	case nil:
		c.JSON200(payload)
		break
	case error, *error, string:
		c.StatusJSON(http.StatusBadRequest, -1, errMsg, payload)
		break
	default:
		logger.Error("invalid errMsg types. %v", ty)
		break
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
	c.JSON200(pager)
}

// JSON500PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 500.
func (c *Context) JSON500PayloadMsg(status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(http.StatusInternalServerError, status, errMsg, payload)
}

// StatusJSON response a JSON formatter to client.
// Auto call c.Abort() when code < 200 || code > 202.
func (c *Context) StatusJSON(code int, status int, errMsg interface{}, payload interface{}) {
	s := hostServer(c.Context)
	c.Context.JSON(code, s.RespBodyBuildFunc(status, c.RequestID, errMsg, payload))
}

func defaultRespBodyBuildFunc(status int, requestID string, errMsg interface{}, payload interface{}) interface{} {
	return gin.H{
		"Status":    status,
		"Error":     errMsg,
		"RequestId": requestID,
		"Payload":   payload,
	}
}
