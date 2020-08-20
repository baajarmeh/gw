package rest

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/tester/dto"
	"github.com/oceanho/gw/logger"
)

type MyTesterRestAPI struct {
}

// gw.IRestAPI
// Name returns the Resource name.
func (MyTesterRestAPI) Name() string {
	return "my-tester-default-RestAPI"
}

// http Get
func (MyTesterRestAPI) Get(c *gw.Context) {
	db := c.State.Store().GetDbStore()
	out := &[]dto.MyTester{}
	err := db.Model(&dto.MyTester{}).Limit(200).Offset(0).Scan(out).Error
	c.JSON(err, out)
}

// http Get before handler
func (MyTesterRestAPI) OnGetBefore() gw.DecoratorHandler {
	return func(ctx *gw.Context) (status int, err error, payload interface{}) {
		logger.Info("Exec func (MyTesterRestAPI) OnGetBefore")
		return gw.DecoratorHandler500.Result()
		// ``return 200, fmt.Errorf("tester before action quit"), gin.H{
		// 	"A": "a",
		// 	"B": "b",
		// }``
	}
}

// http Get after handler
func (MyTesterRestAPI) OnGetAfter() gw.DecoratorHandler {
	return func(ctx *gw.Context) (status int, err error, payload interface{}) {
		logger.Info("Exec func (MyTesterRestAPI) OnGetAfter")
		return gw.DecoratorHandler500.Result()
	}
}

// http Get pager Query.
func (MyTesterRestAPI) Query(c *gw.Context) {
	expr := &gw.QueryExpr{}
	if c.Bind(expr) != nil {
		return
	}
	db := c.State.Store().GetDbStore()
	out := &[]dto.MyTester{}
	var total int64
	err := db.Model(&dto.MyTester{}).Count(&total).Limit(expr.PageSize).Offset(expr.PageOffset()).Scan(out).Error
	if err != nil {
		c.Fault(err)
		return
	}
	c.PagerJSON(total, expr.PagerExpr, out)
}

type MyTesterGlobalDecorator struct {
}

func (m MyTesterGlobalDecorator) OnBeforeCall(ctx *gw.Context) (friendlyMsg string, err error) {
	return "", nil
}

func (m MyTesterGlobalDecorator) OnAfterCall(ctx *gw.Context) (friendlyMsg string, err error) {
	return "", nil
}

// http Post
func (MyTesterRestAPI) Post(c *gw.Context) {
	obj := &dto.MyTester{}
	err := c.State.Store().GetDbStore().Create(obj).Error
	c.JSON(err, obj)
}

//
// MyTesterRestAPI global decorators.
func (m MyTesterRestAPI) SetupDecorator() []gw.Decorator {
	var d []gw.Decorator
	return d
}

//
// More http Methods(Put,Delete and so so.)
//
