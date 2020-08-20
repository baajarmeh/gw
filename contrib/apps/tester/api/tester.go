package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/tester/biz"
	"github.com/oceanho/gw/contrib/apps/tester/dto"
)

func CreateMyTester(c *gw.Context) {
	obj := &dto.MyTester{}
	if c.Bind(obj) != nil {
		return
	}
	err := biz.CreateMyTester(c.State.Store().GetDbStore(), obj)
	c.JSON(err, obj.ID)
}

func QueryMyTester(c *gw.Context) {
	query := &gw.QueryExpr{}
	if c.Bind(query) != nil {
		return
	}
	objs := make([]dto.MyTester, 0)
	err := biz.QueryMyTester(c.State.Store().GetDbStore(), &objs)
	c.JSON(err, objs)
}

func GetTester(c *gw.Context) {
	c.JSON200(struct {
		RequestID string
	}{
		RequestID: c.RequestID,
	})
}

func GetTester400(c *gw.Context) {
	c.JSON400(4000)
}

func GetTester400WithCustomErr(c *gw.Context) {
	c.JSON400Msg(4000, "Custom 400 Err")
}

func GetTester400WithCustomPayload(c *gw.Context) {
	c.JSON400Payload(4000, "Custom 400 Payload")
}

func GetTester400WithCustomPayloadErr(c *gw.Context) {
	c.JSON400PayloadMsg(4000, "Custom 400 Err and Payload", gin.H{"A": "a"})
}

func GetTester401(c *gw.Context) {
	c.JSON401(4001)
}

func GetTester401WithCustomErr(c *gw.Context) {
	c.JSON401Msg(4001, "Custom 401 Err")
}

func GetTester401WithCustomPayload(c *gw.Context) {
	c.JSON401Payload(4001, "Custom 401 Payload")
}

func GetTester401WithCustomPayloadErr(c *gw.Context) {
	c.JSON401PayloadMsg(4001, "Custom 401 Err and Payload", gin.H{"A": "a"})
}

func GetTester403(c *gw.Context) {
	c.JSON403(4003)
}

func GetTester403WithCustomErr(c *gw.Context) {
	c.JSON403Msg(4003, "Custom 403 Err")
}

func GetTester403WithCustomPayload(c *gw.Context) {
	c.JSON403Payload(4003, "Custom 403 Payload")
}

func GetTester403WithCustomPayloadErr(c *gw.Context) {
	c.JSON403PayloadMsg(4003, "Custom 403 Err and Payload", gin.H{"A": "a"})
}

func GetTester404(c *gw.Context) {
	c.JSON404(4004)
}

func GetTester404WithCustomErr(c *gw.Context) {
	c.JSON404Msg(4004, "Custom 404 Err")
}

func GetTester404WithCustomPayload(c *gw.Context) {
	c.JSON404Payload(4004, "Custom 404 Payload")
}

func GetTester404WithCustomPayloadErr(c *gw.Context) {
	c.JSON404PayloadMsg(4004, "Custom 404 Err and Payload", gin.H{"A": "a"})
}

func GetTester500(c *gw.Context) {
	c.JSON500(5000)
}

func GetTester500WithCustomErr(c *gw.Context) {
	c.JSON500Msg(5000, "Custom 500 Err")
}

func GetTester500WithCustomPayload(c *gw.Context) {
	c.JSON500Payload(5000, "Custom 500 Payload")
}

func GetTester500WithCustomPayloadErr(c *gw.Context) {
	c.JSON500PayloadMsg(5000, "Custom 500 Err and Payload", gin.H{"A": "a"})
}
