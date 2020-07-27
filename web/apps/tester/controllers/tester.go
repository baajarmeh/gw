package controllers

import (
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/web/apps/tester/dto"
)

type MyTesterController struct {
}

func (MyTesterController) Name() string {
	return "my-tester-default-controller"
}

func (c *MyTesterController) Create(ctx *gw.Context) gw.IController {
	return c
}

func (c MyTesterController) OnDestroy(ctx *gw.Context) error {
	return nil
}

func (MyTesterController) Get(c *gw.Context) {
	db := c.Store.GetDbStore()
	out := &[]dto.MyTester{}
	err := db.Model(&dto.MyTester{}).Limit(200).Offset(0).Scan(out).Error
	c.JSON(err, out)
}

func (MyTesterController) Query(c *gw.Context) {
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

func (MyTesterController) Post(c *gw.Context) {
	obj := &dto.MyTester{}
	err := c.Store.GetDbStore().Create(obj)
	c.JSON(err, obj)
}
