package gw

import (
	"github.com/gin-gonic/gin"
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
