package gw

import "testing"

type TesterController struct {
}

func (controller TesterController) Get(c *Context) {

}

func TestRegisterRouteWithControllers(t *testing.T) {
	r := &RouterGroup{}
	RegisterControllers(r, TesterController{})
}
