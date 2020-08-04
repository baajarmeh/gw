package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Hook struct {
	Name    string
	Handler gin.HandlerFunc
}

func NewHook(name string, handler gin.HandlerFunc) *Hook {
	return &Hook{
		Name:    name,
		Handler: handler,
	}
}

type IDecorator interface {
	Call(ctx *Context) (friendlyMsg string, err error)
}

type PermissionDecoratorImpl struct {
	perms       []Permission
	friendlyMsg string
}

var ErrPermissionDenied = fmt.Errorf("permission denied")

func (p PermissionDecoratorImpl) Call(ctx *Context) (friendlyMsg string, err error) {
	s := hostServer(ctx.Context)
	if !s.permissionManager.HasPermission(ctx.User, p.perms...) {
		return p.friendlyMsg, ErrPermissionDenied
	}
	return "", nil
}

func PermissionDecorator(perms ...Permission) IDecorator {
	names := make([]string, len(perms))
	for idx := 0; idx < len(perms); idx++ {
		names[idx] = perms[idx].Name
	}
	friendlyMsg := fmt.Sprintf("Permission Deined, needs:(%s)", strings.Join(names, "|"))
	return &PermissionDecoratorImpl{
		perms:       perms,
		friendlyMsg: friendlyMsg,
	}
}
