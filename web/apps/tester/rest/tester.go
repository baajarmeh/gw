package rest

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/logger"
	"github.com/oceanho/gw/web/apps/tester/dto"
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
	db := c.Store.GetDbStore()
	out := &[]dto.MyTester{}
	err := db.Model(&dto.MyTester{}).Limit(200).Offset(0).Scan(out).Error
	c.JSON(err, out)
}

// http Get before handler
func (MyTesterRestAPI) OnGetBefore(c *gw.Context) {
	logger.Info("Exec func (MyTesterRestAPI) OnGetBefore(c *gw.Context)")
}

// http Get after handler
func (MyTesterRestAPI) OnGetAfter(c *gw.Context) {
	logger.Info("Exec func (MyTesterRestAPI) OnGetAfter(c *gw.Context)")
}

// http Get pager Query.
func (MyTesterRestAPI) Query(c *gw.Context) {
	expr := &gw.QueryExpr{}
	if c.Bind(expr) != nil {
		return
	}
	db := c.Store.GetDbStore()
	out := &[]dto.MyTester{}
	var total int64
	err := db.Model(&dto.MyTester{}).Count(&total).Limit(expr.PageSize).Offset(expr.PageOffset()).Scan(out).Error
	if err != nil {
		c.Fault(err)
		return
	}
	c.PagerJSON(total, expr.PagerExpr, out)
}

type MyTesterGlobalBeforeDecorator struct {
}
type MyTesterGlobalAfterDecorator struct {
}

func (m MyTesterGlobalBeforeDecorator) Catalog() string {
	return "my-tester-global-decorator"
}

func (m MyTesterGlobalBeforeDecorator) Point() gw.DecoratorPoint {
	return gw.DecoratorPointActionBefore
}

func (m MyTesterGlobalBeforeDecorator) Call(ctx *gw.Context) (friendlyMsg string, err error) {
	logger.Info("requestID %s, func (m MyTesterGlobalBeforeDecorator) Call(ctx *gw.Context) (friendlyMsg string, err error)", ctx.RequestID)
	return "", nil
}

func (m MyTesterGlobalAfterDecorator) Catalog() string {
	return "my-tester-global-decorator"
}

func (m MyTesterGlobalAfterDecorator) Point() gw.DecoratorPoint {
	return gw.DecoratorPointActionBefore
}

func (m MyTesterGlobalAfterDecorator) Call(ctx *gw.Context) (friendlyMsg string, err error) {
	logger.Info("requestID %s, func (m MyTesterGlobalAfterDecorator) Call(ctx *gw.Context) (friendlyMsg string, err error)", ctx.RequestID)
	return "", nil
}

// http Post
func (MyTesterRestAPI) Post(c *gw.Context) {
	obj := &dto.MyTester{}
	err := c.Store.GetDbStore().Create(obj).Error
	c.JSON(err, obj)
}

//
// MyTesterRestAPI global decorators.
func (m MyTesterRestAPI) SetupDecorator() []gw.IDecorator {
	var d []gw.IDecorator
	d = append(d, MyTesterGlobalBeforeDecorator{}, MyTesterGlobalAfterDecorator{})
	return d
}

//
// More http Methods(Put,Delete and so so.)
//
