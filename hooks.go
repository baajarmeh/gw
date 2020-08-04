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

type DecoratorPoint int8

const (
	DecoratorPointActionBefore DecoratorPoint = 0
	DecoratorPointActionAfter
)

type IDecorator interface {
	Catalog() string
	Point() DecoratorPoint
	Call(ctx *Context) (friendlyMsg string, err error)
}

type PermissionDecoratorImpl struct {
	perms       []Permission
	friendlyMsg string
}

const permissionDecoratorCatalog = "permission"

var ErrPermissionDenied = fmt.Errorf("permission denied")

func (p PermissionDecoratorImpl) Catalog() string {
	return permissionDecoratorCatalog
}

func (p PermissionDecoratorImpl) Point() DecoratorPoint {
	return DecoratorPointActionBefore
}

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

// helpers
func filterDecorator(filter func(d IDecorator) bool, decorators ...IDecorator) []IDecorator {
	var result []IDecorator
	for _, dc := range decorators {
		if filter(dc) {
			result = append(result, dc)
		}
	}
	return result
}
