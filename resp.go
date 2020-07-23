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
	c.Err400PayloadMsg(status, errDefaultPayload, errMsg)
}

// Err400Payload response a has payload properties JSON formatter to client with http status = 400.
func (c *Context) Err400Payload(status int, payload interface{}) {
	c.Err400PayloadMsg(status, payload, errDefault400Msg)
}

// Err400PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err400PayloadMsg(status int, payload interface{}, errMsg interface{}) {
	c.ErrJSON(http.StatusBadRequest, status, payload, errMsg)
}

// Err401 response a JSON formatter to client with http status = 401.
func (c *Context) Err401(status int) {
	c.Err401Msg(status, errDefault401Msg)
}

// Err401Msg response a a has errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err401Msg(status int, errMsg interface{}) {
	c.Err401PayloadMsg(status, errDefaultPayload, errMsg)
}

// Err401Payload response a has payload properties JSON formatter to client with http status = 401.
func (c *Context) Err401Payload(status int, payload interface{}) {
	c.Err401PayloadMsg(status, payload, errDefault401Msg)
}

// Err401PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 401.
func (c *Context) Err401PayloadMsg(status int, payload interface{}, errMsg interface{}) {
	c.ErrJSON(http.StatusUnauthorized, status, payload, errMsg)
}

// Err403 response a JSON formatter to client with http status = 403.
func (c *Context) Err403(status int) {
	c.Err403Msg(status, errDefault403Msg)
}

// Err403Msg response a has errMsg properties JSON formatter to client with http status = 403.
func (c *Context) Err403Msg(status int, errMsg interface{}) {
	c.Err403PayloadMsg(status, errDefaultPayload, errMsg)
}

// Err403Payload response a has payload properties JSON formatter to client with http status = 403.
func (c *Context) Err403Payload(status int, payload interface{}) {
	c.Err403PayloadMsg(status, payload, errDefault400Msg)
}

// Err403PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 403.
func (c *Context) Err403PayloadMsg(status int, payload interface{}, errMsg interface{}) {
	c.ErrJSON(http.StatusForbidden, status, payload, errMsg)
}

// Err404 response a JSON formatter to client with http status = 404.
func (c *Context) Err404(status int) {
	c.Err404Msg(status, errDefault404Msg)
}

// Err404Msg response a a has errMsg properties JSON formatter to client with http status = 404.
func (c *Context) Err404Msg(status int, errMsg interface{}) {
	c.Err404PayloadMsg(status, errDefaultPayload, errMsg)
}

// Err404Payload response a has payload properties JSON formatter to client with http status = 404.
func (c *Context) Err404Payload(status int, payload interface{}) {
	c.Err404PayloadMsg(status, payload, errDefault404Msg)
}

// Err404PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 404.
func (c *Context) Err404PayloadMsg(status int, payload interface{}, errMsg interface{}) {
	c.ErrJSON(http.StatusNotFound, status, payload, errMsg)
}

// Err500 response a JSON formatter to client with http status = 500.
func (c *Context) Err500(status int) {
	c.Err500Msg(status, nil)
}

// Err500Msg response a has errMsg JSON formatter to client with http status = 500.
func (c *Context) Err500Msg(status int, errMsg interface{}) {
	c.Err500PayloadMsg(status, errDefaultPayload, errMsg)
}

// Err500Payload response a has payload properties JSON formatter to client with http status = 500.
func (c *Context) Err500Payload(status int, payload interface{}) {
	c.Err500PayloadMsg(status, payload, errDefault500Msg)
}

// Err500PayloadMsg response a has payload,errMsg properties JSON formatter to client with http status = 500.
func (c *Context) Err500PayloadMsg(status int, payload interface{}, errMsg interface{}) {
	c.ErrJSON(http.StatusInternalServerError, status, payload, errMsg)
}

// ErrJSON response a has code,status,errMsg,payload properties JSON formatter to client.
func (c *Context) ErrJSON(code, status int, errMsg interface{}, payload interface{}) {
	c.StatusJSON(code, status, errMsg, payload)
}

// StatusJSON response a JSON formatter to client.
func (c *Context) StatusJSON(code int, status int, errMsg interface{}, payload interface{}) {
	c.Context.JSON(code, resp(status, c.RequestID, errMsg, payload))
}

func resp(status int, requestID string, errMsg interface{}, payload interface{}) interface{} {
	return gin.H{
		"Status":    status,
		"Err":       errMsg,
		"RequestID": requestID,
		"Payload":   payload,
	}
}
