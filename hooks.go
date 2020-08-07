package gw

import (
	"github.com/gin-gonic/gin"
)

// Hook represents a global gin engine http Middleware.
type Hook struct {
	Name     string
	OnBefore gin.HandlerFunc
	OnAfter  gin.HandlerFunc
}

func NewBeforeHook(name string, before gin.HandlerFunc) *Hook {
	return NewHook(name, before, nil)
}

func NewAfterHook(name string, after gin.HandlerFunc) *Hook {
	return NewHook(name, nil, after)
}

func NewHook(name string, before gin.HandlerFunc, after gin.HandlerFunc) *Hook {
	return &Hook{
		Name:     name,
		OnBefore: before,
		OnAfter:  after,
	}
}
